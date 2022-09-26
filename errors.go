package libwekan

import (
	"fmt"
)

type libwekanError struct {
	err string
}

func (e libwekanError) Error() string {
	return e.err
}

type UserAlreadyExistsError struct {
	libwekanError
}

func NewUserAlreadyExistsError(user User) UserAlreadyExistsError {
	return UserAlreadyExistsError{libwekanError{err: fmt.Sprintf("l'utilisateur existe déjà (UserID: %s, Username: %s)", user.ID, user.Username)}}
}

type UnknownBoardError struct {
	libwekanError
}

func NewUnknownBoardError(board Board) UnknownBoardError {
	return UnknownBoardError{libwekanError{err: fmt.Sprintf("la board est inconnue (BoardID: %s, Title: %s, Slug: %s", board.ID, board.Title, board.Slug)}}
}

type InsertEmptyRuleError struct {
	libwekanError
}

func NewInsertEmptyRuleError() InsertEmptyRuleError {
	return InsertEmptyRuleError{libwekanError{err: "l'insertion d'un objet Rule vide est impossible"}}
}
