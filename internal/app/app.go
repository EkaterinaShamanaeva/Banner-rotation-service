package app

import (
	"context"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/storage/sqlstorage"
)

type App struct {
	logger  Logger
	storage sqlstorage.Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

func New(logger Logger, storage sqlstorage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Info("storage closing...")
	return a.storage.Close(ctx)
}
