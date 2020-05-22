[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/cyrusaf/ratelimit)

# ratelimit
`ratelimit` implements the `x/sync/errgroup` interface, but enforces a max throughput
that the functions will be executed at. This is useful when making concurrent requests to
backends that have either have rate limiting or cannot handle excessive throughput.

## Usage
See the [godoc](https://godoc.org/github.com/cyrusaf/ratelimit) page for full documentation.
```golang
import (
    "context"
    "fmt"
    "time"
    
    "github.com/cyrusaf/ratelimit"
)

func ExampleErrGroup() {
	ctx := context.Background()

	// Create ratelimited errgroup with max 10 executions per second
	eg, ctx := ratelimit.WithContext(ctx, 10)

	// Kick off 10 goroutines.
	for i := 0; i < 10; i++ {
		// Shadow i so that it can be used in concurrent goroutines without
		// future loop iterations changing its value.
		i := i
		eg.Go(func() error {
			fmt.Printf("%d: %s\n", i, time.Now().Format(time.StampMilli))
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		panic(err)
	}

	// Output:
	// 0: May 22 14:16:34.144
	// 1: May 22 14:16:34.244
	// 2: May 22 14:16:34.346
	// 3: May 22 14:16:34.441
	// 4: May 22 14:16:34.544
	// 5: May 22 14:16:34.643
	// 6: May 22 14:16:34.741
	// 7: May 22 14:16:34.842
	// 8: May 22 14:16:34.945
	// 9: May 22 14:16:35.041
}
```