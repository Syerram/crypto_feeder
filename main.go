package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

var LOGGER *Logger

func init() {
	LOGGER = NewLogger(TRACE, "")
	LOGGER.Trace.Println("Setting up data feeds")
}

func main() {
	// Close channel will be used by exchanges and writers to shut down before ending the program
	var closeChannel = make(chan bool)
	// Channel will be used for transferring data between exchanges to writers
	var broadcastChannel = make(chan []TradeRow)

	// Handle interrupt. Let others know its time to close
	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c // waits infinitely till interrupt signal is sent
		fmt.Println("\nCleaning up. Sending close signal")
		closeChannel <- true // Inform others to close now. Sleep few seconds before exiting
		time.Sleep(5)
		close(closeChannel)
		close(broadcastChannel)
		os.Exit(0)
	}()

	// Let the show begin.
	// Writers and Exchanges are unaware of each other and only communicate via channel
	WRITERS.Listen(closeChannel)
	EXCHANGES.Listen(closeChannel, broadcastChannel)
	fmt.Println("Working...")

	// Loop infinitely till close channel is called
	// Otherwise, listen for new trades and on its arrival have writers handle it in a separate routine
	for {
		select {
		case <-closeChannel:
			break
		case tradeRows := <-broadcastChannel:
			// Launch a new go routine for this, so its not blocking the write
			go func() {
				WRITERS.Write(tradeRows)
			}()
		}
	}
}
