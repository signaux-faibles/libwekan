package libwekan

import "fmt"

type UserAlreadyExistsError error

func NewUserAlreadyExistsError(user User) UserAlreadyExistsError {
	return fmt.Errorf("l'utilisateur existe déjà (UserID: %s, Username: %s)", user.ID, user.Username)

}

type UnknownBoardError error

func NewUnknownBoardError(board Board) UnknownBoardError {
	return fmt.Errorf("la board est inconnue (BoardID: %s, Title: %s, Slug: %s", board.ID, board.Title, board.Slug)
}
