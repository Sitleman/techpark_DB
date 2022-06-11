package psql

import (
	"strconv"
	"techpark_db/internal/domain/entity"
	"techpark_db/internal/utils"
)

const querySaveThread = "INSERT INTO Thread(Title, Author, Message, Forum, Slug, Created) VALUES ($1, $2, $3, $4, $5, now())"

func (store *Storage) SaveThread(thread entity.CreateThread, slugForum string) error {
	if _, err := store.DB.Exec(querySaveThread, thread.Title, thread.Author, thread.Message, slugForum, utils.RandSlug()); err != nil {
		return err
	}
	return nil
}

const queryGetThread = "SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread WHERE Slug = $1 OR Id = $2"

func (store *Storage) GetThread(slugOrId string) (*entity.Thread, error) {
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
