package libwekan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Users []User

type User struct {
	ID                   UserID       `bson:"_id"`
	CreateAt             time.Time    `bson:"createAt"`
	Services             UserServices `bson:"services"`
	Username             Username     `bson:"username"`
	Emails               []UserEmail  `bson:"emails"`
	Profile              UserProfile  `bson:"profile"`
	AuthenticationMethod string       `bson:"authenticationMethod"`
	ModifiedAt           time.Time    `bson:"modifiedAt"`
	IsAdmin              bool         `bson:"isAdmin"`
	LoginDisabled        bool         `bson:"loginDisabled"`
}

type UserTemplates struct {
	TemplateBoard         Board
	CardTemplateSwimlane  Swimlane
	ListTemplateSwimlane  Swimlane
	BoardTemplateSwimlane Swimlane
}
type UserServicesOIDC struct {
	ID           string   `bson:"id"`
	Username     Username `bson:"username"`
	Fullname     string   `bson:"fullname"`
	AccessToken  string   `bson:"accessToken"`
	ExpiresAt    int      `bson:"expiresAt"`
	Email        string   `bson:"email"`
	RefreshToken string   `bson:"refreshToken"`
}

type UserServicesResume struct {
	LoginTokens []UserServicesResumeLoginToken
}

type UserServicesResumeLoginToken struct {
	When        time.Time `bson:"when"`
	HashedToken string    `bson:"hashedToken"`
}

type UserServicesPassword struct {
	Bcrypt string `bson:"bcrypt"`
}

type UserServices struct {
	OIDC     UserServicesOIDC     `bson:"oidc"`
	Resume   UserServicesResume   `bson:"resume"`
	Password UserServicesPassword `bson:"password"`
}

type UserEmail struct {
	Address  string `json:"address"`
	Verified bool
}

type UserProfileNotification struct {
	Activity string `bson:"activity"`
}

type UserProfile struct {
	Initials                 string                    `bson:"initials"`
	Fullname                 string                    `bson:"fullname"`
	BoardView                string                    `bson:"boardView"`
	ListSortBy               string                    `bson:"-modifiedAt"`
	TemplatesBoardId         BoardID                   `bson:"templatesBoardId"`
	CardTemplatesSwimlaneId  SwimlaneID                `bson:"cardTemplatesSwimlaneId"`
	ListTemplatesSwimlaneId  SwimlaneID                `bson:"listTemplatesSwimlaneId"`
	BoardTemplatesSwimlaneId SwimlaneID                `bson:"boardTemplatesSwimlaneId"`
	InvitedBoards            []string                  `bson:"invitedBoards"`
	StarredBoards            []string                  `bson:"starredBoards"`
	Language                 string                    `bson:"language"`
	CardMaximized            bool                      `bson:"cardMaximized"`
	EmailBuffer              []string                  `bson:"emailBuffer"`
	Notifications            []UserProfileNotification `bson:"notifications"`
	HiddenSystemMessages     bool                      `bson:"hiddenSystemMessages"`
}

type Username string
type UserID string

func (u Username) toString() string {
	return string(u)
}

func (u UserID) toString() string {
	return string(u)
}

func (u User) getUsername() Username {
	return u.Username
}

func (u User) getID() UserID {
	return u.ID
}

// ListUsers returns all wekan users
func (wekan *Wekan) ListUsers(ctx context.Context) ([]User, error) {
	cursor, err := wekan.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []User
	err = cursor.All(ctx, &users)
	return users, err
}

// GetUserFromUsername retourne l'objet utilisateur correspond au champ .username
func (wekan *Wekan) GetUserFromUsername(ctx context.Context, username Username) (User, error) {
	var user User
	err := wekan.db.Collection("users").FindOne(ctx, bson.M{
		"username": username,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, UnknownUserError{key: string("username = " + username)}
		}
		return User{}, UnexpectedMongoError{err}
	}
	return user, nil
}

// GetUserFromID retourne l'objet utilisateur correspond au champ ._id
func (wekan *Wekan) GetUserFromID(ctx context.Context, id UserID) (User, error) {
	var user User
	err := wekan.db.Collection("users").FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, UnknownUserError{key: string("id = " + id)}
		}
		return User{}, UnexpectedMongoError{err}
	}
	return user, nil
}

// GetUsersFromUsernames retourne les objets users correspondant aux usernames en une seule requête
func (wekan *Wekan) GetUsersFromUsernames(ctx context.Context, usernames []Username) ([]User, error) {
	usernameSet := uniq(usernames)
	cur, err := wekan.db.Collection("users").Find(ctx, bson.M{
		"username": bson.M{"$in": usernameSet},
	})
	if err != nil {
		return Users{}, UnexpectedMongoError{err}
	}
	var users []User
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, UnexpectedMongoError{err}
		}
		users = append(users, user)
	}
	if len(users) != len(usernameSet) {
		selectedUsernamesString := mapSlice(mapSlice(users, User.getUsername), Username.toString)
		usernameSetString := mapSlice(usernameSet, Username.toString)
		_, missing, _ := intersect(usernameSetString, selectedUsernamesString)
		sort.Strings(missing)
		return Users{}, UnknownUserError{key: fmt.Sprintf("usernames in (%s)", strings.Join(missing, ", "))}
	}
	return users, nil
}

// GetUsersFromIDs retourne les objets users correspondant aux usernames en une seule requête
func (wekan *Wekan) GetUsersFromIDs(ctx context.Context, userIDs []UserID) ([]User, error) {
	userIDSet := uniq(userIDs)
	cur, err := wekan.db.Collection("users").Find(ctx, bson.M{
		"_id": bson.M{"$in": userIDSet},
	})
	if err != nil {
		return Users{}, UnexpectedMongoError{err}
	}
	var users []User
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, UnexpectedMongoError{err}
		}
		users = append(users, user)
	}
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, UnexpectedMongoError{err}
		}
		users = append(users, user)
	}
	if len(users) != len(userIDSet) {
		selectedUsernamesString := mapSlice(mapSlice(users, User.getID), UserID.toString)
		userIDSetString := mapSlice(userIDSet, UserID.toString)
		_, missing, _ := intersect(userIDSetString, selectedUsernamesString)
		sort.Strings(missing)
		return Users{}, UnknownUserError{key: fmt.Sprintf("ids in (%s)", strings.Join(missing, ", "))}
	}
	return users, nil
}

// GetUsers retourne tous les utilisateurs
func (wekan *Wekan) GetUsers(ctx context.Context) (Users, error) {
	var users Users
	cursor, err := wekan.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (wekan *Wekan) UsernameExists(ctx context.Context, username Username) (bool, error) {
	_, err := wekan.GetUserFromUsername(ctx, username)
	if _, ok := err.(UnknownUserError); ok {
		return false, nil
	}
	return err == nil, err
}

func (wekan *Wekan) InsertUser(ctx context.Context, user User) (User, error) {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return User{}, err
	}

	userAlreadyExists, err := wekan.UsernameExists(ctx, user.Username)
	if err != nil || userAlreadyExists {
		return User{}, UserAlreadyExistsError{user}
	}
	if err = wekan.InsertTemplates(ctx, user.BuildTemplates()); err != nil {
		return User{}, err
	}
	if _, err = wekan.db.Collection("users").InsertOne(ctx, user); err != nil {
		return User{}, err
	}
	err = wekan.EnsureUserIsActiveBoardMember(ctx, user.Profile.TemplatesBoardId, user.ID)
	return user, err
}

func (wekan *Wekan) InsertTemplates(ctx context.Context, templates UserTemplates) error {
	if err := wekan.InsertSwimlane(ctx, templates.CardTemplateSwimlane); err != nil {
		return err
	}
	if err := wekan.InsertSwimlane(ctx, templates.ListTemplateSwimlane); err != nil {
		return err
	}
	if err := wekan.InsertSwimlane(ctx, templates.BoardTemplateSwimlane); err != nil {
		return err
	}
	if err := wekan.InsertBoard(ctx, templates.TemplateBoard); err != nil {
		return err
	}
	return nil
}

func (user *User) BuildTemplates() UserTemplates {
	templateBoard := newBoard("Template", "templates", "template-container")
	cardTemplateSwimlane := newCardTemplateSwimlane(templateBoard.ID)
	listTemplateSwimlane := newListTemplateSwimlane(templateBoard.ID)
	boardTemplateSwimlane := newBoardTemplateSwimlane(templateBoard.ID)

	user.Profile.TemplatesBoardId = templateBoard.ID
	user.Profile.CardTemplatesSwimlaneId = cardTemplateSwimlane.ID
	user.Profile.ListTemplatesSwimlaneId = listTemplateSwimlane.ID
	user.Profile.BoardTemplatesSwimlaneId = boardTemplateSwimlane.ID

	return UserTemplates{
		TemplateBoard:         templateBoard,
		CardTemplateSwimlane:  cardTemplateSwimlane,
		ListTemplateSwimlane:  listTemplateSwimlane,
		BoardTemplateSwimlane: boardTemplateSwimlane,
	}
}

// BuildUser retourne un objet User à insérer/updater avec la fonction Wekan.UpsertUser
func BuildUser(email, initials, fullname string) User {
	newUser := User{
		ID:       UserID(newId()),
		CreateAt: time.Now(),

		Services: UserServices{
			OIDC: UserServicesOIDC{
				ID:           email,
				Username:     Username(email),
				Fullname:     fullname,
				AccessToken:  "",
				ExpiresAt:    int(time.Now().UnixMilli()),
				Email:        email,
				RefreshToken: "",
			},
			Resume: UserServicesResume{
				LoginTokens: []UserServicesResumeLoginToken{},
			},
		},
		Username: Username(email),
		Emails: []UserEmail{
			{
				Address:  email,
				Verified: true,
			},
		},
		Profile: UserProfile{
			Initials:             initials,
			Fullname:             fullname,
			BoardView:            "board-view-swimlanes",
			ListSortBy:           "-modifiedAt",
			InvitedBoards:        []string{},
			StarredBoards:        []string{},
			EmailBuffer:          []string{},
			HiddenSystemMessages: true,
			Language:             "fr",
			Notifications:        []UserProfileNotification{},
			CardMaximized:        false,
		},
		AuthenticationMethod: "oauth2",
		ModifiedAt:           time.Now(),
		IsAdmin:              false,
		LoginDisabled:        false,
	}

	return newUser
}

// EnableUser: active un utilisateur dans la base `users` et active la participation à son tableau templates
func (wekan *Wekan) EnableUser(ctx context.Context, user User) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	_, err := wekan.db.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"loginDisabled": false,
		},
	})

	if err != nil {
		return err
	}

	// enable BoardMember on template board
	_, err = wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": user.Profile.TemplatesBoardId},
		bson.M{
			"$set": bson.M{"members.$[member].isActive": true},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": user.ID}}},
		},
	)
	return err
}

func (wekan *Wekan) CreateUsers(ctx context.Context, users Users) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	for _, user := range users {
		_, err := wekan.InsertUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wekan *Wekan) EnableUsers(ctx context.Context, users Users) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	for _, user := range users {
		err := wekan.EnableUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

// DisableUser désactive l'utilisateur dans la base `users` et désactive la participation à tous les tableaux
func (wekan *Wekan) DisableUser(ctx context.Context, user User) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	_, err := wekan.db.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"loginDisabled": true,
		},
	})

	if err != nil {
		return err
	}

	_, err = wekan.db.Collection("boards").UpdateMany(ctx, bson.M{},
		bson.M{
			"$set": bson.M{"members.$[member].isActive": false},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": user.ID}}},
		},
	)
	if err != nil {
		return err
	}

	return err
}

func (wekan *Wekan) DisableUsers(ctx context.Context, users Users) error {
	for _, user := range users {
		err := wekan.DisableUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}
