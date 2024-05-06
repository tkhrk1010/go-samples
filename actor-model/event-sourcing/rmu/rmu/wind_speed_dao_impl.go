package rmu

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type WindSpeedDaoImpl struct {
	db *sqlx.DB
}

func NewWindSpeedDaoImpl(db *sqlx.DB) WindSpeedDaoImpl {
	return WindSpeedDaoImpl{db}
}

func (dao *WindSpeedDaoImpl) InsertWindSpeed(id string, value float64, createdAt time.Time) error {
	stmt, err := dao.db.Prepare(`INSERT INTO wind_speeds (id, value, created_at, updated_at) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err.Error())
		}
	}(stmt)
	dt := createdAt.Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(id, value, dt, dt)
	if err != nil {
		return err
	}
	return nil
}

func (dao *WindSpeedDaoImpl) UpdateWindSpeed(id string, value float64, updatedAt time.Time) error {
	stmt, err := dao.db.Prepare(`UPDATE wind_speeds SET value = ?, updated_at = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err.Error())
		}
	}(stmt)
	dt := updatedAt.Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(value, dt, id)
	if err != nil {
		return err
	}
	return nil
}
