package model

import (
	"time"
)

// URL describes the URL model
type URL struct {
	ID         int        `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	URL        string     `db:"url" json:"url"`
	Hits       int        `db:"hits" json:"hits"`
	LastHitAt  *time.Time `db:"last_hit_at" json:"last_hit_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ModifiedAt time.Time  `db:"modified_at" json:"modified_at"`
}
