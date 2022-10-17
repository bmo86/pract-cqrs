package database

import (
	"context"
	"cqrs/models"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PostgresReposytory struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresReposytory, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &PostgresReposytory{db}, nil
}

func (repo *PostgresReposytory) Close() {
	repo.db.Close()
}

func (repo *PostgresReposytory) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO feeds (id, tittle, description) VALUES ($1, $2, $3)", feed.ID, feed.Title, feed.Description)
	return err
}

func (repo *PostgresReposytory) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, title, description, created_at FROM feeds")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	feeds := []*models.Feed{}

	for rows.Next() {
		feed := &models.Feed{}
		if err := rows.Scan(&feed.ID, &feed.Title, &feed.Description, &feed.CreatedAt); err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	return feeds, nil
}
