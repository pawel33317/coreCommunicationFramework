package sys_signal_receiver

import (
	"os"
	"os/signal"
	"syscall"
)

func ReceiveTerminationSignal(receivedChannel chan<- bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-sigs

			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				receivedChannel <- true
			}
		}
	}()

}
