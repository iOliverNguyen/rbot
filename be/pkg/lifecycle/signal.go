package lifecycle

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olvrng/rbot/be/pkg/l"
)

var ll = l.New()

func ListenForSignal(cancel func(), max time.Duration, sigs ...os.Signal) {
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, sigs...)
		ll.Info("Received OS signal", l.Stringer("signal", <-osSignal))

		if max > 0 {
			go func() {
				timer := time.NewTimer(max)
				<-timer.C
				ll.Fatal("force shutdown due to timeout!")
			}()
		}
		cancel()
	}()
}
