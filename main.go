package main

import (
	"bloom/bloom"
	"bloom/view"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func getShutdownSignal() *chan os.Signal {
	shutdownSignal := make(chan os.Signal)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return &shutdownSignal
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	shutdown := getShutdownSignal()

	manager := view.CreateNodes(
		[]*view.NodeInput{
			view.GetViewInputRest(view.RestDaemonInput{Addr: ":8080"}),
		},
		bloom.NodeInitInput{
			HashPrefix:    []byte("prefix"),
			MinEntryCount: 1000000000,
		},
	)

	err := manager.Setup()
	if err != nil {
		log.Fatal(err)
	}

	<-*shutdown
	err = manager.Shutdown()
	if err != nil {
		log.Fatal(err)
	}
}
