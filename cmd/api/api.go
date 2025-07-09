package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type Application struct {
	Port     string
	Host     string
	Env      string
	DbConfig DBConf
}

type DBConf struct {
	Addr        string
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime string
}

func (app *Application) Mux() http.Handler {

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

	})

	return r
}

func (app *Application) Run(mux http.Handler) error {

	srv := http.Server{
		Addr:         net.JoinHostPort(app.Host, app.Port),
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		log.Info("signal cought, %v", s.String())

		shutdown <- srv.Shutdown(ctx)

	}()

	log.Info("server has started at", "host", app.Host, "port", app.Port)
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	log.Info("server has stopped at", "host", app.Host, "env", app.Env)

	return nil
}
