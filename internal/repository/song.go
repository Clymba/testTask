package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"time"
)

type DataBaseMethods interface {
	CreateSong(ctx context.Context, song *Song) error
	GetSongByID(ctx context.Context, song *Song) error
	GetAllSongs(ctx context.Context) ([]Song, error)
	UpdateSong(ctx context.Context, song *Song) error
	DeleteSong(ctx context.Context, song *Song) error
	GetFilteredSongs(ctx context.Context, groupName, text, genre string) ([]Song, error)
}

type Song struct {
	ID         uuid.UUID `json:"id"`
	GroupName  string    `json:"groupName"`
	Text       string    `json:"text"`
	Genre      string    `json:"genre"`
	Date_added time.Time `json:"date_added"`
	Link       string    `json:"link"`
}

func (s *Repository) CreateSong(ctx context.Context, song *Song) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `INSERT INTO songs (group_name, text, genre, date_added, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if err := s.pool.QueryRow(ctx, query, song.GroupName, song.Text, song.Genre, song.Date_added, song.Link).Scan(&song.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, State: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			fmt.Println(newErr)
			return nil
		}
		fmt.Println("Ошибка при сохранении песни:", err)
		return err
	}

	return nil
}

func (s *Repository) GetSongByID(ctx context.Context, song *Song) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `SELECT group_name, text, genre, date_added, link 
              FROM songs WHERE id = $1`

	if err := s.pool.QueryRow(ctx, query, song.ID).Scan(&song.GroupName, &song.Text, &song.Genre, &song.Date_added, &song.Link); err != nil {
		return err
	}

	return nil
}

func (s *Repository) GetAllSongs(ctx context.Context) ([]Song, error) {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	query := `SELECT id, group_name, text, genre, date_added, link 
              FROM songs`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.Text, &song.Genre, &song.Date_added, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *Repository) GetFilteredSongs(ctx context.Context, groupName, text, genre, link string, page, limit int) ([]Song, error) {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	query := "SELECT id, group_name, text, genre, date_added, link FROM songs"
	args := []interface{}{}
	argCount := 1
	conditions := []string{}

	if groupName != "" {
		conditions = append(conditions, fmt.Sprintf("group_name ILIKE $%d", argCount)) // Исправлено на group_name
		args = append(args, "%"+groupName+"%")
		argCount++
	}

	if text != "" {
		conditions = append(conditions, fmt.Sprintf("text ILIKE $%d", argCount))
		args = append(args, "%"+text+"%")
		argCount++
	}

	if genre != "" {
		conditions = append(conditions, fmt.Sprintf("genre ILIKE $%d", argCount))
		args = append(args, "%"+genre+"%")
		argCount++
	}

	if link != "" {
		conditions = append(conditions, fmt.Sprintf("link ILIKE $%d", argCount)) // Исправлено на link
		args = append(args, "%"+link+"%")
		argCount++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.Text, &song.Genre, &song.Date_added, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *Repository) UpdateSong(ctx context.Context, song *Song) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	query := `UPDATE songs 
			  SET group_name = $1, text = $2, genre = $3, date_added = $4 
			  WHERE id = $5`

	if _, err := s.pool.Exec(ctx, query, song.GroupName, song.Text, song.Genre, song.Date_added, song.ID); err != nil {
		return err
	}

	return nil
}

func (s *Repository) DeleteSong(ctx context.Context, song *Song) error {
	ctx, timeout := context.WithTimeout(ctx, s.config.Timeout)
	defer timeout()

	query := `DELETE from songs 
       		  WHERE id = $1`

	_, err := s.pool.Exec(ctx, query, &song.ID)

	return err
}
