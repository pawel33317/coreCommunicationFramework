package test_objects

import (
	"time"

	"github.com/pawel33317/coreCommunicationFramework/logger"
)

func Print_data_time_parallely(l logger.Logger) {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			l.Log(logger.INFO, "DATE", "Current date and time is: ", time.Now().String())
		}
	}()
}
