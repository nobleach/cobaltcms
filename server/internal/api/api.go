package api

import (
	"net/http"
	"os"

	"github.com/knadh/koanf"
	"github.com/labstack/echo/v4"
	"github.com/nobleach/cobaltcms/internal/storage"
	"github.com/nobleach/cobaltcms/internal/types"
	"github.com/rs/zerolog"
)

type APIServer struct {
	config koanf.Koanf
	store  storage.Storage
	logger zerolog.Logger
}

type Health struct {
	Status string `json:"status"`
}

type PublishedContent struct {
	Id                 string      `json:"id"`
	FragmentType       string      `json:"fragmentType"`
	Name               string      `json:"name"`
	Body               types.JSONB `json:"body"`
	ExtendedAttributes types.JSONB `json:"extendedAttributes"`
}

func NewApiServer(config *koanf.Koanf, store storage.Storage) *APIServer {
	return &APIServer{
		config: *config,
		store:  store,
		logger: zerolog.New(os.Stdout),
	}
}

func (s *APIServer) Run() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/health", s.handleGetHealthcheck)
	e.GET("/published-statuses", s.handleGetPublishedStatuses)
	e.GET("/published-page", s.handleGetPublishedForId)
	e.GET("/published-content", s.handleGetPublishedForDate)
	e.POST("/content", s.handlePostContent)

	port := s.config.String("server.port")

	e.Logger.Fatal(e.Start(port))
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

func (s *APIServer) handleGetPublishedForDate(c echo.Context) error {
	date := c.QueryParam("date")
	time := c.QueryParam("time")
	res, err := s.store.GetPublishedContentForDate(date + " " + time)

	if err != nil {
		s.logger.Error().Msg("Could not fetch content")
	}

	var publishedContents []PublishedContent

	for _, element := range res {
		content := PublishedContent{
			Id:                 element.ID.String(),
			FragmentType:       element.FragmentType,
			Name:               element.Name,
			Body:               element.Body,
			ExtendedAttributes: element.ExtendedAttributes,
		}

		publishedContents = append(publishedContents, content)
	}

	return c.JSON(http.StatusOK, publishedContents)
}

type PublishedPageContent struct {
	Id      string         `json:"id"`
	Content map[string]any `json:"content"`
}

func (s *APIServer) handleGetPublishedForId(c echo.Context) error {
	// TODO: We need to get this from the client as it'll likely
	// be in another timezone
	date := c.QueryParam("date")
	time := c.QueryParam("time")
	pageId := c.QueryParam("id")
	res, err := s.store.GetPublishedContentForId(pageId, date+" "+time)

	if err != nil {
		s.logger.Error().Msg("Could not fetch content")
	}

	publishedContents := PublishedPageContent{
		Id: pageId,
	}

	contentMap := make(map[string]any)

	for _, element := range res {
		contentMap[element.Name.String] = element.Body
	}

	publishedContents.Content = contentMap

	return c.JSON(http.StatusOK, publishedContents)
}

func (s *APIServer) handlePostContent(c echo.Context) error {
	var content types.NewContent
	if err := c.Bind(&content); err != nil {
		s.logger.Error().Msg("Invalid input")
		return c.JSON(http.StatusBadRequest, InvalidInputError{Message: "Invalid input"})
	}

	// TODO: Handle err
	createdUuid, err := s.store.SaveContent(content)
	if err != nil {
		type Error struct {
			Error string `json:"error"`
		}

		newError := Error{
			Error: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, newError)
	}

	type NewUuid struct {
		Uuid string `json:"uuid"`
	}

	newUuid := NewUuid{
		Uuid: createdUuid,
	}

	return c.JSON(http.StatusOK, newUuid)
}
