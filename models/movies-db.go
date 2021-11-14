package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type DBModel struct{
	DB *sql.DB
}

//for one movie
func (m *DBModel) GetOneMovie(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
	defer cancel()
	query := `select id, title, description, year, release_date, rating, runtime, mpaa_rating,
				created_at, updated_at from movies where id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	var movie Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	//get genres
	query = `select
				mg.id, mg.movie_id, mg.genre_id, g.genre_name
			from
				movies_genres mg
			left join genres g
			on (g.id = mg.genre_id) 
			where mg.movie_id = $1
			`
	//var genres []MovieGenre
	genres := make(map[int]string)
	rows,_ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
			)
		if err != nil {
			return nil, err
		}
		//genres = append(genres, mg)
		genres[mg.ID] = mg.Genre.GenreName
	}
	movie.MovieGenre = genres

	return &movie, nil
}

//for all movies
func (m *DBModel) AllMovies() (*[]Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
	defer cancel()
	query := `select id, title, description, year, release_date, rating, runtime, mpaa_rating,
				created_at, updated_at from movies
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var movies_list []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies_list = append(movies_list, movie)
	}
	return &movies_list, nil
}

func (m *DBModel) AllGenres() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
	defer cancel()
	query := `select id, genre_name from genres`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var genre_list []*Genre
	for rows.Next() {
		var genre Genre
		err = rows.Scan(
			&genre.ID,
			&genre.GenreName,
		)
		if err != nil {
			return nil, err
		}
		genre_list = append(genre_list, &genre)
	}
	return genre_list, nil
}

func (m *DBModel) AllMoviesByGenre(id int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
	defer cancel()
	query := `SELECT m.id, m.title
		FROM movies_genres mg INNER JOIN movies m 
		ON mg.movie_id = m.id
		WHERE mg.genre_id = $1;
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	var movie_list []*Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(
			&movie.ID,
			&movie.Title,
		)
		if err != nil {
			return nil, err
		}
		movie_list = append(movie_list, &movie)
	}
	return movie_list, nil
}

func (m *DBModel) InsertMovie(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
	defer cancel()

	var count int
	query := `SELECT MAX(id) FROM movies;`
	if err := m.DB.QueryRow(query).Scan(&count); err != nil{
		log.Println(err)
		return err
	}


	query = `INSERT INTO movies(id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9, $10);`

	_, err := m.DB.ExecContext(ctx, query, count+1, movie.Title, movie.Description, movie.Year,
		movie.ReleaseDate, movie.Runtime, movie.Rating, movie.MPAARating, movie.CreatedAt, movie.UpdatedAt)

	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("SUCCESS")
	return nil
}

func (m *DBModel) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE
		  FROM movies_genres 
		  WHERE movie_id = $1;`
	_, err := m.DB.ExecContext(ctx, query, id)


	if err != nil {
		log.Println(err)
		return err
	}
	query = `DELETE
		  FROM movies 
		  WHERE id = $1;`
	_, err = m.DB.ExecContext(ctx, query, id)


	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("SUCCESS")
	return nil
}