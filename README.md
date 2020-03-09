goalbatch
================
[![goalbatch badge](https://github.com/upbit/goalbatch/workflows/goalbatch/badge.svg)](https://github.com/upbit/goalbatch/actions?query=workflow%3Agoalbatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/upbit/goalbatch)](https://goreportcard.com/report/github.com/upbit/goalbatch)
[![codecov](https://codecov.io/gh/upbit/goalbatch/branch/master/graph/badge.svg)](https://codecov.io/gh/upbit/goalbatch)
[![](https://godoc.org/github.com/upbit/goalbatch?status.svg)](http://godoc.org/github.com/upbit/goalbatch)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/upbit/goalbatch/blob/master/LICENSE)

goalbatch - A simple way to execute functions asynchronously and waits for results

## Batch
> Batch method returns when all of the callbacks passed or context is done, returned responses and errors are ordered according to callback order

```go
	timeout := time.Duration(100) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	g := goalbatch.New(ctx)
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

	fmt.Println(rs)
	fmt.Println(errs)
	// Output:
	// [<nil> <nil> 3 4]
	// [<nil> "failed" <nil> <nil>]
```

Or generate closure functions with parameters:

```go
	newAsyncFunc := func(ctx context.Context, param1 string, param2 int) AsyncFunc {
		return func(ctx context.Context) (interface{}, error) {
			// deal with param1, param2...
			result := fmt.Sprintf("p1=%s p2=%d", param1, param2)
			return result, nil
		}
	}

	fns := make([]AsyncFunc, 2)
	fns[0] = newAsyncFunc(ctx, "foo", 1)
	fns[1] = newAsyncFunc(ctx, "bar", 2)

	g := goalbatch.New(ctx)
	rs, errs := g.Batch(fns...)

	fmt.Println(rs)
	fmt.Println(errs)
	// Output:
	// ["p1=foo p2=1" "p1=bar p2=2"]
	// [<nil> <nil>]
```
