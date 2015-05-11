package backoff

import (
	"math/rand"
	"time"
)

/*
ExponentialBackOff is an implementation of BackOff that backf of and increases
it's the back off period for each retry attempt using a randomization function
that grows exponentially.
Backoff() time is calculated using the following formula:
    randomized_interval =
        retry_interval * (random value in range [1 - randomization_factor, 1 + randomization_factor])
In other words BackOff() will sleep for times between the randomization factor
percentage below and above the retry interval.
For example, using 2 seconds as the base retry interval and 0.5 as the
randomization factor, the actual back off period used in the next retry
attempt will be between 1 and 3 seconds.

Note: max_interval caps the retry_interval and not the randomized_interval.

If the time elapsed since an ExponentialBackOff instance is created goes past the
max_elapsed_time then the method NextBackOff() starts returning backoff.Stop.
The elapsed time can be reset by calling Reset().

Example: The default retry_interval is .5 seconds, default randomization_factor is 0.5, default
multiplier is 1.5 and the default max_interval is 1 minute. For 10 tries the sequence will be
(values in seconds) and assuming we go over the max_elapsed_time on the 10th try:

    request#     retry_interval     randomized_interval

    1             0.5                [0.25,   0.75]
    2             0.75               [0.375,  1.125]
    3             1.125              [0.562,  1.687]
    4             1.687              [0.8435, 2.53]
    5             2.53               [1.265,  3.795]
    6             3.795              [1.897,  5.692]
    7             5.692              [2.846,  8.538]
    8             8.538              [4.269, 12.807]
    9            12.807              [6.403, 19.210]
    10           19.210              backoff.Stop
Implementation is not thread-safe.
*/
type ExponentialBackOff struct {
	InitialInterval time.Duration
	currentInterval time.Duration
	MaxInterval     time.Duration

	RandomizationFactor float64
	Multiplier          float64
}

// Default values for ExponentialBackOff.
const (
	DefaultInitialInterval     = 500 * time.Millisecond
	DefaultRandomizationFactor = 0.5
	DefaultMultiplier          = 1.5
	DefaultMaxInterval         = 60 * time.Second
)

// NewExponentialBackOff creates an instance of ExponentialBackOff using default values.
func NewExponentialBackOff() *ExponentialBackOff {
	b := &ExponentialBackOff{
		InitialInterval:     DefaultInitialInterval,
		RandomizationFactor: DefaultRandomizationFactor,
		Multiplier:          DefaultMultiplier,
		MaxInterval:         DefaultMaxInterval,
		currentInterval:     DefaultInitialInterval,
	}
	b.Reset()
	return b
}

// Reset the interval back to the initial retry interval and restarts the timer.
func (b *ExponentialBackOff) Reset() {
	b.currentInterval = b.InitialInterval
}

func (b *ExponentialBackOff) GetSleepTime() time.Duration {
	return getRandomValueFromInterval(b.RandomizationFactor, rand.Float64(), b.currentInterval)
}

func (b *ExponentialBackOff) BackOff() {
	time.Sleep(b.GetSleepTime())

	b.IncrementCurrentInterval()
}

// Increments the current interval by multiplying it with the multiplier.
func (b *ExponentialBackOff) IncrementCurrentInterval() {
	// Check for overflow, if overflow is detected set the current interval to the max interval.
	if float64(b.currentInterval) >= float64(b.MaxInterval)/b.Multiplier {
		b.currentInterval = b.MaxInterval
	} else {
		b.currentInterval = time.Duration(float64(b.currentInterval) * b.Multiplier)
	}
}

// Returns a random value from the interval:
//  [randomizationFactor * currentInterval, randomizationFactor * currentInterval].
func getRandomValueFromInterval(randomizationFactor, random float64, currentInterval time.Duration) time.Duration {
	var delta = randomizationFactor * float64(currentInterval)
	var minInterval = float64(currentInterval) - delta
	var maxInterval = float64(currentInterval) + delta
	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	return time.Duration(minInterval + (random * (maxInterval - minInterval + 1)))
}
