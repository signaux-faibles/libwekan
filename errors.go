package libwekan

import (
	"fmt"
)

type UserAlreadyExistsError struct {
	user User
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("l'utilisateur existe déjà (UserID: %s, Username: %s)", e.user.ID, e.user.Username)
}

type UnknownUserError struct {
	key string
}

func (e UnknownUserError) Error() string {
	return fmt.Sprintf("l'utilisateur n'est pas connu (%s)", e.key)
}

type UnknownBoardError struct {
	board Board
}

func (e UnknownBoardError) Error() string {
	return fmt.Sprintf("la board est inconnue (BoardID: %s, Title: %s, Slug: %s", e.board.ID, e.board.Title, e.board.Slug)
}

type UserIsNotAdminError struct {
	id UserID
}

func (e UserIsNotAdminError) Error() string {
	return fmt.Sprintf("l'utilisateur n'est pas admin: id = %s", e.id)
}

type InsertEmptyRuleError struct {
}

func (e InsertEmptyRuleError) Error() string {
	return "l'insertion d'un objet Rule vide est impossible"
}

type BoardLabelAlreadyExistsError struct {
	boardLabel BoardLabel
	board      Board
}

func (e BoardLabelAlreadyExistsError) Error() string {
	return fmt.Sprintf("un objet BoardLabel existe déjà dans la board (%s) avec le même nom (%s)", e.board.ID, e.boardLabel.Name)
}

type UnexpectedMongoError struct {
	err error
}

func (e UnexpectedMongoError) Error() string {
	return e.err.Error()
}
