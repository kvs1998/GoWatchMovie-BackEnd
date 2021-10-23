package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct{
	DB *sql.DB
}

//for one movie
func (m *DBModel) Get(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
	defer cancel()
	query := `SELECT id, title, description, year, release_date FROM movies WHERE id=$1;`
	row := m.DB.QueryRowContext(ctx, query, id)
	var movie Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Year,
		&movie.ReleaseDate,
		)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

//for all movies
func (m *DBModel) All(id int) (*Movie, error) {
	return nil, nil
}