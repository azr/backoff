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

// Indicates that no more retries should be made for use in NextBackOff().
const Stop time.Duration = -1

// ZeroBackOff is a fixed back-off policy whose back-off time is always zero,
// meaning that the operation is retried immediately without waiting.
type ZeroBackOff struct{}

var _ BackOffer = (*ZeroBackOff)(nil)

func (b *ZeroBackOff) Reset() {}

func (b *ZeroBackOff) BackOff() {}

type ConstantBackOff struct {
	Interval time.Duration
}

func (b *ConstantBackOff) Reset()                     {}
func (b *ConstantBackOff) NextBackOff() time.Duration { return b.Interval }

func NewConstantBackOff(d time.Duration) *ConstantBackOff {
	return &ConstantBackOff{Interval: d}
}
