package main

import (
	"banner_service/api"
	"banner_service/internal/config"
	"banner_service/internal/domains"
	"banner_service/internal/handlers"
	"banner_service/internal/logger"
	"banner_service/internal/repositories/postgres"
	"banner_service/internal/services"
	tokenManager "banner_service/internal/token_manager"
	"banner_service/pkg/cache/lfu"
	adminmw "banner_service/pkg/middlewares/admin_mw"
	"banner_service/pkg/middlewares/auth"
	loggermw "banner_service/pkg/middlewares/logger_mw"
	"banner_service/pkg/mux"
	"fmt"
	"os"

	// _ "banner_service/docs"

	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

const configPath = "configs/config.yaml"

func main() {
	cfg, err := config.New(configPath)
	exitOnError(err)
	tokenManager := tokenManager.New(cfg.Server.UserToken, cfg.Server.AdminToken)

	log := logger.New()

	repository, err := postgres.New(&cfg.Database)
	exitOnError(err)

	fmt.Println(cfg.Server.BanneLifeCycle)
	cache := lfu.NewWithLifeCycle[domains.BannerKey, *domains.Banner](1000, cfg.Server.BanneLifeCycle)

	service := services.New(log, repository, cache)

	handler := handlers.New(service, log)

	swagger, err := api.GetSwagger()
	exitOnError(err)
	swagger.Servers = nil

	r := mux.New()

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

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		b, _ := swagger.MarshalJSON()
		w.Write(b)
	})

	r.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))

	fmt.Printf("Starting server on %s:%s...\n", cfg.Server.Host, cfg.Server.Port)
	http.ListenAndServe(":8080", r)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
