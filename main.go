package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"template/httpLocal"
	"template/log"
	token "template/token"

	_ "go.uber.org/automaxprocs"
)

func init() {
	err := os.Setenv("TZ", "Asia/Tehran")
	if err != nil {
		log.ErrLogger.Sugar().Panicf("panic: %w", err)
	}
}

func main() {
	ctx, exit := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup

	defer func() {
		httpLocal.HttpServer.Shutdown(ctx)
		httpLocal.HttpClient.CloseIdleConnections()
		log.Logger.Sync()
		log.ServLogger.Sync()
		log.ErrLogger.Sync()
		wg.Wait()
		exit()
	}()
	httpLocal.HttpServer.POST("/generate", token.Generate)
	httpLocal.HttpServer.POST("/authorize", token.Authorize)
	httpLocal.HttpServer.POST("/verify", token.Verify)
	httpLocal.HttpServer.POST("/revoke", token.Revoke)
	httpLocal.HttpServer.Logger.Fatal(httpLocal.HttpServer.Start(":7050"))
}
