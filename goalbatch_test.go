package goalbatch

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBatch(t *testing.T) {
	Convey("It should return rs in the same order as provided by Batch()", t, func() {
		timeout := time.Duration(100) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		g := New(ctx)
		rs, errs := g.Batch(
			func(ctx context.Context) (interface{}, error) {
				time.Sleep(5 * time.Second)
				return 1, nil
			},
			func(ctx context.Context) (interface{}, error) {
				return nil, errors.New("failed")
			},
			func(ctx context.Context) (interface{}, error) {
				return 3, nil
			},
			func(ctx context.Context) (interface{}, error) {
				return 4, nil
			},
		)

		t.Log(rs)
		t.Log(errs)
		// Output:
		//   rs [<nil> <nil> 3 4]
		// errs [<nil> "failed" <nil> <nil>]

		So(rs, ShouldHaveLength, 4)
		So(rs[0], ShouldBeNil)
		So(rs[1], ShouldBeNil)
		So(rs[2], ShouldEqual, 3)
		So(rs[3], ShouldEqual, 4)

		So(errs, ShouldHaveLength, 4)
		So(errs[0], ShouldBeNil)
		So(errs[1].Error(), ShouldEqual, "failed")
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
