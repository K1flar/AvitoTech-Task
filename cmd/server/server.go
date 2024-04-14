package server

import (
	"banner_service/api"
	"banner_service/internal/config"
	"banner_service/internal/domains"
	"banner_service/internal/handlers"
	"banner_service/internal/repositories/postgres"
	"banner_service/internal/services"
	tokenmanager "banner_service/internal/token_manager"
	"banner_service/pkg/cache/lfu"
	adminmw "banner_service/pkg/middlewares/admin_mw"
	"banner_service/pkg/middlewares/auth"
	loggermw "banner_service/pkg/middlewares/logger_mw"
	ratelimit "banner_service/pkg/middlewares/rate_limit"
	timelimit "banner_service/pkg/middlewares/time_limit"
	"banner_service/pkg/mux"
	"fmt"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	cfg    *config.Config
	Server *http.Server
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	tokenManager := tokenmanager.New(cfg.Server.UserToken, cfg.Server.AdminToken)

	repository, err := postgres.New(&cfg.Database)
	if err != nil {
		return nil, err
	}

	cache := lfu.NewWithLifeCycle[domains.BannerKey, *domains.Banner](1000, cfg.Server.BanneLifeCycle)

	service := services.New(log, repository, cache)

	handler := handlers.New(service, log)

	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}
	swagger.Servers = nil

	r := mux.New()

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		b, _ := swagger.MarshalJSON()
		w.Write(b)
	})
	r.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))

	r.Use(ratelimit.New(cfg.Server.RPS))
	r.Use(timelimit.New(cfg.Server.ResponseTime))
	r.Use(loggermw.New(log))

	r.Group(func(m *mux.Mux) {
		m.Use(auth.New(tokenManager))

		m.HandleFunc("GET /user_banner", handler.GetUserBanner)

		m.Group(func(adminHandler *mux.Mux) {
			adminHandler.Use(adminmw.New())

			adminHandler.HandleFunc("POST /banner", handler.PostBanner)
			adminHandler.HandleFunc("GET /banner", handler.GetBanner)
			adminHandler.HandleFunc("PATCH /banner/{id}", handler.PatchBannerId)
			adminHandler.HandleFunc("DELETE /banner/{id}", handler.DeleteBannerId)
		})
	})

	return &App{
		cfg: cfg,
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler: r,
		},
	}, nil
}

func (a *App) Run() error {
	fmt.Printf("Starting server on %s:%s...\n", a.cfg.Server.Host, a.cfg.Server.Port)
	err := a.Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
