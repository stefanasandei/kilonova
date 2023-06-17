package db

import (
	"context"
	"errors"

	"github.com/KiloProjects/kilonova"
	"github.com/jackc/pgx/v5"
)

// Functions that interact with views and functions

func (s *DB) IsProblemViewer(ctx context.Context, problemID, userID int) (bool, error) {
	var cnt int
	err := s.pgconn.QueryRow(ctx, "SELECT COUNT(*) FROM visible_pbs($1) WHERE problem_id = $2", userID, problemID).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *DB) IsFullProblemViewer(ctx context.Context, problemID, userID int) (bool, error) {
	var cnt int
	err := s.pgconn.QueryRow(ctx, "SELECT COUNT(*) FROM persistently_visible_pbs($1) WHERE problem_id = $2", userID, problemID).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *DB) IsProblemEditor(ctx context.Context, problemID, userID int) (bool, error) {
	var cnt int
	err := s.pgconn.QueryRow(ctx, "SELECT COUNT(*) FROM problem_editors WHERE problem_id = $1 AND user_id = $2", problemID, userID).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *DB) IsContestViewer(ctx context.Context, contestID, userID int) (bool, error) {
	var cnt int
	err := s.pgconn.QueryRow(ctx, "SELECT COUNT(*) FROM contest_visibility WHERE contest_id = $1 AND user_id = $2", contestID, userID).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *DB) RefreshProblemStats(ctx context.Context) error {
	_, err := s.pgconn.Exec(ctx, "REFRESH MATERIALIZED VIEW problem_statistics")
	return err
}

func (s *DB) ProblemChecklist(ctx context.Context, problemID int) (*kilonova.ProblemChecklist, error) {
	rows, _ := s.pgconn.Query(ctx, "SELECT * FROM problem_checklist WHERE problem_id = $1 LIMIT 1", problemID)
	chk, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[kilonova.ProblemChecklist])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return chk, nil
}
