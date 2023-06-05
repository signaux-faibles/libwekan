package libwekan

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestErrors_UserAlreadyExistsError(t *testing.T) {
	e := UserAlreadyExistsError{User{ID: "testID", Username: "testID"}}
	expected := fmt.Sprintf("l'utilisateur existe déjà (UserID: %s, Username: %s)", e.user.ID, e.user.Username)
	assert.EqualError(t, e, expected)
}

func TestErrors_UserNotFoundError(t *testing.T) {
	e := UserNotFoundError{key: "test"}
	expected := fmt.Sprintf("l'utilisateur n'est pas connu (%s)", e.key)
	assert.EqualError(t, e, expected)
}

func TestErrors_BoardNotFoundError(t *testing.T) {
	e := boardNotFoundWithSlug("test")
	expected := "aucun tableau n'a été trouvé avec le slug : 'test'"
	assert.EqualError(t, e, expected)
	e = boardNotFoundWithId("test")
	expected = "aucun tableau n'a été trouvé avec le boardID : 'test'"
	assert.EqualError(t, e, expected)
	e = boardNotFoundWithTitle("test")
	expected = "aucun tableau n'a été trouvé avec le titre : 'test'"
	assert.EqualError(t, e, expected)
}

func TestErrors_NotPrivilegedError(t *testing.T) {
	e := NotPrivilegedError{"test", errors.New("")}
	expected := fmt.Sprintf("l'utilisateur n'est pas admin: id = %s", e.id)
	assert.EqualError(t, e, expected)
}

func TestErrors_ProtectedUserError(t *testing.T) {
	e := ProtectedUserError{"test"}
	expected := fmt.Sprintf("cet action est interdite sur cet utilisateur (%s)", e.id)
	assert.EqualError(t, e, expected)
}

func TestErrors_InsertEmptyRuleError(t *testing.T) {
	e := InsertEmptyRuleError{}
	expected := "l'insertion d'un objet Rule vide est impossible"
	assert.EqualError(t, e, expected)
	assert.Error(t, e)
}

func TestErrors_BoardLabelAlreadyExistsError(t *testing.T) {
	e := BoardLabelAlreadyExistsError{BoardLabel{}, Board{}}
	expected := fmt.Sprintf("un objet BoardLabel existe déjà dans la board (%s) avec le même nom (%s)", e.board.ID, e.boardLabel.Name)
	assert.EqualError(t, e, expected)
}

func TestErrors_UnexpectedMongoError(t *testing.T) {
	e := UnexpectedMongoError{err: mongo.ErrNoDocuments}
	expected := "erreur survenue lors de l'exécution de la requête"
	assert.ErrorContains(t, e, expected)
	assert.ErrorIs(t, e, mongo.ErrNoDocuments)
}

func TestErrors_UnexpectedMongoDecodeError(t *testing.T) {
	e := UnexpectedMongoDecodeError{err: mongo.ErrNoDocuments}
	expected := "erreur survenue lors du décodage du résultat de la requête"
	assert.ErrorContains(t, e, expected)
	assert.ErrorIs(t, e, mongo.ErrNoDocuments)
}

func TestErrors_AlreadySetActityError(t *testing.T) {
	e := AlreadySetActivityError{"test"}
	expected := fmt.Sprintf("l'activité est déjà définie: activityType = %s", e.activityType)
	assert.EqualError(t, e, expected)
}

func TestErrors_UnreachableMongoError(t *testing.T) {
	e := UnreachableMongoError{mongo.ErrNilValue}
	expected := "la connexion a échoué"
	assert.ErrorContains(t, e, expected)
	assert.ErrorIs(t, e, mongo.ErrNilValue)
}

func TestErrors_InvalidMongoConfigurationError(t *testing.T) {
	e := InvalidMongoConfigurationError{mongo.ErrNilValue}
	expected := "les paramètres de connexion sont invalides"
	assert.ErrorContains(t, e, expected)
	assert.ErrorIs(t, e, mongo.ErrNilValue)
}

func TestErrors_ForbiddenOperationError(t *testing.T) {
	e := ForbiddenOperationError{UserIsNotMemberError{"test"}}
	expected := "operation interdite"
	assert.ErrorContains(t, e, expected)
	assert.ErrorIs(t, e, UserIsNotMemberError{"test"})
}

func TestErrors_NotImplemented(t *testing.T) {
	e := NotImplemented{"test"}
	expected := fmt.Sprintf("not implemented : " + e.method)
	assert.EqualError(t, e, expected)
}

func TestErrors_CardNotFoundError(t *testing.T) {
	e := CardNotFoundError{"test"}
	expected := fmt.Sprintf("la carte n'existe pas (ID: %s)", e.cardID)
	assert.EqualError(t, e, expected)
}

func TestErrors_RuleNotFoundError(t *testing.T) {
	e := RuleNotFoundError{"test"}
	expected := fmt.Sprintf("la règle n'existe pas (ID: %s)", e.ruleID)
	assert.EqualError(t, e, expected)
}

func TestErrors_ActionNotFoundError(t *testing.T) {
	e := ActionNotFoundError{"test"}
	expected := fmt.Sprintf("l'action n'existe pas (ID: %s)", e.actionId)
	assert.EqualError(t, e, expected)
}

func TestErrors_TriggerNotFoundError(t *testing.T) {
	e := TriggerNotFoundError{"test"}
	expected := fmt.Sprintf("le trigger n'existe pas (ID: %s)", e.triggerId)
	assert.EqualError(t, e, expected)
}

func TestErrors_NothingDoneError(t *testing.T) {
	e := NothingDoneError{}
	expected := "le traitement n'a eu aucun effet"
	assert.EqualError(t, e, expected)
}
func TestErrors_UnknownActivityError(t *testing.T) {
	e := UnknownActivityError{"test"}
	expected := fmt.Sprintf("l'activité n'existe pas (ID: %s)", e.key)
	assert.EqualError(t, e, expected)
}
func TestErrors_UserIsNotMemberError(t *testing.T) {
	e := UserIsNotMemberError{"test"}
	expected := fmt.Sprintf("cette action nécessite que l'utilisateur soit membre (id=%s)", e.userID)
	assert.EqualError(t, e, expected)
}
