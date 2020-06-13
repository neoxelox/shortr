package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Repo describes the URLs repository
type Repo struct {
	db *pgxpool.Pool
}

// Connect tries to connect to an specified database via the dsn connection string
func Connect(dsn string, retries int) (*Repo, error) {
	db, err := connect(context.Background(), dsn, retries)
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

// GetByID retrieves the url by its id
func (r *Repo) GetByID(id int) (string, error) {
	var url string
	err := r.db.QueryRow(context.Background(), `SELECT "url" FROM "urls"
												WHERE "id" = $1`, id).Scan(&url)
	return url, err
}

// GetByName retrieves the url by its name
func (r *Repo) GetByName(name string) (string, error) {
	var url string
	err := r.db.QueryRow(context.Background(), `SELECT "url" FROM "urls"
												WHERE "name" = $1`, name).Scan(&url)
	return url, err
}

// GetMetricsByID retrieves the url metrics by its id
func (r *Repo) GetMetricsByID(id int) (int, *time.Time, time.Time, time.Time, error) {
	var hits int
	var lastHitAt *time.Time
	var createdAt time.Time
	var modifiedAt time.Time
	err := r.db.QueryRow(context.Background(), `SELECT "hits", "last_hit_at", "created_at", "modified_at" FROM "urls"
												WHERE "name" = $1`, id).Scan(&hits, &lastHitAt, &createdAt, &modifiedAt)
	return hits, lastHitAt, createdAt, modifiedAt, err
}

// GetMetricsByName retrieves the url metrics by its name
func (r *Repo) GetMetricsByName(name string) (int, *time.Time, time.Time, time.Time, error) {
	var hits int
	var lastHitAt *time.Time
	var createdAt time.Time
	var modifiedAt time.Time
	err := r.db.QueryRow(context.Background(), `SELECT "hits", "last_hit_at", "created_at", "modified_at" FROM "urls"
												WHERE "name" = $1`, name).Scan(&hits, &lastHitAt, &createdAt, &modifiedAt)
	return hits, lastHitAt, createdAt, modifiedAt, err
}

// Create creates a new entry for the url and returns its newly created id
func (r *Repo) Create(url string) (int, error) {
	var id int
	createdAt := time.Now()
	modifiedAt := createdAt
	err := r.db.QueryRow(context.Background(), `INSERT INTO "urls" ("url", "created_at", "modified_at")
											    VALUES ($1, $2, $3) RETURNING "id"`, url, createdAt, modifiedAt).Scan(&id)
	return id, err
}

// UpdateNameByID updates the name for the url by its id
func (r *Repo) UpdateNameByID(id int, name string) error {
	modifiedAt := time.Now()
	ret, err := r.db.Exec(context.Background(), `UPDATE "urls"
											     SET "name" = $1, "modified_at" = $2
											     WHERE "id" = $3`, name, modifiedAt, id)
	if ret.RowsAffected() == 0 {
		err = errors.New("Update query must affect rows")
	}
	return err
}

// UpdateURLByID updates the url for an entry by its id
func (r *Repo) UpdateURLByID(id int, url string) error {
	modifiedAt := time.Now()
	ret, err := r.db.Exec(context.Background(), `UPDATE "urls"
											     SET "url" = $1, "modified_at" = $2
											     WHERE "id" = $3`, url, modifiedAt, id)
	if ret.RowsAffected() == 0 {
		err = errors.New("Update query must affect rows")
	}
	return err
}

// UpdateURLByName updates the url for an entry by its name
func (r *Repo) UpdateURLByName(name string, url string) error {
	modifiedAt := time.Now()
	ret, err := r.db.Exec(context.Background(), `UPDATE "urls"
											     SET "url" = $1, "modified_at" = $2
											     WHERE "name" = $3`, url, modifiedAt, name)
	if ret.RowsAffected() == 0 {
		err = errors.New("Update query must affect rows")
	}
	return err
}

// UpdateMetricsByID updates the metrics for the url by its id
func (r *Repo) UpdateMetricsByID(id int) error {
	lastHitAt := time.Now()
	ret, err := r.db.Exec(context.Background(), `UPDATE "urls"
											     SET "hits" = "hits" + 1, "last_hit_at" = $1
											     WHERE "id" = $2`, lastHitAt, id)
	if ret.RowsAffected() == 0 {
		err = errors.New("Update query must affect rows")
	}
	return err
}

// UpdateMetricsByName updates the metrics for the url by its name
func (r *Repo) UpdateMetricsByName(name string) error {
	lastHitAt := time.Now()
	ret, err := r.db.Exec(context.Background(), `UPDATE "urls"
											     SET "hits" = "hits" + 1, "last_hit_at" = $1
											     WHERE "name" = $2`, lastHitAt, name)
	if ret.RowsAffected() == 0 {
		err = errors.New("Update query must affect rows")
	}
	return err
}

// DeleteByID deletes de url entry by its id
func (r *Repo) DeleteByID(id int) error {
	ret, err := r.db.Exec(context.Background(), `DELETE FROM "urls"
											     WHERE "id" = $1`, id)
	if ret.RowsAffected() == 0 {
		err = errors.New("Delete query must affect rows")
	}
	return err
}

// DeleteByName deletes de url entry by its name
func (r *Repo) DeleteByName(name string) error {
	ret, err := r.db.Exec(context.Background(), `DELETE FROM "urls"
											     WHERE "name" = $1`, name)
	if ret.RowsAffected() == 0 {
		err = errors.New("Delete query must affect rows")
	}
	return err
}
