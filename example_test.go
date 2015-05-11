package backoff_test

import (
	"fmt"
	"time"

	"github.com/azr/backoff"
)

func ExampleExponentialSleepTimesOutput() {
	exp := backoff.NewExponential()
	exp.MaxInterval = time.Second * 25

	for i := 0; i < 12; i++ {
		d := exp.GetSleepTime()
		fmt.Printf("Random duration was %s, interval is %s in [%s , %s]\n",
			d.String(),
			exp.Inverval().String(),
			(exp.Inverval() - time.Duration(exp.RandomizationFactor*float64(exp.Inverval()))).String(),
			(exp.Inverval() + time.Duration(exp.RandomizationFactor*float64(exp.Inverval()))).String(),
		)
		exp.IncrementCurrentInterval()
	}
	// Output:
	// Random duration was 507.606314ms, interval is 500ms in [250ms , 750ms]
	// Random duration was 985.229971ms, interval is 750ms in [375ms , 1.125s]
	// Random duration was 803.546856ms, interval is 1.125s in [562.5ms , 1.6875s]
	// Random duration was 1.486109007s, interval is 1.6875s in [843.75ms , 2.53125s]
	// Random duration was 2.070709754s, interval is 2.53125s in [1.265625s , 3.796875s]
	// Random duration was 3.67875363s, interval is 3.796875s in [1.8984375s , 5.6953125s]
	// Random duration was 4.459624189s, interval is 5.6953125s in [2.84765625s , 8.54296875s]
	// Random duration was 6.775444383s, interval is 8.54296875s in [4.271484375s , 12.814453125s]
	// Random duration was 15.10932531s, interval is 12.814453125s in [6.407226563s , 19.221679687s]
	// Random duration was 13.811796615s, interval is 19.221679687s in [9.610839844s , 28.83251953s]
	// Random duration was 17.579671916s, interval is 25s in [12.5s , 37.5s]
	// Random duration was 21.521785421s, interval is 25s in [12.5s , 37.5s]
}
