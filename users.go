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

func (userId UserID) GetDocument(ctx context.Context, wekan *Wekan) (User, error) {
	return wekan.GetUserFromID(ctx, userId)
}

func (userId UserID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetUserFromID(ctx, userId)
	return err
}

func (username Username) String() string {
	return string(username)
}

func (userId UserID) String() string {
	return string(userId)
}

func (user User) GetUsername() Username {
	return user.Username
}

func (user User) GetID() UserID {
	return user.ID
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
			return User{}, UserNotFoundError{key: string("username = " + username), err: err}
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
			return User{}, UserNotFoundError{key: string("id = " + id)}
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
		selectedUsernamesString := mapSlice(mapSlice(users, User.GetUsername), Username.String)
		usernameSetString := mapSlice(usernameSet, Username.String)
		_, missing, _ := intersect(usernameSetString, selectedUsernamesString)
		sort.Strings(missing)
		return Users{}, UserNotFoundError{key: fmt.Sprintf("usernames in (%s)", strings.Join(missing, ", "))}
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
		selectedUsernamesString := mapSlice(mapSlice(users, User.GetID), UserID.String)
		userIDSetString := mapSlice(userIDSet, UserID.String)
		_, missing, _ := intersect(userIDSetString, selectedUsernamesString)
		sort.Strings(missing)
		return Users{}, UserNotFoundError{key: fmt.Sprintf("ids in (%s)", strings.Join(missing, ", "))}
	}
	return users, nil
}

// GetUsers retourne tous les utilisateurs
func (wekan *Wekan) GetUsers(ctx context.Context) (Users, error) {
	var users Users
	cursor, err := wekan.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, UnexpectedMongoError{err}
		}
		users = append(users, user)
	}
	return users, nil
}

func (wekan *Wekan) UsernameExists(ctx context.Context, username Username) (bool, error) {
	_, err := wekan.GetUserFromUsername(ctx, username)
	if _, ok := err.(UserNotFoundError); ok {
		return false, nil
	}
	return err == nil, err
}

func (wekan *Wekan) InsertUser(ctx context.Context, user User) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	userAlreadyExists, err := wekan.UsernameExists(ctx, user.Username)
	if err != nil || userAlreadyExists {
		return UserAlreadyExistsError{user}
	}
	if err = wekan.InsertTemplates(ctx, user.BuildTemplates()); err != nil {
		return err
	}
	if _, err = wekan.db.Collection("users").InsertOne(ctx, user); err != nil {
		return UnexpectedMongoError{err}
	}
	_, err = wekan.EnsureUserIsActiveBoardMember(ctx, user.Profile.TemplatesBoardId, user.ID)
	return err
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
	templateBoard := BuildBoard("Template", "templates", "template-container")
	cardTemplateSwimlane := buildCardTemplateSwimlane(templateBoard.ID)
	listTemplateSwimlane := buildListTemplateSwimlane(templateBoard.ID)
	boardTemplateSwimlane := buildBoardTemplateSwimlane(templateBoard.ID)

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

func (user User) Admin(admin bool) User {
	user.IsAdmin = admin
	return user
}

// EnableUser active un utilisateur dans la base `users` et active la participation à son tableau templates
func (wekan *Wekan) EnableUser(ctx context.Context, user User) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	// enable BoardMember on template board
	_, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": user.Profile.TemplatesBoardId},
		bson.M{
			"$set": bson.M{"members.$[member].isActive": true},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": user.ID}}},
		},
	)
	if err != nil {
		return UnexpectedMongoError{err}
	}

	stats, err := wekan.db.Collection("users").UpdateOne(ctx,
		bson.M{
			"_id":           user.ID,
			"loginDisabled": true,
		},
		bson.M{
			"$set": bson.M{
				"loginDisabled": false,
			},
		})

	if err == mongo.ErrNoDocuments {
		return NothingDoneError{}
	}
	if err != nil {
		return UnexpectedMongoError{err}
	}
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	return nil
}

func (wekan *Wekan) InsertUsers(ctx context.Context, users Users) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	for _, user := range users {
		err := wekan.InsertUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wekan *Wekan) EnableUsers(ctx context.Context, users Users) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
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
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	// désactivation de l'utilisateur dans la collection users
	stats, err := wekan.db.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"loginDisabled": true,
		},
	})
	if err != nil {
		return UnexpectedMongoError{err}
	}

	// désactivation du BoardMember sur toutes les boards où il est présent
	boards, err := wekan.SelectBoardsFromMemberID(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, board := range boards {
		if err = wekan.DisableBoardMember(ctx, board.ID, user.ID); err != nil {
			return err
		}
	}

	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	return nil
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

func (wekan *Wekan) RemoveMemberFromCard(ctx context.Context, cardID CardID, memberID UserID) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	if err := wekan.CheckDocuments(ctx, cardID, memberID); err != nil {
		return err
	}

	stats, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id": cardID,
	}, bson.M{
		"$pull": bson.M{
			"members": memberID,
		},
	})
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) EnsureMemberOutOfCard(ctx context.Context, cardId CardID, memberID UserID) (bool, error) {
	err := wekan.RemoveMemberFromCard(ctx, cardId, memberID)
	if _, ok := err.(NothingDoneError); ok {
		return false, nil
	}
	return err == nil, err
}

func (wekan *Wekan) AddMemberToCard(ctx context.Context, cardID CardID, memberID UserID) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	card, err := wekan.GetCardFromID(ctx, cardID)
	if err != nil {
		return err
	}
	board, err := wekan.GetBoardFromID(ctx, card.BoardID)
	if err != nil {
		return err
	}
	user, err := wekan.GetUserFromID(ctx, memberID)
	if err != nil {
		return err
	}
	if !board.UserIsActiveMember(user) {
		return ForbiddenOperationError{
			UserIsNotMemberError{memberID},
		}
	}

	stats, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id": cardID,
	}, bson.M{
		"$addToSet": bson.M{
			"members": memberID,
		},
	})

	if err != nil {
		return UnexpectedMongoError{err}
	}
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	return nil
}

func (wekan *Wekan) EnsureMemberInCard(ctx context.Context, cardID CardID, memberID UserID) (bool, error) {
	err := wekan.AddMemberToCard(ctx, cardID, memberID)
	if _, ok := err.(NothingDoneError); ok {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
