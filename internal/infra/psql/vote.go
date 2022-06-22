package psql

import (
	"techpark_db/internal/domain/entity"
)

const queryGetVote = "SELECT IdThread, Nickname, Voice FROM Vote WHERE IdThread = $1 AND LOWER(Nickname) = LOWER($2)"

func (store *Storage) GetVote(idThread int, nickname string) (*entity.Vote, error) {
	row := store.DB.QueryRow(queryGetVote, idThread, nickname)
	vote := entity.Vote{}
	if err := row.Scan(&vote.IdThread, &vote.Nickname, &vote.Voice); err != nil {
		return nil, err
	}
	return &vote, nil
}

const querySaveVote = "INSERT INTO Vote(IdThread, Nickname, Voice) VALUES ($1, $2, $3)"

func (store *Storage) SaveVote(voteReq entity.Vote) error {
	_, err := store.DB.Exec(querySaveVote, voteReq.IdThread, voteReq.Nickname, voteReq.Voice)
	if err != nil {
		return err
	}
	return nil
}

const queryUpdateVote = "UPDATE Vote SET Voice = $3 WHERE IdThread = $1 AND LOWER(Nickname) = LOWER($2)"

func (store *Storage) UpdateVote(voteReq entity.Vote) error {
	_, err := store.DB.Exec(queryUpdateVote, voteReq.IdThread, voteReq.Nickname, voteReq.Voice)
	if err != nil {
		return err
	}
	return nil
}
