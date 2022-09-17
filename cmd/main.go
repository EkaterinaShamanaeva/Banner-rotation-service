package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/app"
	configuration "github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/config"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/logger"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/server"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/storage/sqlstorage"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := configuration.NewConfig()

	if err := config.BuildConfig(configFile); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("Logger error: %v", err)
	}

	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.Database.Username,
		config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name,
		config.Database.SSLMode)
	fmt.Println(dsn)

	sqlStorage := sqlstorage.New()
	dbpool, err := sqlstorage.Connect(ctx, dsn)
	if err != nil {
		logg.Error("failed to connect DB: " + err.Error())
	}
	sqlStorage.Pool = dbpool
	defer sqlStorage.Pool.Close()

	logg.Info("DB connected...")

	rotator := app.New(logg, *sqlStorage)

	httpServer := server.NewServer(logg, rotator)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err = rotator.Close(ctx); err != nil {
			logg.Error("failed close storage: " + err.Error())
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		addrServer := net.JoinHostPort(config.Server.Host, config.Server.Port)
		if err = httpServer.Start(ctx, addrServer); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	<-ctx.Done()
	wg.Wait()
}
