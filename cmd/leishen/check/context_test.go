package check

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestContext(t *testing.T) {
	parallelism := 4
	grounp, ctx := errgroup.WithContext(context.Background())

	ch := make(chan uint64, parallelism)
	for i := 0; i < parallelism; i++ {
		grounp.Go(func() error {
			return runer(ctx, ch)
		})
	}

	grounp.Go(func() error {
		for i := uint64(0); i < 12; i++ {
			select {
			case <-ctx.Done():
				return nil
			case ch <- i:
			}
		}
		close(ch)
		return nil
	})

	fmt.Println(grounp.Wait())
}

func runer(ctx context.Context, ch chan uint64) error {
	var ok bool
	var number uint64
	for {
		select {
		case <-ctx.Done():
			return nil
		case number, ok = <-ch:
			if !ok {
				return nil
			}
		}

		if number > 10 {
			return fmt.Errorf("number greater than 10")
		}
	}
}
