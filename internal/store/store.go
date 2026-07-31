package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"split-tunnel-manager/internal/domain"
)

// Store holds bypass rules.
type Store struct {
	db *sql.DB
}

// NewStore opens SQLite at path and migrates schema.
func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS bypass_rules (
			id TEXT PRIMARY KEY,
			cidr TEXT NOT NULL,
			label TEXT NOT NULL,
			enabled INTEGER NOT NULL DEFAULT 1,
			source TEXT NOT NULL DEFAULT 'user'
		)`,
	}
	for _, q := range stmts {
		if _, err := s.db.Exec(q); err != nil {
			return fmt.Errorf("store migrate: %w", err)
		}
	}
	return nil
}

// Close closes the DB.
func (s *Store) Close() error {
	return s.db.Close()
}

// UpsertBypass inserts or replaces a bypass rule.
func (s *Store) UpsertBypass(r domain.BypassRule) error {
	en := 0
	if r.Enabled {
		en = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO bypass_rules (id, cidr, label, enabled, source) VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET cidr=excluded.cidr, label=excluded.label,
		 enabled=excluded.enabled, source=excluded.source`,
		r.ID, r.CIDR, r.Label, en, r.Source,
	)
	return err
}

// ListBypass returns all bypass rules.
func (s *Store) ListBypass() ([]domain.BypassRule, error) {
	rows, err := s.db.Query(`SELECT id, cidr, label, enabled, source FROM bypass_rules ORDER BY cidr`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.BypassRule
	for rows.Next() {
		var r domain.BypassRule
		var en int
		if err := rows.Scan(&r.ID, &r.CIDR, &r.Label, &en, &r.Source); err != nil {
			return nil, err
		}
		r.Enabled = en != 0
		out = append(out, r)
	}
	return out, rows.Err()
}

// DeleteBypass removes a bypass rule by id.
func (s *Store) DeleteBypass(id string) error {
	_, err := s.db.Exec(`DELETE FROM bypass_rules WHERE id = ?`, id)
	return err
}
