package hsdp

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func StartLogger(group, stream string) (chan cloudwatchlogs.OutputLogEvent, chan bool, error) {
	logChan := make(chan cloudwatchlogs.OutputLogEvent)
	doneChan := make(chan bool)

	go func() {
		for {
			select {
			case msg := <-logChan:
				fmt.Printf("CustomLogEvent|%s|%s|%s|%s|%s\n", "1", "2", "3", stream, *msg.Message)
			case <-doneChan:
				fmt.Printf("exiting logger\n")
				return
			}
		}
	}()

	return logChan, doneChan, nil
}
