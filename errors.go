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

type UserNotFoundError struct {
	key string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("l'utilisateur n'est pas connu (%s)", e.key)
}

type BoardNotFoundError struct {
	board Board
}

func (e BoardNotFoundError) Error() string {
	return fmt.Sprintf("la board est inconnue (BoardID: %s, Title: %s, Slug: %s", e.board.ID, e.board.Title, e.board.Slug)
}

type NotPrivilegedError struct {
	id  UserID
	err error
}

func (e NotPrivilegedError) Unwrap() error {
	return e.err
}

func (e NotPrivilegedError) Error() string {
	return fmt.Sprintf("l'utilisateur n'est pas admin: id = %s", e.id)
}

type ProtectedUserError struct {
	id UserID
}

func (e ProtectedUserError) Error() string {
	return fmt.Sprintf("cet action est interdite sur cet utilisateur (%s)", e.id)
}

type AdminUserIsNotAdminError struct {
	username Username
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

type BoardLabelNotFoundError struct {
	boardLabelID BoardLabelID
	board        Board
}

func (e BoardLabelNotFoundError) Error() string {
	return fmt.Sprintf("l'objet BoardLabel (id=%s) n'a pas été trouvé dans la board (%s)", e.boardLabelID, e.board.ID)
}

type UnexpectedMongoError struct {
	err error
}

func (e UnexpectedMongoError) Error() string {
	return "une erreur est survenue lors de l'exécution de la requête"
}

func (e UnexpectedMongoError) Unwrap() error {
	return e.err
}

type AlreadySetActivityError struct {
	activityType string
}

func (e AlreadySetActivityError) Error() string {
	return fmt.Sprintf("l'activité est déjà définie: activityType = %s", e.activityType)
}

type UnreachableMongoError struct {
	err error
}

func (e UnreachableMongoError) Error() string {
	return "la connexion a échoué"
}

func (e UnreachableMongoError) Unwrap() error {
	return e.err
}

type InvalidMongoConfigurationError struct {
	err error
}

func (e InvalidMongoConfigurationError) Error() string {
	return "les paramètres de connexion sont invalides"
}

func (e InvalidMongoConfigurationError) Unwrap() error {
	return e.err
}

type ForbiddenOperationError struct {
	err error
}

func (e ForbiddenOperationError) Error() string {
	return "operation interdite"
}

func (e ForbiddenOperationError) Unwrap() error {
	return e.err
}

type NotImplemented struct {
	method string
}

func (e NotImplemented) Error() string {
	return "not implemented : " + e.method
}

type CardNotFoundError struct {
	cardID CardID
}

func (e CardNotFoundError) Error() string {
	return fmt.Sprintf("la carte n'existe pas (ID: %s)", e.cardID)
}

type RuleNotFoundError struct {
	ruleID RuleID
}

func (e RuleNotFoundError) Error() string {
	return fmt.Sprintf("la règle n'existe pas (ID: %s)", e.ruleID)
}

type ActionNotFoundError struct {
	actionId ActionID
}

func (e ActionNotFoundError) Error() string {
	return fmt.Sprintf("l'action n'existe pas (ID: %s)", e.actionId)
}

type TriggerNotFoundError struct {
	triggerId TriggerID
}

func (e TriggerNotFoundError) Error() string {
	return fmt.Sprintf("le trigger n'existe pas (ID: %s)", e.triggerId)
}

type NothingDoneError struct{}

func (e NothingDoneError) Error() string {
	return "le traitement n'a eu aucun effet"
}

type UnknownActivityError struct {
	key string
}

func (e UnknownActivityError) Error() string {
	return fmt.Sprintf("l'activité n'existe pas (ID: %s)", e.key)
}

type UserIsNotMemberError struct {
	userID UserID
}

func (e UserIsNotMemberError) Error() string {
	return fmt.Sprintf("cette action nécessite que l'utilisateur soit membre (id=%s)", e.userID)
}
