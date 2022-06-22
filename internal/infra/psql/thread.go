package psql

import (
	"errors"
	"strconv"
	"techpark_db/internal/domain/entity"
)

const querySaveThread = "INSERT INTO Thread(Title, Author, Message, Forum, Slug, Created) VALUES ($1, $2, $3, $4, $5, $6::TIMESTAMP WITH TIME ZONE) RETURNING id"

func (store *Storage) SaveThread(thread entity.CreateThread, slugForum string) (int, error) {
	row := store.DB.QueryRow(querySaveThread, thread.Title, thread.Author, thread.Message, slugForum, thread.Slug, thread.Created)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

const queryUpdateThreadVote = "UPDATE Thread SET Votes = $2 WHERE Id = $1"

func (store *Storage) UpdateThreadVote(thread entity.Thread) error {
	_, err := store.DB.Exec(queryUpdateThreadVote, thread.Id, thread.Votes)
	return err
}

const queryUpdateThread = "UPDATE Thread SET Title = $2, Author = $3, Forum = $4, Message = $5, Slug = $6 WHERE Id = $1"

func (store *Storage) UpdateThread(thread entity.Thread) error {
	_, err := store.DB.Exec(queryUpdateThread, thread.Id, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug)
	return err
}

const queryGetThread = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE LOWER(Slug) = LOWER($1) OR Id = $2"

func (store *Storage) GetThread(slugOrId string) (*entity.Thread, error) {
	if slugOrId == "" {
		return nil, errors.New("Empty slug")
	}
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		id = 0
	}
	row := store.DB.QueryRow(queryGetThread, slugOrId, id)
	thread := entity.Thread{}
	if err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
		return nil, err
	}
	return &thread, nil
}

const queryGetThreadByTitle = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE Title = $1"

func (store *Storage) GetThreadByTitle(title string) (*entity.Thread, error) {
	row := store.DB.QueryRow(queryGetThreadByTitle, title)
	thread := entity.Thread{}
	if err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
		return nil, err
	}
	return &thread, nil
}

const queryGetThreadByID = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE Id = $1"

func (store *Storage) GetThreadById(id int) (*entity.Thread, error) {
	row := store.DB.QueryRow(queryGetThreadByID, id)
	thread := entity.Thread{}
	if err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
		return nil, err
	}
	return &thread, nil
}
