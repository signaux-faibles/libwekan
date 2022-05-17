package libwekan

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID                   string       `bson:"_id"`
	CreateAt             time.Time    `bson:"createAt"`
	Services             UserServices `bson:"services"`
	Username             string       `bson:"username"`
	Emails               UserEmails   `bson:"emails"`
	Profile              UserProfile  `bson:"profile"`
	AuthenticationMethod string       `bson:"authenticationMethod"`
	ModifiedAt           time.Time    `bson:"modifiedAt"`
	IsAdmin              bool         `bson:"isAdmin"`
	SessionData          []interface{}
}

type UserServices struct {
	OIDC struct {
		ID           string `bson:"id"`
		Username     string `bson:"username"`
		Fullname     string `bson:"fullname"`
		AccessToken  string `bson:"accessToken"`
		ExpiresAt    int    `bson:"expiresAt"`
		Email        string `bson:"email"`
		RefreshToken string `bson:"refreshToken"`
	} `bson:"oidc"`
	Resume struct {
		LoginTokens []struct {
			When        time.Time `bson:"when"`
			HashedToken string    `bson:"hashedToken"`
		}
	} `bson:"resume"`
	Password struct {
		Bcrypt string `bson:"bcrypt"`
	}
}

type UserEmails struct {
	Address  string `json:"address"`
	Verified bool
}

type UserProfile struct {
	Initials                 string   `bson:"initials"`
	Fullname                 string   `bson:"fullname"`
	BoardView                string   `bson:"boardView"`
	ListSortBy               string   `bson:"-modifiedAt"`
	TemplatesBoardId         string   `bson:"templatesBoardId"`
	CardTemplatesSwimlaneId  string   `bson:"cardTemplatesSwimlaneId"`
	ListTemplatesSwimlaneId  string   `bson:"listTemplatesSwimlaneId"`
	BoardTemplatesSwimlaneId string   `bson:"boardTemplatesSwimlaneId"`
	InvitedBoards            []string `bson:"invitedBoards"`
	StarredBoards            []string `bson:"starredBoards"`
	CardMaximized            bool     `bson:"cardMaximized"`
	EmailBuffer              []string `bson:"emailBuffer"`
	Notifications            []struct {
		Activity string `bson:"activity"`
	} `bson:"notifications"`
	HiddenSystemMessages bool `bson:"hiddenSystemMessages"`
}

func (w Wekan) GetUser(ctx context.Context, username string) (User, error) {
	var u User
	err := w.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&u)
	return u, err
}

// ListUsers returns all wekan users
func (w Wekan) ListUsers(ctx context.Context) ([]User, error) {
	cursor, err := w.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []User
	err = cursor.All(ctx, &users)
	return users, err
}

// GetUserFromUsername retourne l'objet utilisateur correspond au champ .username
func (w Wekan) GetUserFromUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := w.db.Collection("users").FindOne(ctx, bson.M{
		"username": username,
	}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUserFromID retourne l'objet utilisateur correspond au champ ._id
func (w Wekan) GetUserFromID(ctx context.Context, id string) (User, error) {
	var user User
	err := w.db.Collection("users").FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
