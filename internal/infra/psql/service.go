package psql

import (
	log "github.com/sirupsen/logrus"
	"techpark_db/internal/domain/entity"
)

const queryGetUserCount = "SELECT COUNT(*) FROM Users"
const queryGetForumCount = "SELECT COUNT(*) FROM Forum"
const queryGetThreadCount = "SELECT COUNT(*) FROM Thread"
const queryGetPostCount = "SELECT COUNT(*) FROM Posts"

func (store *Storage) GetServiceStatus() (*entity.ServStatus, error) {
	var servStatus entity.ServStatus
	row := store.DB.QueryRow(queryGetUserCount)
	if err := row.Scan(&servStatus.User); err != nil {
		return nil, err
	}

	row = store.DB.QueryRow(queryGetForumCount)
	if err := row.Scan(&servStatus.Forum); err != nil {
		return nil, err
	}

	row = store.DB.QueryRow(queryGetThreadCount)
	if err := row.Scan(&servStatus.Thread); err != nil {
		return nil, err
	}

	row = store.DB.QueryRow(queryGetPostCount)
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

func (store *Storage) ClearData() error {
	if _, err := store.DB.Exec(queryClearVote); err != nil {
		return err
	}
	log.Info("clear vote")
	if _, err := store.DB.Exec(queryClearPosts); err != nil {
		return err
	}
	log.Info("clear posts")
	if _, err := store.DB.Exec(queryClearThread); err != nil {
		return err
	}
	log.Info("clear thread")
	if _, err := store.DB.Exec(queryClearForum); err != nil {
		return err
	}
	log.Info("clear forum")
	if _, err := store.DB.Exec(queryClearUsers); err != nil {
		return err
	}
	log.Info("clear users")
	return nil
}
