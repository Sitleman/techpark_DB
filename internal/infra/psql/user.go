package psql

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"techpark_db/internal/domain/entity"
)

const queryGetUser = "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1"

func (store *Storage) GetUser(tx *sql.Tx, nickname string) (*entity.User, error) {
	row := tx.QueryRow(queryGetUser, nickname)
	user := entity.User{}
	if err := row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

const queryFindUser = "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2"

func (store *Storage) FindUser(tx *sql.Tx, nickname string, email string) (*[]entity.User, error) {
	rows, err := tx.Query(queryFindUser, nickname, email)
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

const querySaveUser = "INSERT INTO Users(Nickname, Fullname, About, Email) VALUES ($1, $2, $3, $4)"

func (store *Storage) SaveUser(tx *sql.Tx, user entity.CreateUser, nickname string) error {
	if _, err := tx.Exec(querySaveUser, nickname, user.Fullname, user.About, user.Email); err != nil {
		return err
	}
	return nil
}

const queryUpdateUser = "UPDATE Users SET Fullname = $1, About = $2, Email = $3 WHERE LOWER(Nickname) = LOWER($4)"

func (store *Storage) UpdateUser(tx *sql.Tx, user entity.UpdateUser, nickname string) error {
	if _, err := tx.Exec(queryUpdateUser, user.Fullname, user.About, user.Email, nickname); err != nil {
		return err
	}
	return nil
}
