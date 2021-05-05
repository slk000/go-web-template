package main

import (
	"context"
	"fmt"
	"go-web-template/dao/mysql"
	"go-web-template/dao/redis"
	"go-web-template/logger"
	"go-web-template/routes"
	"go-web-template/settings"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configFilename := ""
	if len(os.Args) == 2 {
		configFilename = os.Args[1]
	}
	// Load config
	if err := settings.Init(configFilename); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
	}
	// Init logger
	if err := logger.Init(settings.Conf.LoggerConfig); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
	}
	defer zap.L().Sync()

	// Init database
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
	}
	defer mysql.Close()

	// Init redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
	}
	defer redis.Close()

	// Register routers
	r := routes.Setup()
	// Start service
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server shutdown failed", zap.Error(err))
	}
	zap.L().Info("Server exiting")

}
