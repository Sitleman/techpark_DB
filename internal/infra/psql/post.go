package psql

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"techpark_db/internal/domain/entity"
)

const queryCheckParentPost = "SELECT count(Id) FROM Posts WHERE Id = $1 AND Thread = $2"

func (store *Storage) CheckParentPost(tx *sql.Tx, parent int, threadId int) (bool, error) {
	row := tx.QueryRow(queryCheckParentPost, parent, threadId)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

const querySavePost = "INSERT INTO Posts(Parent, Author, Message, Forum, Thread, Created) VALUES "

func (store *Storage) SavePosts(tx *sql.Tx, posts []entity.CreatePost, forum string, thread int, created string) (*[]int, error) {
	query := querySavePost
	args := make([]interface{}, 0, len(posts))
	for i, post := range posts {
		query += fmt.Sprintf(" ($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		args = append(args, post.Parent, post.Author, post.Message, forum, thread, created)
	}
	query = query[:len(query)-1]
	query += " RETURNING Id"

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0, len(posts))
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Error(err)
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}
	return &ids, nil
}

//const queryGetPosts = "SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts WHERE Created = $1::TIMESTAMP WITH TIME ZONE ORDER BY Id"
//
//func (store *Storage) GetPostsByCreated(tx *sql.Tx, created string) (*[]entity.Post, error) {
//	rows, err := tx.Query(queryGetPosts, created)
//	if err != nil {
//		log.Error(err, "[created ", created, "]")
//		return nil, err
//	}
//	defer rows.Close()
//
//	posts := make([]entity.Post, 0)
//	for rows.Next() {
//		post := entity.Post{}
//		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		posts = append(posts, post)
//	}
//	if err := rows.Err(); err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	return &posts, nil
//}

const queryGetPostById = "SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts WHERE Id = $1"

func (store *Storage) GetPostById(tx *sql.Tx, id int) (*entity.Post, error) {
	row := tx.QueryRow(queryGetPostById, id)

	var post entity.Post
	if err := row.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
		return nil, err
	}
	return &post, nil
}

const queryUpdatePost = "UPDATE Posts SET Message = $2, IsEdited = true WHERE Id = $1"

func (store *Storage) UpdatePost(tx *sql.Tx, id int, message string) error {
	_, err := tx.Exec(queryUpdatePost, id, message)
	return err
}

const queryGetPostsFlat = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1
ORDER BY Id
LIMIT $2
`

const queryGetPostsFlatDesc = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1
ORDER BY Id DESC
LIMIT $2
`

const queryGetPostsFlatSince = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1 AND Id > $3
ORDER BY Id
LIMIT $2
`

const queryGetPostsFlatSinceDesc = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1 AND Id < $3
ORDER BY Id DESC
LIMIT $2
`

func (store *Storage) GetPostsByThreadFlat(tx *sql.Tx, thread int, limit int, since int, sort string, order string) (*[]entity.Post, error) {
	var rows *sql.Rows
	err := errors.New("undefined")
	if since == 0 {
		if order == "ASC" {
			rows, err = tx.Query(queryGetPostsFlat, thread, limit)
		}
		if order == "DESC" {
			rows, err = tx.Query(queryGetPostsFlatDesc, thread, limit)
		}
	} else {
		if order == "ASC" {
			rows, err = tx.Query(queryGetPostsFlatSince, thread, limit, since)
		}
		if order == "DESC" {
			rows, err = tx.Query(queryGetPostsFlatSinceDesc, thread, limit, since)
		}
	}

	if err != nil {
		log.Error(err, "[thread ", thread, "]", "[limit ", limit, "] [since ", since, "] [sort ", sort, "] [order ", order, "]")
		return nil, err
	}
	defer rows.Close()

	posts := make([]entity.Post, 0, 100)
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			log.Error(err)
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	return &posts, nil
}

const queryGetPostsTree = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1
ORDER BY TreePath 
LIMIT $2
`

const queryGetPostsTreeDesc = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1
ORDER BY TreePath DESC
LIMIT $2
`

const queryGetPostsTreeSince = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1 AND Treepath > (SELECT Treepath FROM Posts WHERE Id = $3)
ORDER BY TreePath 
LIMIT $2
`

const queryGetPostsTreeSinceDesc = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE Thread = $1 AND Treepath < (SELECT Treepath FROM Posts WHERE Id = $3)
ORDER BY TreePath DESC
LIMIT $2
`

func (store *Storage) GetPostsTree(tx *sql.Tx, thread int, limit int, since int, sort string, order string) (*[]entity.Post, error) {
	var rows *sql.Rows
	err := errors.New("undefined")
	if since == 0 {
		if order == "ASC" {
			rows, err = tx.Query(queryGetPostsTree, thread, limit)
		} else {
			rows, err = tx.Query(queryGetPostsTreeDesc, thread, limit)
		}
	} else {
		if order == "ASC" {
			rows, err = tx.Query(queryGetPostsTreeSince, thread, limit, since)
		} else {
			rows, err = tx.Query(queryGetPostsTreeSinceDesc, thread, limit, since)
		}
	}

	if err != nil {
		log.Error(err, "[sort ", sort, "] [order ", order, "]")
		return nil, err
	}
	defer rows.Close()

	posts := make([]entity.Post, 0)
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			log.Error(err)
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}
	return &posts, nil
}

const queryGetPostsParentTree = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE TreePath[1] IN (SELECT Id FROM Posts WHERE Thread = $1 AND Parent = 0 ORDER BY Id LIMIT $2)
ORDER BY TreePath
`
const queryGetPostsParentTreeDesc = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE TreePath[1] IN (SELECT Id FROM Posts WHERE Thread = $1 AND Parent = 0 ORDER BY Id DESC LIMIT $2)
ORDER BY TreePath[1] DESC, TreePath
`
const queryGetPostsParentTreeSince = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE TreePath[1] IN 
(SELECT Id FROM Posts WHERE Thread = $1 AND Parent = 0 AND Id > (SELECT TreePath[1] FROM Posts WHERE Id = $3) 
ORDER BY Id LIMIT $2)
ORDER BY TreePath
`
const queryGetPostsParentTreeSinceDesc = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
WHERE TreePath[1] IN 
(SELECT Id FROM Posts WHERE Thread = $1 AND Parent = 0 AND Id < (SELECT TreePath[1] FROM Posts WHERE Id = $3) 
ORDER BY Id DESC LIMIT $2)
ORDER BY TreePath[1] DESC, TreePath
`

func (store *Storage) GetPostsParentTree(tx *sql.Tx, thread int, limit int, since int, sort string, order string) (*[]entity.Post, error) {
	var rows *sql.Rows
	err := errors.New("undefined")

	if since == 0 {
		if order == "ASC" {
			rows, err = tx.Query(queryGetPostsParentTree, thread, limit)
		} else {
			rows, err = tx.Query(queryGetPostsParentTreeDesc, thread, limit)
		}
	} else {
		if order == "ASC" {
			rows, err = tx.Query(queryGetPostsParentTreeSince, thread, limit, since)
		} else {
			rows, err = tx.Query(queryGetPostsParentTreeSinceDesc, thread, limit, since)
		}
	}

	if err != nil {
		log.Error(err, "[sort ", sort, "] [order ", order, "]")
		return nil, err
	}
	defer rows.Close()

	posts := make([]entity.Post, 0)
	if sort == "parent_tree" {
		limit = INF
	}
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			log.Error(err)
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}
	return &posts, nil
}

//const queryGetPostsByParent = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
//WHERE Parent = $1
//ORDER BY Id
//LIMIT $2
//`
//
//const queryGetPostsByParentDesc = `SELECT Id, Parent, Author, Message, IsEdited, Forum, Thread, Created FROM Posts
//WHERE Parent = $1
//ORDER BY Id DESC
//LIMIT $2
//`
//
//func (store *Storage) GetPostsByParent(parent int, limit int, order string, since int, sinceFlag *bool) (*[]entity.Post, error) {
//	var rows *sql.Rows
//	err := errors.New("undefined")
//	if order == "ASC" {
//		rows, err = store.DB.Query(queryGetPostsByParent, parent, limit)
//	}
//	if order == "DESC" {
//		rows, err = store.DB.Query(queryGetPostsByParentDesc, parent, limit)
//	}
//	if err != nil {
//		log.Error(err, "[parent ", parent, "] [order ", order, "]", "[limit ", limit, "]")
//		return nil, err
//	}
//	defer rows.Close()
//
//	count := 0
//	posts := make([]entity.Post, 0)
//	for rows.Next() {
//		post := entity.Post{}
//		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		if *sinceFlag {
//			count++
//			posts = append(posts, post)
//		}
//
//		if post.Id == since {
//			*sinceFlag = true
//		}
//
//		childPosts, err := store.GetPostsByParent(post.Id, limit-count, order, since, sinceFlag)
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//
//		for _, childPost := range *childPosts {
//			if *sinceFlag {
//				posts = append(posts, childPost)
//				count++
//			}
//		}
//
//		if count >= limit {
//			break
//		}
//	}
//	if err := rows.Err(); err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	return &posts, nil
//}
