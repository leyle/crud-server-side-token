package sstapp

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbTableRevokeList = "revoke_list"
)

const createRevokeTable string = `
	CREATE TABLE IF NOT EXISTS revoke_list (
		id INTEGER NOT NULL PRIMARY KEY,
	 	token TEXT UNIQUE,
	 	userId TEXT,
		t INTEGER
	);
`

func (sst *SSTokenOption) getDb() error {
	db, err := sql.Open("sqlite3", sst.sqliteFile)
	if err != nil {
		sst.logger.Error().Err(err).Msg("open sqlite3 db file failed")
		return err
	}

	err = db.Ping()
	if err != nil {
		sst.logger.Error().Err(err).Msg("ping sqlite3 failed")
		return err
	}

	sst.logger.Debug().Msgf("get sqlite db ok, db name[%s]", sst.sqliteFile)

	sst.db = db
	return nil
}

func (sst *SSTokenOption) createTable() error {
	result, err := sst.db.Exec(createRevokeTable)
	if err != nil {
		sst.logger.Error().Err(err).Msg("create sqlite table failed")
		return err
	}
	affected, _ := result.RowsAffected()
	sst.logger.Debug().Msgf("create sqlite3 table, affected rows[%d]", affected)
	return nil
}

func (sst *SSTokenOption) insertIntoRevokeList(token, userId string, t int64) error {
	ctx := context.Background()
	tx, err := sst.db.BeginTx(ctx, nil)
	if err != nil {
		sst.logger.Error().Err(err).Msg("sqlite3 start transaction failed")
		return err
	}
	defer tx.Rollback()

	var dataId int
	row := tx.QueryRowContext(ctx, "SELECT id from revoke_list where token = ?", token)
	err = row.Scan(&dataId)
	if err != nil && err != sql.ErrNoRows {
		sst.logger.Error().Err(err).Msg("query token from sqlite3 failed")
		return err
	}

	if dataId != 0 {
		// return ok
		sst.logger.Warn().Str("token", token).Msg("insert revoke token, token already revoked")
		return nil
	}

	// insert into db
	_, err = tx.ExecContext(ctx, "INSERT INTO revoke_list(token, t, userId) VALUES(?, ?, ?)", token, t, userId)
	if err != nil {
		sst.logger.Error().Err(err).Msg("insert revoke list failed")
		return err
	}

	err = tx.Commit()
	if err != nil {
		sst.logger.Error().Err(err).Msg("commit transaction failed")
		return err
	}

	sst.logger.Info().Str("token", token).Msg("save token into revoke list succeed")
	return nil
}

func (sst *SSTokenOption) loadRevokeList() error {
	if sst.db == nil {
		err := sst.getDb()
		if err != nil {
			sst.logger.Error().Err(err).Msg("load revoke list failed")
			return err
		}
	}

	rows, err := sst.db.Query("SELECT token, t, userId FROM revoke_list")
	if err != nil {
		sst.logger.Error().Err(err).Msg("select from revoke list failed")
		return err
	}

	for rows.Next() {
		rv := &revokedToken{}
		err = rows.Scan(&rv.token, &rv.t, &rv.userId)
		if err != nil {
			sst.logger.Error().Err(err).Msg("scan sqlite rows failed")
			return err
		}
		sst.revokeList = append(sst.revokeList, rv)
	}

	sst.logger.Info().Msgf("load revoke list from sqlite succeed, total num[%d]", len(sst.revokeList))
	return nil
}
