package repo

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"shortr/model"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgxutil"
)

var ErrNoRows = pgx.ErrNoRows
var rErrNoRows = regexp.MustCompile(fmt.Sprintf("^%s$", pgx.ErrNoRows))
var ErrIntegrityViolation = errors.New("integrity constraint violation")
var rErrIntegrityViolation = regexp.MustCompile(fmt.Sprintf("(SQLSTATE %s)", pgerrcode.UniqueViolation))

// Repo describes the URLs repository
type Repo struct {
	db *pgxpool.Pool
}

// Connect tries to connect to an specified database via the dsn connection string
func Connect(dsn string, retries int, logger pgx.Logger) (*Repo, error) {
	db, err := connect(context.Background(), dsn, retries, logger)
	if err != nil {
		return nil, err
	}
	return &Repo{
		db: db,
	}, nil
}

// Disconnect closes the connection with the database
func (r *Repo) Disconnect() {
	r.db.Close()
}

// GetByID retrieves the URL by its id
func (r *Repo) GetByID(id int) (model.URL, error) {
	var URL model.URL
	query := `SELECT * FROM "urls"
			  WHERE "id" = $1;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		id)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// GetByName retrieves the URL by its name
func (r *Repo) GetByName(name string) (model.URL, error) {
	var URL model.URL
	query := `SELECT * FROM "urls"
			  WHERE "name" = $1;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		name)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// Create creates a new entry for the url and returns the new URL
func (r *Repo) Create(url string) (model.URL, error) {
	var URL model.URL
	createdAt := time.Now()
	modifiedAt := createdAt
	query := `INSERT INTO "urls" ("url", "created_at", "modified_at")
			  VALUES ($1, $2, $3)
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		url, createdAt, modifiedAt)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// UpdateNameByID updates the name for the url by its id and returns the updated URL
func (r *Repo) UpdateNameByID(id int, name string) (model.URL, error) {
	var URL model.URL
	modifiedAt := time.Now()
	query := `UPDATE "urls"
			  SET "name" = $1, "modified_at" = $2
			  WHERE "id" = $3
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		name, modifiedAt, id)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// UpdateURLByID updates the url by its id and returns the updated URL
func (r *Repo) UpdateURLByID(id int, url string) (model.URL, error) {
	var URL model.URL
	modifiedAt := time.Now()
	query := `UPDATE "urls"
	          SET "url" = $1, "modified_at" = $2
			  WHERE "id" = $3
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		url, modifiedAt, id)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// UpdateURLByName updates the url by its name and returns the updated URL
func (r *Repo) UpdateURLByName(name string, url string) (model.URL, error) {
	var URL model.URL
	modifiedAt := time.Now()
	query := `UPDATE "urls"
			  SET "url" = $1, "modified_at" = $2
			  WHERE "name" = $3
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		url, modifiedAt, name)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// UpdateMetricsByID updates the metrics for the url by its id and returns the updated URL
func (r *Repo) UpdateMetricsByID(id int) (model.URL, error) {
	var URL model.URL
	lastHitAt := time.Now()
	query := `UPDATE "urls"
			  SET "hits" = "hits" + 1, "last_hit_at" = $1
			  WHERE "id" = $2
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		lastHitAt, id)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// UpdateMetricsByName updates the metrics for the url by its name and returns the updated URL
func (r *Repo) UpdateMetricsByName(name string) (model.URL, error) {
	var URL model.URL
	lastHitAt := time.Now()
	query := `UPDATE "urls"
			  SET "hits" = "hits" + 1, "last_hit_at" = $1
			  WHERE "name" = $2
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		lastHitAt, name)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// DeleteByID deletes de url entry by its id and returns the deleted URL
func (r *Repo) DeleteByID(id int) (model.URL, error) {
	var URL model.URL
	query := `DELETE FROM "urls"
			  WHERE "id" = $1
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		id)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// DeleteByName deletes de url entry by its name and returns the deleted URL
func (r *Repo) DeleteByName(name string) (model.URL, error) {
	var URL model.URL
	query := `DELETE FROM "urls"
			  WHERE "name" = $1
			  RETURNING *;`
	err := pgxutil.SelectStruct(context.Background(), r.db, &URL, query,
		name)
	if err != nil {
		switch {
		case rErrNoRows.MatchString(err.Error()):
			return URL, ErrNoRows
		case rErrIntegrityViolation.MatchString(err.Error()):
			return URL, ErrIntegrityViolation
		}
	}
	return URL, err
}

// Health checks the database connection health
func (r Repo) Health() error {
	if _, err := r.db.Exec(context.Background(), ";"); err != nil {
		return err
	}
	if r.db.Stat().TotalConns() < r.db.Config().MinConns {
		return errors.New("database connections below the threshold")
	}
	return nil
}
