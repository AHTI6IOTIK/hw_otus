package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AHTI6IOTIK/hw_otus/hw11_telnet_client/cmd"
)

func init() {
	log.SetFlags(^log.Lshortfile & ^log.Llongfile & ^log.Ldate & ^log.Ltime & ^log.Lmicroseconds)
	log.SetPrefix("...")
	log.SetOutput(os.Stderr)
}

func main() {
	ctx := context.Background()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go cmd.Execute(ctx, quit)

	<-quit
}
