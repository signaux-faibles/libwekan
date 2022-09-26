package libwekan

import (
	"fmt"
)

//type libwekanError struct {
//	err string
//}
//
//func (e libwekanError) Error() string {
//	return e.err
//}
//
//type UserAlreadyExistsError struct {
//	libwekanError
//}
//
//func NewUserAlreadyExistsError(user User) UserAlreadyExistsError {
//	return UserAlreadyExistsError{libwekanError{err: fmt.Sprintf("l'utilisateur existe déjà (UserID: %s, Username: %s)", user.ID, user.Username)}}
//}
//
//type UnknownBoardError struct {
//	libwekanError
//}
//
//func NewUnknownBoardError(board Board) UnknownBoardError {
//	return UnknownBoardError{libwekanError{err: fmt.Sprintf("la board est inconnue (BoardID: %s, Title: %s, Slug: %s", board.ID, board.Title, board.Slug)}}
//}
//
//type InsertEmptyRuleError struct {
//	libwekanError
//}
//
//func NewInsertEmptyRuleError() InsertEmptyRuleError {
//	return InsertEmptyRuleError{libwekanError{err: "l'insertion d'un objet Rule vide est impossible"}}
//}

type UserAlreadyExistsError struct {
	user User
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("l'utilisateur existe déjà (UserID: %s, Username: %s)", e.user.ID, e.user.Username)
}

type UnknownBoardError struct {
	board Board
}

func (e UnknownBoardError) Error() string {
	return fmt.Sprintf("la board est inconnue (BoardID: %s, Title: %s, Slug: %s", e.board.ID, e.board.Title, e.board.Slug)
}

type InsertEmptyRuleError struct {
}

func (e InsertEmptyRuleError) Error() string {
	return fmt.Sprintf("l'insertion d'un objet Rule vide est impossible")
}

type BoardLabelAlreadyExistsError struct {
	boardLabel BoardLabel
	board      Board
}

func (e BoardLabelAlreadyExistsError) Error() string {
	return fmt.Sprintf("un objet BoardLabel existe déjà dans la board (%s) avec le même nom (%s)", e.board.ID, e.boardLabel.Name)
}
