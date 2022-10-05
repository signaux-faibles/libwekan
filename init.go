package libwekan

import (
	"context"

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
	privileged       bool
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

	if err != nil {
		return Wekan{}, err
	}

	return w, nil
}

func (wekan *Wekan) AssertHasAdmin(ctx context.Context) error {
	//if !wekan.IsPrivileged() {
	//	return false
	//}
	admin, err := wekan.GetUserFromUsername(ctx, wekan.adminUsername)
	if err != nil {
		return err
	}
	if !admin.IsAdmin {
		return UserIsNotAdminError{admin.ID}
	}
	wekan.adminUserID = admin.ID
	return nil
}

func (wekan *Wekan) AdminUsername() Username {
	return wekan.adminUsername
}

func (wekan *Wekan) AdminID() UserID {
	return wekan.adminUserID
}

func (wekan *Wekan) IsPrivileged() bool {
	return wekan.privileged
}

//func (wekan *Wekan) AdminUser(ctx context.Context) (User, error) {
//	admin, err := wekan.GetUserFromUsername(ctx, wekan.adminUsername)
//	if _, ok := err.(UnknownUserError); ok {
//		return User{}, err
//	}
//	if !admin.IsAdmin {
//		return User{}, UserIsNotAdminError{admin.ID}
//	}
//	return admin, nil
//}

//func (wekan *Wekan) CheckAdminUserIsAdmin(ctx context.Context) error {
//	if wekan.adminUserID != "" {
//		return nil
//	}
//	adminUser, err := wekan.AdminUser(ctx)
//	if err != nil {
//		return err
//	}
//	wekan.adminUserID = adminUser.ID
//	return nil
//}

func (wekan *Wekan) Ping(ctx context.Context) error {
	err := wekan.client.Ping(ctx, nil)
	if err != nil {
		return UnreachableMongoError{err}
	}
	return nil
}

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
