package goalbatch

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBatch(t *testing.T) {
	Convey("It should return rs in the same order as provided", t, func() {
		g := New(context.Background())
		rs, errs := g.Batch(
			func(ctx context.Context) (interface{}, error) {
				return 1, nil
			},
			func(ctx context.Context) (interface{}, error) {
				return nil, errors.New("failed")
			},
			func(ctx context.Context) (interface{}, error) {
				return 3, nil
			},
		)

		// t.Log(rs)
		// t.Log(errs)

		So(rs, ShouldHaveLength, 3)
		So(rs[0], ShouldEqual, 1)
		So(rs[1], ShouldBeNil)
		So(rs[2], ShouldEqual, 3)

		So(errs, ShouldHaveLength, 3)
		So(errs[0], ShouldBeNil)
		So(errs[1].Error(), ShouldEqual, "failed")
		So(errs[2], ShouldBeNil)
	})
}

func TestBatchWithTimeout(t *testing.T) {
	newAsyncFunc := func(ctx context.Context, param int) AsyncFunc {
		return func(ctx context.Context) (interface{}, error) {
			switch param {
			case 0:
				// 0 - failed
				return nil, errors.New("failed")
			case 3:
				// 3 - timeoutd
				time.Sleep(5 * time.Second)
			}
			result := fmt.Sprintf("param=%d", param+1)
			return result, nil
		}
	}

	Convey("It should return as soon as possible when context timeoutd", t, func() {
		timeout := time.Duration(100) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Gen async functions
		fns := make([]AsyncFunc, 4)
		for i, val := range []int{0, 1, 2, 3} {
			fns[i] = newAsyncFunc(ctx, val)
		}

		g := New(ctx)

		start := time.Now()
		rs, errs := g.Batch(fns...)
		cost := time.Since(start)

		So(cost, ShouldBeBetweenOrEqual, 100*time.Millisecond, 500*time.Millisecond)

		// t.Log(rs)
		// t.Log(errs)

		So(rs, ShouldHaveLength, 4)
		So(rs[0], ShouldBeNil)
		So(rs[1], ShouldEqual, "param=2")
		So(rs[2], ShouldEqual, "param=3")
		So(rs[3], ShouldBeNil)

		So(errs, ShouldHaveLength, 4)
		So(errs[0].Error(), ShouldEqual, "failed")
		So(errs[1], ShouldBeNil)
		So(errs[2], ShouldBeNil)
		So(errs[3], ShouldBeNil)
	})
}

func TestNew(t *testing.T) {
	g := New(nil)

	if g == nil {
		t.Fail()
	}
}

func TestNewWithContext(t *testing.T) {
	g := New(context.Background())

	if g == nil {
		t.Fail()
	}
}
