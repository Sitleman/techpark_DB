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

const queryClear = "TRUNCATE Vote, Posts, Thread, Forum, Users CASCADE"
const queryClearPosts = "TRUNCATE TABLE Posts"
const queryClearThread = "TRUNCATE TABLE Thread"
const queryClearForum = "TRUNCATE TABLE Forum"
const queryClearVote = "TRUNCATE TABLE Vote"

func (store *Storage) ClearData() error {
	if _, err := store.DB.Exec(queryClear); err != nil {
		return err
	}
	log.Info("clear data")
	//if _, err := tx.Exec(queryClearPosts); err != nil {
	//	return err
	//}
	//log.Info("clear posts")
	//if _, err := tx.Exec(queryClearThread); err != nil {
	//	return err
	//}
	//log.Info("clear thread")
	//if _, err := tx.Exec(queryClearForum); err != nil {
	//	return err
	//}
	//log.Info("clear forum")
	//if _, err := tx.Exec(queryClearUsers); err != nil {
	//	return err
	//}
	//log.Info("clear users")
	return nil
}
