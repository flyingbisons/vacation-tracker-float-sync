package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/flyingbisons/vacation-tracker-float-sync/internal/integrator"
)

type RequestRepositoryDB struct {
	db          *sql.DB
	currentTime time.Time
}

func NewRequestRepositoryDB(db *sql.DB, currentTime time.Time) *RequestRepositoryDB {
	return &RequestRepositoryDB{
		db:          db,
		currentTime: currentTime,
	}
}

func (r *RequestRepositoryDB) CreateRequest(vtRequestID string, floatTimeoffID int64) error {
	_, err := r.db.Exec("INSERT INTO requests (vt_request_id, float_timeoff_id, created) VALUES (?, ?, ?)", vtRequestID, floatTimeoffID, r.currentTime.Unix())
	if err != nil {
		return err
	}
	return nil
}

func (r *RequestRepositoryDB) GetRequest(vtRequestID string) (integrator.Request, error) {
	var request integrator.Request
	q := fmt.Sprintf("SELECT vt_request_id, float_timeoff_id, created FROM requests WHERE vt_request_id = '%s' LIMIT 1", vtRequestID)
	rows, err := r.db.Query(q)
	if err != nil {
		return request, err
	}

	for rows.Next() {
		err := rows.Scan(&request.VtRequestID, &request.FloatTimeoffID, &request.Created)
		if err != nil {
			return request, err
		}
		return request, nil
	}

	return request, integrator.ErrorRequestNotFound
}
