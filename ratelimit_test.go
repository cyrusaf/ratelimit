package ratelimit_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cyrusaf/ratelimit"
)

func TestRateLimitErrGroup(t *testing.T) {
	ctx := context.Background()
	eg, ctx := ratelimit.WithContext(ctx, 100)

	start := time.Now()
	for i := 0; i < 10; i++ {
		eg.Go(func() error {
			return nil
		})
	}
	err := eg.Wait()
	diff := time.Since(start)

	if err != nil {
		t.Fatalf("expected err to be nil, but got %v instead", err)
	}

	if diff.Round(10*time.Millisecond) != 100*time.Millisecond && diff > time.Millisecond*99 {
		t.Fatalf("expected ratelimit to take 100ms to execute all functions, but took %v instead", diff.Round(100*time.Millisecond))
	}
}

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
