package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	consumer "github.com/tientp-floware/mgodb-stream/transport/consumer"
	transport "github.com/tientp-floware/mgodb-stream/transport/http"
	logger "go.uber.org/zap"
)

var (
	log = logger.GetLogger("[Gateway service]")
)

func main() {
	e := transport.NewHTTP().Server()
	mgostream := consumer.NewMgoStream()
	quit := make(chan os.Signal)

	port := "9090"
	// Start server
	go func() {
		if err := e.Start(":" + port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()
	// run stream with concurrency
	go mgostream.FlowChangeStream()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 30 seconds.
	signal.Notify(quit, os.Interrupt)
	quitMsg := <-quit
	log.Error(quitMsg.String())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
