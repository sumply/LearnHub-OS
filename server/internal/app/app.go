package app

import (
	"net/http"
	"server/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() error {
<<<<<<< Updated upstream
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	s, err := config.NewServer()
	if err != nil {
		return err
	}
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Addr, s.Port), r)
	return nil
=======
	app, err := config.NewApp("./configs/app.yaml")
	if err != nil {
		return err
	}

	if err := InitLogger(app); err != nil {
		return err
	}

	r, err := createHandler(app)
	if err != nil {
		return err
	}

	s := config.Server{Addr: "127.0.0.1", Port: 8000}
	return http.ListenAndServe(s.String(), r)
}

func createHandler(app *config.App) (http.Handler, error) {
	uu, err := app.CreateUsecaseUser()
	if err != nil {
		return nil, err
	}
	ua, err := app.CreateUsecaseAnswer()
	if err != nil {
		return nil, err
	}
	ug, err := app.CreateUsecaseGroup()
	if err != nil {
		return nil, err
	}
	us, err := app.CreateUsecaseSubject()
	if err != nil {
		return nil, err
	}
	uq, err := app.CreateUsecaseQuiz()
	if err != nil {
		return nil, err
	}
	tp, err := app.CreateTransportJWTParser()
	if err != nil {
		return nil, err
	}

	r, err := transport.NewRouter(
		uu,
		ug,
		us,
		uq,
		ua,
		tp,
	)
	if err != nil {
		return nil, err
	}
	return r, nil
>>>>>>> Stashed changes
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
