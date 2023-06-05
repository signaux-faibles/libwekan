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
	err error
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("l'utilisateur n'est pas connu (%s)", e.key)
}

func (e UserNotFoundError) Unwrap() error {
	return e.err
}

type BoardNotFoundError struct {
	msg string
	err error
}

func boardNotFoundWithSlug(slug BoardSlug) error {
	return BoardNotFoundError{msg: fmt.Sprintf("le slug : '%s'", slug)}
}

func boardNotFoundWithId(boardID BoardID) error {
	return BoardNotFoundError{msg: fmt.Sprintf("le boardID : '%s'", boardID)}
}

func boardNotFoundWithTitle(boardTitle BoardTitle) error {
	return BoardNotFoundError{msg: fmt.Sprintf("le titre : '%s'", boardTitle)}
}

func (e BoardNotFoundError) Error() string {
	return fmt.Sprintf("aucun tableau n'a été trouvé avec %s", e.msg)
}

func (e BoardNotFoundError) Unwrap() error {
	return e.err
}

type NotPrivilegedError struct {
	id  UserID
	err error
}

func (e NotPrivilegedError) Unwrap() error {
	return e.err
}

func (e NotPrivilegedError) Error() string {
	return fmt.Sprint("l'utilisateur n'est pas admin: id = "+e.id, e.err)
}

type ProtectedUserError struct {
	id UserID
}

func (e ProtectedUserError) Error() string {
	return fmt.Sprintf("cet action est interdite sur cet utilisateur (%s)", e.id)
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
	return fmt.Sprint("erreur survenue lors de l'exécution de la requête : ", e.err)
}

func (e UnexpectedMongoError) Unwrap() error {
	return e.err
}

type UnexpectedMongoDecodeError struct {
	err error
}

func (e UnexpectedMongoDecodeError) Error() string {
	return fmt.Sprint("erreur survenue lors du décodage du résultat de la requête : ", e.err)
}

func (e UnexpectedMongoDecodeError) Unwrap() error {
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
	return fmt.Sprint("la connexion a échoué : ", e.err)
}

func (e UnreachableMongoError) Unwrap() error {
	return e.err
}

type InvalidMongoConfigurationError struct {
	err error
}

func (e InvalidMongoConfigurationError) Error() string {
	return fmt.Sprint("les paramètres de connexion sont invalides : ", e.err)
}

func (e InvalidMongoConfigurationError) Unwrap() error {
	return e.err
}

type ForbiddenOperationError struct {
	err error
}

func (e ForbiddenOperationError) Error() string {
	return fmt.Sprint("operation interdite : ", e.err)
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
