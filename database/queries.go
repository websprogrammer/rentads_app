package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"rentads_app/database/sql"
	"rentads_app/schemas"
)

var PageSize int8 = 20

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
			"page_size":      PageSize,
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
