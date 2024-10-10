package server

import (
	"context"
	"errors"
	"github.com/Clymba/testTask/internal/repository"
	"github.com/Clymba/testTask/internal/service"
	"github.com/Clymba/testTask/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrInvalidRequest = fiber.NewError(fiber.StatusBadRequest, "Invalid request: unable to parse body")
	ErrInvalidID      = fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	ErrSongNotFound   = fiber.NewError(fiber.StatusNotFound, "Song not found")
)

type Server struct {
	app     *fiber.App
	service *service.Service
}

func New(ctx context.Context, app *fiber.App, service *service.Service) *Server {
	srv := &Server{
		app:     app,
		service: service,
	}
	srv.registerHandlers(ctx)
	return srv
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) registerHandlers(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	api := s.app.Group("/api")

	api.Post("/add_song", s.addSong)
	api.Get("/find_song/:id", s.findSong)
	api.Get("/view_song", s.viewSongs)
	api.Get("/view_song_with_filter", s.viewSongsWithFilter)
	api.Put("/update_song/:id", s.updateSong)
	api.Delete("/delete_song/:id", s.deleteSong)
}

func (s *Server) addSong(c *fiber.Ctx) error {
	var song repository.Song
	if err := c.BodyParser(&song); err != nil {
		logger.Log.Errorf("Ошибка парсинга: %v", err)
		return s.handleError("Ошибка парсинга реквеста", err, fiber.StatusBadRequest)
	}

	if err := s.service.CreateSong(c.Context(), &song); err != nil {
		logger.Log.Errorf("Не создалась песня: %v", err)
		return s.handleError("Не создалась песня", err, fiber.StatusInternalServerError)
	}

	logger.Log.Infof("Песня создалась: %+v", song)
	return c.Status(fiber.StatusCreated).JSON(song)
}

func (s *Server) findSong(c *fiber.Ctx) error {
	id, err := s.parseUUID(c.Params("id"))
	if err != nil {
		logger.Log.Warnf("Неправильный формат ID: %v", err)
		return ErrInvalidID
	}

	song, err := s.service.GetSongByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			logger.Log.Warnf("Песня не найдена: ID=%s", id)
			return ErrSongNotFound
		}
		logger.Log.Errorf("Не нашлась песня: %v", err)
		return s.handleError("Не нашлась песня", err, fiber.StatusInternalServerError)
	}

	logger.Log.Infof("Песня: %+v", song)
	return c.JSON(song)
}

func (s *Server) viewSongs(c *fiber.Ctx) error {
	songs, err := s.service.GetAllSongs(c.Context())
	if err != nil {
		logger.Log.Errorf("Не нашлись песни: %v", err)
		return s.handleError("Не нашлись песни", err, fiber.StatusInternalServerError)
	}

	logger.Log.Infof("Песни: %d", len(songs))
	return c.JSON(songs)
}

func (s *Server) viewSongsWithFilter(c *fiber.Ctx) error {
	groupName := c.Query("group_name")
	text := c.Query("text")
	genre := c.Query("genre")
	link := c.Query("link") // Добавляем извлечение параметра link

	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	if groupName == "" && text == "" && genre == "" && link == "" {
		return s.viewSongs(c)
	}

	songs, err := s.service.GetFilteredSongs(c.Context(), groupName, text, genre, link, page, limit)
	if err != nil {
		logger.Log.Errorf("Ошибка фильтрации: %v", err)
		return s.handleError("Ошибка фильтрации", err, fiber.StatusInternalServerError)
	}

	logger.Log.Infof("Обнаруженные %d песни: %+v", len(songs), songs)
	return c.JSON(songs)
}

func (s *Server) updateSong(c *fiber.Ctx) error {
	id, err := s.parseUUID(c.Params("id"))
	if err != nil {
		logger.Log.Warnf("Неправильный id: %v", err)
		return ErrInvalidID
	}

	var song repository.Song
	if err := c.BodyParser(&song); err != nil {
		logger.Log.Errorf("Ошибка парсинга: %v", err)
		return s.handleError("Ошибка парсинга", err, fiber.StatusBadRequest)
	}

	song.ID = id
	if err := s.service.UpdateSong(c.Context(), &song); err != nil {
		logger.Log.Errorf("Не обновилась песня: %v", err)
		return s.handleError("Не обновилась песня", err, fiber.StatusInternalServerError)
	}

	logger.Log.Infof("Успешно обновилась: %+v", song)
	return c.JSON(song)
}

func (s *Server) deleteSong(c *fiber.Ctx) error {
	id, err := s.parseUUID(c.Params("id"))
	if err != nil {
		logger.Log.Warnf("Не тот формат: %v", err)
		return ErrInvalidID
	}

	if err := s.service.DeleteSong(c.Context(), id); err != nil {
		logger.Log.Errorf("Не удалось удалить: %v", err)
		return s.handleError("Не удалось удалить", err, fiber.StatusInternalServerError)
	}

	logger.Log.Infof("Удалось удалить: ID=%s", id)
	return c.SendStatus(http.StatusNoContent)
}

func (s *Server) parseUUID(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Не тот формат: %v", err)
	}
	return id, err
}

func (s *Server) handleError(message string, err error, status int) error {
	log.Printf("%s: %v", message, err)
	return fiber.NewError(status, message+": "+err.Error())
}
