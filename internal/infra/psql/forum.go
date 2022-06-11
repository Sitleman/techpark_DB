package psql

import (
	log "github.com/sirupsen/logrus"
	"techpark_db/internal/domain/entity"
)

const querySaveForum = "INSERT INTO Forum(Slug, Title, Nickname) VALUES ($1, $2, $3)"

func (store *Storage) SaveForum(forum entity.CreateForum) error {
	if _, err := store.DB.Exec(querySaveForum, forum.Slug, forum.Title, forum.User); err != nil {
		return err
	}
	return nil
}

const queryGetForum = "SELECT Slug, Title, User, Posts, Threads FROM Forum WHERE Slug = $1"

func (store *Storage) GetForum(slug string) (*entity.Forum, error) {
	row := store.DB.QueryRow(queryGetForum, slug)
	forum := entity.Forum{}
	if err := row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads); err != nil {
		log.Info(err, "[slug: ", slug, "]")
		return nil, err
	}
	return &forum, nil
}
