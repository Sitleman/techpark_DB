package psql

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"techpark_db/internal/domain/entity"
)

const querySaveForum = "INSERT INTO Forum(Slug, Title, Nickname) VALUES ($1, $2, $3)"

func (store *Storage) SaveForum(tx *sql.Tx, forum entity.CreateForum) error {
	if _, err := tx.Exec(querySaveForum, forum.Slug, forum.Title, forum.User); err != nil {
		return err
	}
	return nil
}

const queryGetForum = "SELECT Slug, Title, Nickname, Posts, Threads FROM Forum WHERE Slug = $1"

func (store *Storage) GetForum(tx *sql.Tx, slug string) (*entity.Forum, error) {
	row := tx.QueryRow(queryGetForum, slug)
	forum := entity.Forum{}
	if err := row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads); err != nil {
		//log.Info(err, "[slug: ", slug, "]")
		return nil, err
	}
	return &forum, nil
}

const queryGetForumThreads = `
SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread
WHERE Forum = $1 AND Created >= $2::TIMESTAMP WITH TIME ZONE
ORDER BY Created 
LIMIT $3
`

const queryGetForumThreadsDesc = `
SELECT Id, Title, Author, Forum, Message, Votes, Slug, Created FROM Thread
WHERE Forum = $1 AND Created <= $2::TIMESTAMP WITH TIME ZONE
ORDER BY Created DESC
LIMIT $3
`

func (store *Storage) GetForumThreads(tx *sql.Tx, slug string, order string, limit int, since string) (*[]entity.Thread, error) {
	var rows *sql.Rows
	var err error
	if order == "ASC" {
		rows, err = tx.Query(queryGetForumThreads, slug, since, limit)
	} else {
		rows, err = tx.Query(queryGetForumThreadsDesc, slug, since, limit)
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
SELECT UsersForum.Nickname, Fullname, About, Email FROM UsersForum
JOIN Users ON Users.Nickname = UsersForum.Nickname
WHERE Forum = $1
ORDER BY UsersForum.Nickname
LIMIT $2
`

const queryGetForumUsersDesc = `
SELECT UsersForum.Nickname, Fullname, About, Email FROM UsersForum
JOIN Users ON Users.Nickname = UsersForum.Nickname
WHERE Forum = $1
ORDER BY UsersForum.Nickname DESC
LIMIT $2
`

const queryGetForumUsersSince = `
SELECT UsersForum.Nickname, Fullname, About, Email FROM UsersForum
JOIN Users ON Users.Nickname = UsersForum.Nickname
WHERE Forum = $1 AND UsersForum.Nickname > $3
ORDER BY UsersForum.Nickname
LIMIT $2
`

const queryGetForumUsersSinceDesc = `
SELECT UsersForum.Nickname, Fullname, About, Email FROM UsersForum 
JOIN Users ON Users.Nickname = UsersForum.Nickname
WHERE Forum = $1 AND UsersForum.Nickname < $3
ORDER BY UsersForum.Nickname DESC
LIMIT $2
`

func (store *Storage) GetForumUsers(tx *sql.Tx, slug string, order string, limit int, since string) (*[]entity.User, error) {
	var rows *sql.Rows
	var err error
	if since == "" {
		if order == "ASC" {
			rows, err = tx.Query(queryGetForumUsers, slug, limit)
		} else {
			rows, err = tx.Query(queryGetForumUsersDesc, slug, limit)
		}
	} else {
		if order == "ASC" {
			rows, err = tx.Query(queryGetForumUsersSince, slug, limit, since)
		} else {
			rows, err = tx.Query(queryGetForumUsersSinceDesc, slug, limit, since)
		}
	}

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
