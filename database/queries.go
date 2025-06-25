package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"rentads_app/config"
	"rentads_app/database/sql"
	"rentads_app/schemas"
)

func (db *DB) GetAdvertsFromDB(
	ctx context.Context,
	olderThan uint64,
	city string,
	rentType int,
	roomType int,
	districtsList []string,
	subDistrictsList []string,
	metroStations []string,
	keyWords string,
) ([]schemas.Advert, error) {
	conn, err := db.Client.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, err := sql.GetSQLFromFile("get_adverts.sql")
	rows, err := conn.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"older_than":     olderThan,
			"page_size":      config.AdvertsPageSize,
			"city":           city,
			"rent_type":      rentType,
			"room_type":      roomType,
			"districts":      districtsList,
			"sub_districts":  subDistrictsList,
			"metro_stations": metroStations,
			"key_words":      keyWords,
		},
	)
	if err != nil {
		slog.Error("Failed to get adverts from DB", err)
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[schemas.Advert])
}

func (db *DB) UpsertNotificationToDB(
	ctx context.Context,
	notification *schemas.Notification,
) error {
	conn, err := db.Client.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query, err := sql.GetSQLFromFile("add_notification.sql")
	_, err = conn.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"device_token": notification.DeviceToken,
			"city":         notification.City,
			"rent_type":    notification.RentType,
			"room_type":    notification.RoomType,
			"districts":    notification.Districts,
			"key_words":    notification.KeyWords,
			"enabled":      notification.Enabled,
		},
	)
	if err != nil {
		slog.Error("Failed to upsert notification to DB", err)
		return err
	}
	return nil
}
