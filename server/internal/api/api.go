package api

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nobleach/cobaltcms/internal/storage"
	"github.com/rs/zerolog"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
	logger     zerolog.Logger
}

type Health struct {
	Status string `json:"status"`
}

type Content struct {
	Id                 string      `json:"id"`
	ContentType        string      `json:"contentType"`
	Name               string      `json:"name"`
	Body               interface{} `json:"body"`
	ExtendedAttributes interface{} `json:"extendedAttributes"`
	PublishedStatus    string      `json:"publishedStatus"`
	PublishStart       time.Time   `json:"publishStart"`
	PublishEnd         time.Time   `json:"publishEnd"`
}

func NewApiServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		logger:     zerolog.New(os.Stderr),
	}
}

func (s *APIServer) Run() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/health", s.handleGetHealthcheck)
	e.GET("/published-statuses", s.handleGetPublishedStatuses)

	e.Logger.Fatal(e.Start(s.listenAddr))
}

func (s *APIServer) handleGetHealthcheck(c echo.Context) error {
	health := &Health{
		Status: "UP",
	}
	return c.JSON(http.StatusOK, health)
}

func (s *APIServer) handleGetPublishedStatuses(c echo.Context) error {
	res, err := s.store.GetPublishedStatuses()

	if err != nil {
		s.logger.Error().Msg("Could not fetch statuses")
	}

	return c.JSON(http.StatusOK, res)
}

type InvalidInputError struct {
	Message string `json:"message"`
}
