goalbatch
================
[![](https://godoc.org/github.com/upbit/goalbatch?status.svg)](http://godoc.org/github.com/upbit/goalbatch)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/upbit/goalbatch/blob/master/LICENSE.md)

goalbatch - A simple way to execute functions asynchronously and waits for results

## Batch
> Batch method returns when all of the callbacks passed or context is done, returned responses and errors are ordered according to callback order

```go
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)
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
