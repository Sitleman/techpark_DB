package psql

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"techpark_db/internal/domain/entity"
)

const queryGetUserCount = "SELECT COUNT(*) FROM Users"
const queryGetForumCount = "SELECT COUNT(*) FROM Forum"
const queryGetThreadCount = "SELECT COUNT(*) FROM Thread"
const queryGetPostCount = "SELECT COUNT(*) FROM Posts"

func (store *Storage) GetServiceStatus(tx *sql.Tx) (*entity.ServStatus, error) {
	var servStatus entity.ServStatus
	row := tx.QueryRow(queryGetUserCount)
	if err := row.Scan(&servStatus.User); err != nil {
		return nil, err
	}

	row = tx.QueryRow(queryGetForumCount)
	if err := row.Scan(&servStatus.Forum); err != nil {
		return nil, err
	}

	row = tx.QueryRow(queryGetThreadCount)
	if err := row.Scan(&servStatus.Thread); err != nil {
		return nil, err
	}

	row = tx.QueryRow(queryGetPostCount)
	if err := row.Scan(&servStatus.Post); err != nil {
		return nil, err
	}

	return &servStatus, nil
}

const queryClearUsers = "DELETE FROM Users"
const queryClearPosts = "DELETE FROM Posts"
const queryClearThread = "DELETE FROM Thread"
const queryClearForum = "DELETE FROM Forum"
const queryClearVote = "DELETE FROM Vote"

func (store *Storage) ClearData(tx *sql.Tx) error {
	if _, err := tx.Exec(queryClearVote); err != nil {
		return err
	}
	log.Info("clear vote")
	if _, err := tx.Exec(queryClearPosts); err != nil {
		return err
	}
	log.Info("clear posts")
	if _, err := tx.Exec(queryClearThread); err != nil {
		return err
	}
	log.Info("clear thread")
	if _, err := tx.Exec(queryClearForum); err != nil {
		return err
	}
	log.Info("clear forum")
	if _, err := tx.Exec(queryClearUsers); err != nil {
		return err
	}
	log.Info("clear users")
	return nil
}
