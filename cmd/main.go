package main

import (
	"banner_service/api"
	"banner_service/internal/config"
	"banner_service/internal/repositories/postgres"
	"banner_service/pkg/mux"
	"fmt"
	"os"

	// _ "banner_service/docs"

	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct{}

func (s *Server) GetBanner(w http.ResponseWriter, r *http.Request, params api.GetBannerParams) {
	w.Write([]byte("GetBanner"))
}

// Создание нового баннера
// (POST /banner)
func (s *Server) PostBanner(w http.ResponseWriter, r *http.Request, params api.PostBannerParams) {
	w.Write([]byte("PostBanner"))
}

// Удаление баннера по идентификатору
// (DELETE /banner/{id})
func (s *Server) DeleteBannerId(w http.ResponseWriter, r *http.Request, id int, params api.DeleteBannerIdParams) {
	w.Write([]byte("DeleteBannerId"))
}

// Обновление содержимого баннера
// (PATCH /banner/{id})
func (s *Server) PatchBannerId(w http.ResponseWriter, r *http.Request, id int, params api.PatchBannerIdParams) {
	w.Write([]byte("PatchBannerId"))
}

// Получение баннера для пользователя
// (GET /user_banner)
func (s *Server) GetUserBanner(w http.ResponseWriter, r *http.Request, params api.GetUserBannerParams) {
	w.Write([]byte("GetUserBanner"))
}

const configPath = "configs/config.yaml"

func main() {
	cfg := config.New(configPath)

	repository, err := postgres.New(&cfg.Database)
	_ = repository
	exitOnError(err)

	swagger, err := api.GetSwagger()
	exitOnError(err)
	swagger.Servers = nil

	r := mux.New()
	api.HandlerFromMux(&Server{}, r.GetMux())

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		b, _ := swagger.MarshalJSON()
		w.Write(b)
	})

	r.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))

	http.ListenAndServe(":8080", r)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
