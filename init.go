package libwekan

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Wekan struct {
	url              string
	databaseName     string
	client           *mongo.Client
	db               *mongo.Database
	adminUsername    Username
	adminUserID      UserID
	privileged       *bool
	slugDomainRegexp string
}

// Init retourne un objet de type `Wekan`
func Init(ctx context.Context, uri string, databaseName string, adminUsername Username, slugDomainRegexp string) (Wekan, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return Wekan{}, InvalidMongoConfigurationError{err}
	}
	w := Wekan{
		url:              uri,
		databaseName:     databaseName,
		client:           client,
		db:               client.Database(databaseName),
		adminUsername:    adminUsername,
		slugDomainRegexp: slugDomainRegexp,
	}
	return w, nil
}

func (wekan *Wekan) Ping(ctx context.Context) error {
	return wekan.client.Ping(ctx, nil)
}

// AssertPrivileged s'assure que l'utilisateur déclaré dans la propriété
// Wekan.adminUsername est bien un utilisateur admin dans la base de données
func (wekan *Wekan) AssertPrivileged(ctx context.Context) error {
	if wekan.privileged != nil {
		if *wekan.privileged {
			return nil
		}
		return NotPrivilegedError{wekan.adminUserID, errors.New("l'utilisateur n'est pas administrateur")}
	}
	admin, err := wekan.GetUserFromUsername(ctx, wekan.adminUsername)
	if err != nil {
		return NotPrivilegedError{"inconnu", err}
	}
	wekan.adminUserID = admin.ID
	wekan.privileged = &admin.IsAdmin
	//if !admin.IsAdmin {
	//	return ForbiddenOperationError{
	//		NotPrivilegedError{admin.ID, errors.New("L'utilisateur n'est pas administration")},
	//	}
	//}
	return nil
}

func (wekan *Wekan) AdminUsername() Username {
	return wekan.adminUsername
}

func (wekan *Wekan) AdminID() UserID {
	return wekan.adminUserID
}

//
//func (wekan *Wekan) IsPrivileged() bool {
//	return *wekan.privileged
//}

type Document interface {
	Check(context.Context, *Wekan) error
}

func (wekan *Wekan) CheckDocuments(ctx context.Context, documents ...Document) error {
	for _, document := range documents {
		if err := document.Check(ctx, wekan); err != nil {
			return err
		}
	}
	return nil
}
