package psql

import (
	"database/sql"
	"errors"
	"strconv"
	"techpark_db/internal/domain/entity"
)

const querySaveThread = "INSERT INTO Thread(Title, Author, Message, Forum, Slug, Created) VALUES ($1, $2, $3, $4, $5, $6::TIMESTAMP WITH TIME ZONE) RETURNING id"

func (store *Storage) SaveThread(tx *sql.Tx, thread entity.CreateThread, slugForum string) (int, error) {
	row := tx.QueryRow(querySaveThread, thread.Title, thread.Author, thread.Message, slugForum, thread.Slug, thread.Created)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

const queryUpdateThreadVote = "UPDATE Thread SET Votes = $2 WHERE Id = $1"

func (store *Storage) UpdateThreadVote(tx *sql.Tx, thread entity.Thread) error {
	_, err := tx.Exec(queryUpdateThreadVote, thread.Id, thread.Votes)
	return err
}

const queryUpdateThread = "UPDATE Thread SET Title = $2, Author = $3, Forum = $4, Message = $5, Slug = $6 WHERE Id = $1"

func (store *Storage) UpdateThread(tx *sql.Tx, thread entity.Thread) error {
	_, err := tx.Exec(queryUpdateThread, thread.Id, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug)
	return err
}

const queryGetThreadId = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE Id = $1"
const queryGetThreadSlug = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE Slug = $1"

func (store *Storage) GetThread(tx *sql.Tx, slugOrId string) (*entity.Thread, error) {
	if slugOrId == "" {
		return nil, errors.New("Empty slug")
	}
	var row *sql.Row
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		row = tx.QueryRow(queryGetThreadSlug, slugOrId)
	} else {
		row = tx.QueryRow(queryGetThreadId, id)
	}
	thread := entity.Thread{}
	if err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
		return nil, err
	}
	return &thread, nil
}

const queryCountVote = "SELECT Votes FROM Thread WHERE Id = $1"

func (store *Storage) CountVote(tx *sql.Tx, id int) (*int, error) {
	row := tx.QueryRow(queryCountVote, id)
	var count int
	if err := row.Scan(&count); err != nil {
		return nil, err
	}
	return &count, nil
}

const queryGetThreadByTitle = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE Title = $1"

func (store *Storage) GetThreadByTitle(tx *sql.Tx, title string) (*entity.Thread, error) {
	row := tx.QueryRow(queryGetThreadByTitle, title)
	thread := entity.Thread{}
	if err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
		return nil, err
	}
	return &thread, nil
}

const queryGetThreadByID = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE Id = $1"

func (store *Storage) GetThreadById(tx *sql.Tx, id int) (*entity.Thread, error) {
	row := tx.QueryRow(queryGetThreadByID, id)
	thread := entity.Thread{}
	if err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
		return nil, err
	}
	return &thread, nil
}
