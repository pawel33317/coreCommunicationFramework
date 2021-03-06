package sys_signal_listener

import (
	"os"
	"os/signal"
	"syscall"
)

func ListenOnTerminationSignal(receivedSignalChan chan<- struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-sigs

			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				receivedSignalChan <- struct{}{}
			}
		}
	}()

}
