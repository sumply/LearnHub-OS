package app

import (
	"net/http"
	"server/internal/config"
	"server/internal/logger"
	"server/internal/rest"
)

func Run() error {
	app, err := config.NewApp("./configs/app.yaml")
	if err != nil {
		return err
	}

	if err := InitLogger(app); err != nil {
		return err
	}

	r, err := rest.NewRouterStub()
	if err != nil {
		return err
	}

	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}

func InitLogger(app *config.App) error {
	new, err := app.CreateLoggerNewFunc()
	if err != nil {
		return err
	}
	logger.SetNewFunc(func() logger.Logger {
		return new()
	})
	level, err := app.CreateLoggerLevel()
	if err != nil {
		return err
	}
	logger.SetLayer(level)
	return nil
}
