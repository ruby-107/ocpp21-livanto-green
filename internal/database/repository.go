package database

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) SaveStation(
	ctx context.Context,
	id, model, vendor string,
) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO charging_stations
			(station_id, model, vendor_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (station_id)
		DO UPDATE SET
			model = EXCLUDED.model,
			vendor_name = EXCLUDED.vendor_name
	`, id, model, vendor)

	if err != nil {
		return err
	}

	return r.SaveLog(ctx, id, "BootNotification", "Accepted", nil)
}

func (r *Repository) SaveHeartbeat(
	ctx context.Context,
	id string,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE charging_stations
		SET last_heartbeat = $1
		WHERE station_id = $2
	`, time.Now().UTC(), id)

	if err != nil {
		return err
	}

	return r.SaveLog(ctx, id, "Heartbeat", "Received", nil)
}

func (r *Repository) SaveStatus(
	ctx context.Context,
	id, status string,
	connectorID int,
) error {
	return r.SaveLog(
		ctx, id, "StatusNotification", status, &connectorID,
	)
}

func (r *Repository) SaveLog(
	ctx context.Context,
	id, messageType, status string,
	connectorID *int,
) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO station_logs
			(station_id, message_type, status, connector_id)
		VALUES ($1, $2, $3, $4)
	`, id, messageType, status, connectorID)

	return err
}
