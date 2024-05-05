package rmu

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"github.com/tkhrk1010/go-samples/actor-model/event-sourcing/test"
)

func TestReadModelDao_InsertWindSpeed(t *testing.T) {
	ctx := context.Background()
	container, err := test.CreateMySQLContainer(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)
	dataSourceName := test.GetDataSourceName(port)

	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	require.NoError(t, err)
	defer func(db *sqlx.DB) {
		if db != nil {
			err := db.Close()
			require.NoError(t, err)
		}
	}(db)

	err = test.MigrateDB(t, err, db, "./../")
	require.NoError(t, err)

	dao := NewWindSpeedDaoImpl(db)
	windSpeedId := "testWindSpeed"
	value := 1.0
	now := time.Now()
	err = dao.InsertWindSpeed(windSpeedId, value, now)
	require.NoError(t, err)

	windSpeed, err := getWindSpeed(db, windSpeedId)
	require.NoError(t, err)
	require.NotNil(t, windSpeed)
	require.Equal(t, windSpeedId, windSpeed["ID"])
	require.Equal(t, value, windSpeed["Value"])
}

func TestReadModelDao_UpdateWindSpeed(t *testing.T) {
	ctx := context.Background()
	container, err := test.CreateMySQLContainer(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)
	dataSourceName := test.GetDataSourceName(port)

	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	require.NoError(t, err)
	defer func(db *sqlx.DB) {
		if db != nil {
			err := db.Close()
			require.NoError(t, err)
		}
	}(db)

	err = test.MigrateDB(t, err, db, "./../")
	require.NoError(t, err)

	dao := NewWindSpeedDaoImpl(db)
	windSpeedId := "testWindSpeed2"
	value := 2.0
	now := time.Now()
	err = dao.InsertWindSpeed(windSpeedId, value, now)
	require.NoError(t, err)

	value = 2.0
	err = dao.UpdateWindSpeed(windSpeedId, value, now)
	require.NoError(t, err)

	windSpeed, err := getWindSpeed(db, windSpeedId)
	require.NoError(t, err)
	require.NotNil(t, windSpeed)
	require.Equal(t, windSpeedId, windSpeed["ID"])
	require.Equal(t, value, windSpeed["Value"])
}

func getWindSpeed(db *sqlx.DB, windSpeedId string) (map[string]any, error) {
	stmt, err := db.Prepare(`SELECT ws.id, ws.value, ws.created_at, ws.updated_at FROM wind_speeds AS ws WHERE ws.id = ?`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err.Error())
		}
	}(stmt)
	row := stmt.QueryRow(windSpeedId)
	if row != nil {
		var id string
		var value float64
		var createdAt time.Time
		var updatedAt time.Time
		err = row.Scan(&id, &value, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"ID":        id,
			"Value":     value,
			"CreatedAt": createdAt.String(),
			"UpdatedAt": updatedAt.String(),
		}, nil
	}
	return nil, nil
}
