//Package backoff helps you at backing off !
//
//It was forked from github.com/cenkalti/backoff which is very awesome.
//
//I wanted a BackOff that'd just sleep and do nothing else.
package backoff

import "time"

// BackOffer interface to use after a retryable operation failed.
// A Backoffer.BackOff sleeps.
type BackOffer interface {
	// Example usage:
	//
	//   for ;; {
	//       err, canRetry := somethingThatCanFail()
	//       if err != nil && canRetry {
	//           backoffer.Backoff()
	//       }
	//   }
	BackOff()

	// Reset to initial state.
	Reset()
}

// ZeroBackOff is a fixed back-off policy whose back-off time is always zero,
// meaning that the operation is retried immediately without waiting.
type ZeroBackOff struct{}

var _ BackOffer = (*ZeroBackOff)(nil)

func (b *ZeroBackOff) Reset() {}

func (b *ZeroBackOff) BackOff() {}

type ConstantBackOff struct {
	Interval time.Duration
}

var _ BackOffer = (*ConstantBackOff)(nil)

func (b *ConstantBackOff) Reset() {}

func (b *ConstantBackOff) BackOff() {
	time.Sleep(b.Interval)
}

func NewConstant(d time.Duration) *ConstantBackOff {
	return &ConstantBackOff{Interval: d}
}
