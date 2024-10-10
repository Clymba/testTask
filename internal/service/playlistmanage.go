package service

import (
	"context"
	"errors"
	"github.com/Clymba/testTask/internal/repository"
	"github.com/Clymba/testTask/logger"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("song not found")
)

type SongService interface {
	CreateSong(ctx context.Context, song *repository.Song) error
	GetSongByID(ctx context.Context, id uuid.UUID) (*repository.Song, error)
	GetAllSongs(ctx context.Context) ([]repository.Song, error)
	UpdateSong(ctx context.Context, song *repository.Song) error
	DeleteSong(ctx context.Context, id uuid.UUID) error
	GetFilteredSongs(ctx context.Context, groupName, text, genre string, page, limit int) ([]repository.Song, error)
}

func (s *Service) CreateSong(ctx context.Context, song *repository.Song) error {
	if song == nil {
		return errors.New("bad id")
	}
	song.ID = uuid.New()
	return s.repository.CreateSong(ctx, song)
}

func (s *Service) GetSongByID(ctx context.Context, id uuid.UUID) (*repository.Song, error) {
	if id == uuid.Nil {
		return nil, errors.New("bad id")
	}

	song := &repository.Song{ID: id}
	if err := s.repository.GetSongByID(ctx, song); err != nil {
		return nil, err
	}
	return song, nil
}

func (s *Service) GetAllSongs(ctx context.Context) ([]repository.Song, error) {
	return s.repository.GetAllSongs(ctx)
}

func (s *Service) GetFilteredSongs(ctx context.Context, groupName, text, genre, link string, page, limit int) ([]repository.Song, error) {
	logger.Log.Infof("filter parameters: groupName=%s, text=%s, genre=%s, page=%d, limit=%d", groupName, text, genre, page, limit)

	songs, err := s.repository.GetFilteredOrPaginationSongs(ctx, groupName, text, genre, link, page, limit)
	if err != nil {
		logger.Log.Errorf("Error filtere songs: %v", err)
		return nil, err
	}

	logger.Log.Infof("Fetched %d songs", len(songs))

	return songs, nil
}

func (s *Service) UpdateSong(ctx context.Context, song *repository.Song) error {
	if song == nil {
		return errors.New("empty song")
	}
	return s.repository.UpdateSong(ctx, song)
}

func (s *Service) DeleteSong(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("bad id")
	}
	song := &repository.Song{ID: id}
	return s.repository.DeleteSong(ctx, song)
}
