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

	"faissal.com/blogSpace/internal/auth"
	"faissal.com/blogSpace/internal/services"
	"faissal.com/blogSpace/internal/uploader"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type DBConf struct {
	Addr        string
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime string
}

type JwtConfig struct {
	SecretKey string
	Iss       string
	Sub       string

	// Unix time
	Exp int64
}

type R2Conf struct {
	BucketName      string
	AccountId       string
	AccessKeyId     string
	AccessKeySecret string
}

type Application struct {
	Port string

	Host string

	Env string

	DbConfig DBConf

	Services services.Services

	JwtAuth JwtConfig

	SwaggerUrl string

	Authentication auth.Authenticator

	Uploading uploader.Uploader

	R2Config R2Conf
}

func (app *Application) Mux() http.Handler {

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Route("/", func(r chi.Router) {
			r.Use(app.AuthMiddleware)

			// Users Resource
			r.Get("/users/profile", app.GetUserProfileHandler)

			r.Delete("/users/delete", app.DeleteUserAccountHandler)

			// Categories Resource
			r.Post("/categories", app.CreateCategoryHandler)

			r.Get("/categories/{ID}", app.GetCategoryMiddleware(app.GetCategoryByIdHandler))

			r.Patch("/categories/{ID}", app.GetCategoryMiddleware(app.UpdateCategoryHandler))

			r.Delete("/categories/{ID}", app.GetCategoryMiddleware(app.DeleteCategoryHandler))

			// Blogs Resource
			r.Post("/blogs", app.CreateBlogHandler)

			// Comments Resource
			r.Post("/comments", app.CreateCommentHandler)
			r.Delete("/comments/{ID}", app.CheckOwnerCommentMiddleware(app.DeleteCommentHandler))
		})

		// Public Routes
		r.Get("/ping", app.PingHandler)

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(app.SwaggerUrl)))

		r.Get("/categories", app.GetCategoriesHandler)

		r.Route("/authentication", func(r chi.Router) {

			r.Post("/sign-up", app.SignUpHandler)

			r.Post("/activation/{token}", app.ActivateAccountHandler)

			r.Post("/sign-in", app.SignInHandler)

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
