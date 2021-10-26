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
	start := time.Now()

	// Create ratelimited errgroup with max 10 executions per second
	eg, ctx := ratelimit.WithContext(ctx, 10)

	// Kick off 10 goroutines.
	for i := 0; i < 10; i++ {
		// Shadow i so that it can be used in concurrent goroutines without
		// future loop iterations changing its value.
		i := i
		eg.Go(func() error {
			fmt.Printf("%d: %s\n", i, time.Since(start).Round(time.Millisecond*10))
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		panic(err)
	}

	// Output:
	// 0: 100ms
	// 1: 200ms
	// 2: 300ms
	// 3: 400ms
	// 4: 500ms
	// 5: 600ms
	// 6: 700ms
	// 7: 800ms
	// 8: 900ms
	// 9: 1s
}
```
