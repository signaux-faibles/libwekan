package libwekan

import "fmt"

type Err struct {
	err string
	id  string
}

func (e Err) Error() string {
	return fmt.Sprintf("une erreur est survenue: %s, l'id incriminé est: %s", e.err, e.id)
}

type UserAlreadyExistsError struct {
	Err
}

func NewUserAlreadyExistsError(id string) UserAlreadyExistsError {
	return UserAlreadyExistsError{
		Err{
			err: "l'utilisateur existe déjà",
			id:  id,
		},
	}
}
