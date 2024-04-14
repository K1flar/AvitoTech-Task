package integration

import (
	"banner_service/cmd/server"
	"banner_service/internal/config"
	"banner_service/internal/domains"
	"banner_service/internal/logger"
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type writer struct{}

func (w *writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var (
	cfg = &config.Config{
		Server: config.Server{
			Port:           "8081",
			ResponseTime:   50 * time.Millisecond,
			BanneLifeCycle: time.Minute,
			RPS:            1000,
			UserToken:      "user",
			AdminToken:     "admin",
		},
		Database: config.Database{
			DSN: "postgres://postgres:postgres@localhost:5432/test?sslmode=disable",
		},
	}
	app *testApp
)

func TestMain(m *testing.M) {
	a, err := server.New(cfg, logger.New(&writer{}))
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	app = &testApp{a, db}

	code := m.Run()
	app.mustClearDB()
	os.Exit(code)
}

func TestUnauthorized(t *testing.T) {
	req := httptest.NewRequest("GET", "/user_banner", nil)

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
}

func TestForbidden(t *testing.T) {
	req := httptest.NewRequest("GET", "/user_banner", nil)
	req.Header.Add("token", "asdasdf")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Result().StatusCode)
}

func TestInvalidFeatureID(t *testing.T) {
	req := httptest.NewRequest("GET", "/user_banner?feature_id=asdas", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, []byte(`{"error":"invalid feature id"}`), body)
}

func TestInvalidTagID(t *testing.T) {
	req := httptest.NewRequest("GET", "/user_banner?feature_id=123&tag_id=asg", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, []byte(`{"error":"invalid tag id"}`), body)
}

func TestInvalidUseLastRevision(t *testing.T) {
	req := httptest.NewRequest("GET", "/user_banner?feature_id=123&tag_id=123&use_last_revision=bob", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, []byte(`{"error":"use_last_revision parameter must be true or false"}`), body)
}

func TestEmpty(t *testing.T) {
	app.mustClearDB()

	req := httptest.NewRequest("GET", "/user_banner?feature_id=1&tag_id=1", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestCorrect(t *testing.T) {
	app.mustClearDB()
	banner := &domains.Banner{
		ID:        1,
		Content:   `{"text": "bebr"}`,
		IsActive:  true,
		FeatureID: 1000,
	}
	app.mustAddBanner(banner, []int{1, 2, 3})

	req := httptest.NewRequest("GET", "/user_banner?feature_id=1000&tag_id=1", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, banner.Content, domains.Content(body))
}

func TestFromCache(t *testing.T) {
	app.mustClearDB()
	banner := &domains.Banner{
		Content: `{"text": "bebr"}`,
	}

	req := httptest.NewRequest("GET", "/user_banner?feature_id=1000&tag_id=1", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, banner.Content, domains.Content(body))
}

func TestUseLastRevision(t *testing.T) {
	app.mustClearDB()
	banner := &domains.Banner{
		ID:        1,
		Content:   `{"text": "aboba"}`,
		IsActive:  true,
		FeatureID: 1000,
	}
	app.mustAddBanner(banner, []int{1, 2, 3})

	req := httptest.NewRequest("GET", "/user_banner?feature_id=1000&tag_id=1&use_last_revision=true", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, banner.Content, domains.Content(body))
}

func TestNotActive(t *testing.T) {
	app.mustClearDB()
	banner := &domains.Banner{
		ID:        10,
		Content:   `{"text": "bebr"}`,
		FeatureID: 10,
	}
	app.mustAddBanner(banner, []int{1, 2, 3})

	req := httptest.NewRequest("GET", "/user_banner?feature_id=10&tag_id=1", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestTimeLimit(t *testing.T) {
	cfg.Server.ResponseTime = time.Microsecond
	a, _ := server.New(cfg, logger.New(&writer{}))
	app.App = a
	req := httptest.NewRequest("GET", "/user_banner?feature_id=1&tag_id=1", nil)
	req.Header.Add("token", "user")

	rr := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusGatewayTimeout, rr.Result().StatusCode)
}
