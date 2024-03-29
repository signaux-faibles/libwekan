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
	ID                   UserID       `bson:"_id" json:"_id,omitempty"`
	CreatedAt            time.Time    `bson:"createAt" json:"createAt,omitempty"`
	Services             UserServices `bson:"services" json:"services,omitempty"`
	Username             Username     `bson:"username" json:"username,omitempty"`
	Emails               []UserEmail  `bson:"emails" json:"emails,omitempty"`
	Profile              UserProfile  `bson:"profile" json:"profile,omitempty"`
	AuthenticationMethod string       `bson:"authenticationMethod" json:"authenticationMethod,omitempty"`
	ModifiedAt           time.Time    `bson:"modifiedAt" json:"modifiedAt,omitempty"`
	IsAdmin              bool         `bson:"isAdmin" json:"isAdmin,omitempty"`
	LoginDisabled        bool         `bson:"loginDisabled" json:"loginDisabled,omitempty"`
}

type UserTemplates struct {
	TemplateBoard         Board
	CardTemplateSwimlane  Swimlane
	ListTemplateSwimlane  Swimlane
	BoardTemplateSwimlane Swimlane
}

type UserServicesOIDC struct {
	ID           string   `bson:"id" json:"id,omitempty"`
	Username     Username `bson:"username" json:"username,omitempty"`
	Fullname     string   `bson:"fullname" json:"fullname,omitempty"`
	AccessToken  string   `bson:"accessToken" json:"-"`
	ExpiresAt    int      `bson:"expiresAt" json:"-"`
	Email        string   `bson:"email" json:"email,omitempty"`
	RefreshToken string   `bson:"refreshToken" json:"-"`
}

type UserServicesResume struct {
	LoginTokens []UserServicesResumeLoginToken `json:"-"`
}

type UserServicesResumeLoginToken struct {
	When        time.Time `bson:"when" json:"when,omitempty"`
	HashedToken string    `bson:"hashedToken" json:"-"`
}

type UserServicesPassword struct {
	Bcrypt string `bson:"bcrypt" json:"-"`
}

type UserServices struct {
	OIDC     UserServicesOIDC     `bson:"oidc" json:"oidc,omitempty"`
	Resume   UserServicesResume   `bson:"resume" json:"resume,omitempty"`
	Password UserServicesPassword `bson:"password" json:"-"`
}

type UserEmail struct {
	Address  string `json:"address"`
	Verified bool   `json:"verified"`
}

type UserProfileNotification struct {
	Activity string `bson:"activity" json:"activity,omitempty"`
}

type UserProfile struct {
	Initials                 string                    `bson:"initials" json:"initials,omitempty"`
	Fullname                 string                    `bson:"fullname" json:"fullname,omitempty"`
	BoardView                string                    `bson:"boardView" json:"boardView,omitempty"`
	ListSortBy               string                    `bson:"-modifiedAt" json:"-"`
	TemplatesBoardId         BoardID                   `bson:"templatesBoardId" json:"templatesBoardId,omitempty"`
	CardTemplatesSwimlaneId  SwimlaneID                `bson:"cardTemplatesSwimlaneId" json:"cardTemplatesSwimlaneId,omitempty"`
	ListTemplatesSwimlaneId  SwimlaneID                `bson:"listTemplatesSwimlaneId" json:"listTemplatesSwimlaneId,omitempty"`
	BoardTemplatesSwimlaneId SwimlaneID                `bson:"boardTemplatesSwimlaneId" json:"boardTemplatesSwimlaneId,omitempty"`
	InvitedBoards            []string                  `bson:"invitedBoards" json:"invitedBoards,omitempty"`
	StarredBoards            []string                  `bson:"starredBoards" json:"starredBoards,omitempty"`
	Language                 string                    `bson:"language" json:"language,omitempty"`
	CardMaximized            bool                      `bson:"cardMaximized" json:"cardMaximized,omitempty"`
	EmailBuffer              []string                  `bson:"emailBuffer" json:"emailBuffer,omitempty"`
	Notifications            []UserProfileNotification `bson:"notifications" json:"notifications,omitempty"`
	HiddenSystemMessages     bool                      `bson:"hiddenSystemMessages" json:"hiddenSystemMessages,omitempty"`
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
	if len(userIDs) <= 0 {
		return Users{}, nil
	}
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
		ID:        UserID(newId()),
		CreatedAt: time.Now(),

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

func (wekan *Wekan) RemoveSelfMemberFromCard(ctx context.Context, card Card, member User) error {
	return wekan.RemoveMemberFromCard(ctx, card, member, member)
}

func (wekan *Wekan) RemoveMemberFromCard(ctx context.Context, card Card, user User, member User) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	if err := wekan.CheckDocuments(ctx, card.ID, member.ID); err != nil {
		return err
	}

	stats, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id": card.ID,
	}, bson.M{
		"$pull": bson.M{
			"members": member.ID,
		},
	})
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	if err != nil {
		return UnexpectedMongoError{err}
	}

	activity := newActivityCardUnjoinMember(user.ID, member.Username, member.ID, card.BoardID, card.ListID, card.ID, card.SwimlaneID)
	_, err = wekan.insertActivity(ctx, activity)
	if err != nil {
		return UnexpectedMongoError{err: err}
	}

	return nil
}

func (wekan *Wekan) RemoveAssigneeFromCard(ctx context.Context, card Card, user User, assignee User) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	if err := wekan.CheckDocuments(ctx, card.ID, assignee.ID); err != nil {
		return err
	}

	stats, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id":       card.ID,
		"assignees": bson.M{"$type": 4},
	}, bson.M{
		"$pull": bson.M{
			"assignees": assignee.ID,
		},
	})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}

	activity := newActivityCardUnjoinAssignee(user.ID, assignee.Username, assignee.ID, card.BoardID, card.ListID, card.ID, card.SwimlaneID)
	_, err = wekan.insertActivity(ctx, activity)
	if err != nil {
		return UnexpectedMongoError{err: err}
	}

	return nil
}

func (wekan *Wekan) EnsureAssigneeOutOfCard(ctx context.Context, card Card, user User, assignee User) (bool, error) {
	err := wekan.RemoveAssigneeFromCard(ctx, card, user, assignee)
	if _, ok := err.(NothingDoneError); ok {
		return false, nil
	}
	return err == nil, err
}

func (wekan *Wekan) EnsureMemberOutOfCard(ctx context.Context, card Card, user User, member User) (bool, error) {
	err := wekan.RemoveMemberFromCard(ctx, card, user, member)
	if _, ok := err.(NothingDoneError); ok {
		return false, nil
	}
	return err == nil, err
}

func (wekan *Wekan) AddSelfMemberToCard(ctx context.Context, card Card, member User) error {
	return wekan.AddMemberToCard(ctx, card, member, member)
}

func (wekan *Wekan) AddAssigneeToCard(ctx context.Context, card Card, user User, assignee User) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	board, err := wekan.GetBoardFromID(ctx, card.BoardID)
	if err != nil {
		return err
	}

	if !board.UserIsActiveMember(assignee) {
		return ForbiddenOperationError{
			UserIsNotMemberError{assignee.ID},
		}
	}

	if !board.UserIsActiveMember(user) {
		return ForbiddenOperationError{
			UserIsNotMemberError{user.ID},
		}
	}

	_, err = wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id":       card.ID,
		"assignees": nil,
	}, bson.M{
		"$set": bson.M{"assignees": bson.A{}},
	})

	if err != nil {
		return UnexpectedMongoError{err}
	}

	stats, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id": card.ID,
	}, bson.M{
		"$addToSet": bson.M{
			"assignees": assignee.ID,
		},
	})

	if err != nil {
		return UnexpectedMongoError{err}
	}
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}

	activity := newActivityCardJoinAssignee(user.ID, assignee.Username, assignee.ID, card.BoardID, card.ListID, card.ID, card.SwimlaneID)
	_, err = wekan.insertActivity(ctx, activity)
	if err != nil {
		return UnexpectedMongoError{err: err}
	}

	return nil
}

func (wekan *Wekan) AddMemberToCard(ctx context.Context, card Card, user User, member User) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	board, err := wekan.GetBoardFromID(ctx, card.BoardID)
	if err != nil {
		return err
	}

	if !board.UserIsActiveMember(member) {
		return ForbiddenOperationError{
			UserIsNotMemberError{member.ID},
		}
	}

	if !board.UserIsActiveMember(user) {
		return ForbiddenOperationError{
			UserIsNotMemberError{user.ID},
		}
	}

	stats, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id": card.ID,
	}, bson.M{
		"$addToSet": bson.M{
			"members": member.ID,
		},
	})

	if err != nil {
		return UnexpectedMongoError{err}
	}
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}

	activity := newActivityCardJoinMember(user.ID, member.Username, member.ID, card.BoardID, card.ListID, card.ID, card.SwimlaneID)
	_, err = wekan.insertActivity(ctx, activity)
	if err != nil {
		return UnexpectedMongoError{err: err}
	}

	return nil
}

func (wekan *Wekan) EnsureMemberInCard(ctx context.Context, card Card, user User, member User) (bool, error) {
	err := wekan.AddMemberToCard(ctx, card, user, member)
	if _, ok := err.(NothingDoneError); ok {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
