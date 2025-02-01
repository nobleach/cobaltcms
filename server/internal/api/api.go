package api

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type APIServer struct {
	listenAddr string
	// store      storage.Storage
	logger zerolog.Logger
}

type Health struct {
	Status string `json:"status"`
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		// store:      store,
		logger: zerolog.New(os.Stderr),
	}
}

func (s *APIServer) Run() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/health", s.handleGetHealthcheck)

	e.Logger.Fatal(e.Start(s.listenAddr))
}

func (s *APIServer) handleGetHealthcheck(c echo.Context) error {
	health := &Health{
		Status: "UP",
	}
	return c.JSON(http.StatusOK, health)
}

type InvalidInputError struct {
	Message string `json:"message"`
}
