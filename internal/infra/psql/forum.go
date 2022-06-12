package psql

import (
	"database/sql"
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

const queryGetForum = "SELECT Slug, Title, Nickname, Posts, Threads FROM Forum WHERE LOWER(Slug) = LOWER($1)"

func (store *Storage) GetForum(slug string) (*entity.Forum, error) {
	row := store.DB.QueryRow(queryGetForum, slug)
	forum := entity.Forum{}
	if err := row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads); err != nil {
		log.Info(err, "[slug: ", slug, "]")
		return nil, err
	}
	return &forum, nil
}

const queryGetForumThreads = `
SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread
WHERE LOWER(Forum) = LOWER($1) AND Created >= $2::TIMESTAMP WITH TIME ZONE
ORDER BY Created 
LIMIT $3
`

const queryGetForumThreadsDesc = `
SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread
WHERE LOWER(Forum) = LOWER($1) AND Created <= $2::TIMESTAMP WITH TIME ZONE
ORDER BY Created DESC
LIMIT $3
`

func (store *Storage) GetForumThreads(slug string, order string, limit int, since string) (*[]entity.Thread, error) {
	var rows *sql.Rows
	var err error
	if order == "ASC" {
		rows, err = store.DB.Query(queryGetForumThreads, slug, since, limit)
	} else {
		rows, err = store.DB.Query(queryGetForumThreadsDesc, slug, since, limit)
	}

	if err != nil {
		log.Error(err, "[slug ", slug, "] [order ", order, "] [limit ", limit, "] [since ", since, "]")
		return nil, err
	}
	defer rows.Close()

	threads := make([]entity.Thread, 0)
	for rows.Next() {
		thread := entity.Thread{}
		if err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
			log.Error(err)
			return nil, err
		}
		threads = append(threads, thread)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	return &threads, nil
}

const queryGetForumUsers = `
SELECT Nickname, Fullname, About, Email FROM Users
JOIN Thread ON Users.Nickname = Thread.Author
WHERE Forum = $1
ORDER BY LOWER(Nickname) $2
LIMIT $3
OFFSET $4
`

func (store *Storage) GetForumUsers(slug string, order string, limit int, since int) (*[]entity.User, error) {
	rows, err := store.DB.Query(queryGetForumUsers, slug, order, limit, since)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	users := make([]entity.User, 0)
	for rows.Next() {
		user := entity.User{}
		if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	return &users, nil
}
