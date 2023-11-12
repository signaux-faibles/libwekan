package libwekan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CardID string

type Poker struct {
	Question             bool     `bson:"question" json:"question,omitempty"`
	One                  []string `bson:"one" json:"one,omitempty"`
	Two                  []string `bson:"two" json:"two,omitempty"`
	Three                []string `bson:"three" json:"three,omitempty"`
	Five                 []string `bson:"five" json:"five,omitempty"`
	Eight                []string `bson:"eight" json:"eight,omitempty"`
	Thirteen             []string `bson:"thirteen" json:"thirteen,omitempty"`
	Twenty               []string `bson:"twenty" json:"twenty,omitempty"`
	Forty                []string `bson:"forty" json:"forty,omitempty"`
	OneHundred           []string `bson:"oneHundred" json:"oneHundred,omitempty"`
	Unsure               []string `bson:"unsure" json:"unsure,omitempty"`
	End                  *bool    `bson:"end" json:"end,omitempty"`
	AllowNonBoardMembers bool     `bson:"allowNonBoardMembers" json:"allowNonBoardMembers,omitempty"`
}

type Vote struct {
	Question             string   `bson:"question" json:"question,omitempty"`
	Positive             []string `bson:"positive" json:"positive,omitempty"`
	Negative             []string `bson:"negative" json:"negative,omitempty"`
	End                  *bool    `bson:"end" json:"end,omitempty"`
	Public               bool     `bson:"public" json:"public,omitempty"`
	AllowNonBoardMembers bool     `bson:"allowNonBoardMembers" json:"allowNonBoardMembers,omitempty"`
}

type CardCustomFieldID string

type CardCustomField struct {
	ID    CardCustomFieldID `bson:"_id" json:"_id,omitempty"`
	Value string            `bson:"value" json:"value,omitempty"`
}

type Card struct {
	ID               CardID            `bson:"_id" json:"_id,omitempty"`
	Title            string            `bson:"title" json:"title,omitempty"`
	Members          []UserID          `bson:"members" json:"members,omitempty"`
	LabelIDs         []BoardLabelID    `bson:"labelIds" json:"labelIds,omitempty"`
	CustomFields     []CardCustomField `bson:"customFields" json:"customFields,omitempty"`
	ListID           ListID            `bson:"listId" json:"listId,omitempty"`
	BoardID          BoardID           `bson:"boardId" json:"boardId,omitempty"`
	Sort             float64           `bson:"sort" json:"sort,omitempty"`
	SwimlaneID       SwimlaneID        `bson:"swimlaneId" json:"swimlaneId,omitempty"`
	Type             string            `bson:"type" json:"type,omitempty"`
	Archived         bool              `bson:"archived" json:"archived,omitempty"`
	ParentID         CardID            `bson:"parentId,omitempty" json:"parentId,omitempty"`
	CoverID          string            `bson:"coverId" json:"coverId,omitempty"`
	CreatedAt        time.Time         `bson:"createdAt" json:"createdAt,omitempty"`
	ModifiedAt       time.Time         `bson:"modifiedAt" json:"modifiedAt,omitempty"`
	DateLastActivity time.Time         `bson:"dateLastActivity" json:"dateLastActivity,omitempty"`
	Description      string            `bson:"description" json:"description,omitempty"`
	RequestedBy      UserID            `bson:"requestedBy" json:"requestedBy,omitempty"`
	AssignedBy       UserID            `bson:"assignedBy" json:"assignedBy,omitempty"`
	Assignees        []UserID          `bson:"assignees" json:"assignees,omitempty"`
	SpentTime        int               `bson:"spentTime" json:"spentTime,omitempty"`
	IsOverTime       bool              `bson:"isOvertime" json:"isOvertime,omitempty"`
	UserID           UserID            `bson:"userId" json:"userId,omitempty"`
	SubtaskSort      int               `bson:"subtaskSort" json:"subtaskSort,omitempty"`
	LinkedID         CardID            `bson:"linkedId" json:"linkedId,omitempty"`
	Vote             Vote              `bson:"vote" json:"vote,omitempty"`
	Poker            Poker             `bson:"poker" json:"poker,omitempty"`
	TargetIDGantt    []string          `bson:"targetId_gantt" json:"targetId_gantt,omitempty"`
	LinkTypeGantt    []string          `bson:"linkType_gantt" json:"linkType_gantt,omitempty"`
	LinkIDGantt      []string          `bson:"linkId_gantt" json:"linkId_gantt,omitempty"`
	StartAt          time.Time         `bson:"startAt" json:"startAt,omitempty"`
	EndAt            *time.Time        `bson:"endAt" json:"endAt,omitempty"`
}

type CardWithComments struct {
	Card     Card      `bson:"card" json:"card,omitempty"`
	Comments []Comment `bson:"comments" json:"comments,omitempty"`
}

func BuildCard(boardID BoardID, listID ListID, swimlaneID SwimlaneID, title string, description string, userID UserID) Card {
	return Card{
		ID:               CardID(newId()),
		Title:            title,
		ListID:           listID,
		BoardID:          boardID,
		SwimlaneID:       swimlaneID,
		Members:          []UserID{},
		LabelIDs:         []BoardLabelID{},
		Type:             "card",
		CreatedAt:        toMongoTime(time.Now()),
		ModifiedAt:       toMongoTime(time.Now()),
		DateLastActivity: toMongoTime(time.Now()),
		Description:      description,
		UserID:           userID,
		TargetIDGantt:    []string{},
		LinkTypeGantt:    []string{},
		LinkIDGantt:      []string{},
		StartAt:          toMongoTime(time.Now()),
	}
}

func (cardID CardID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetCardFromID(ctx, cardID)
	return err
}

func (cardID CardID) GetDocument(ctx context.Context, wekan *Wekan) (Card, error) {
	return wekan.GetCardFromID(ctx, cardID)
}

func (card *Card) AddMember(memberID UserID) {
	if !(contains(card.Members, memberID)) {
		card.Members = append(card.Members, memberID)
	}
}

func (wekan *Wekan) SelectCardsFromQuery(ctx context.Context, query bson.M) ([]Card, error) {
	var cards []Card
	cur, err := wekan.db.Collection("cards").Find(ctx, query)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	err = cur.All(ctx, &cards)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return cards, nil
}

func (wekan *Wekan) SelectCardsFromUserID(ctx context.Context, userID UserID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"userId": userID})
}

func (wekan *Wekan) SelectCardsFromMemberID(ctx context.Context, userID UserID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"members": userID})
}

func (wekan *Wekan) SelectCardsFromBoardID(ctx context.Context, boardID BoardID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"boardId": boardID})
}

func (wekan *Wekan) SelectCardsFromSwimlaneID(ctx context.Context, swimlaneID SwimlaneID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"swimlaneId": swimlaneID})
}

func (wekan *Wekan) SelectCardsFromListID(ctx context.Context, listID ListID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"listId": listID})
}

// SelectCardsFromPipeline retourne les objets correspondant au modèle Card à partir d'un pipeline mongodb
func (wekan *Wekan) SelectCardsFromPipeline(ctx context.Context, collection string, pipeline Pipeline) ([]Card, error) {
	cur, err := wekan.db.Collection(collection).Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err)
		return nil, UnexpectedMongoError{err}
	}
	var cards []Card
	err = cur.All(ctx, &cards)
	if err != nil {
		return nil, UnexpectedMongoDecodeError{err}
	}
	return cards, nil
}

// SelectCardsWithCommentsFromPipeline retourne les objets correspondant au modèle CardWithComments à partir d'un pipeline mongodb
func (wekan *Wekan) SelectCardsWithCommentsFromPipeline(ctx context.Context, collection string, pipeline Pipeline) ([]CardWithComments, error) {
	cur, err := wekan.db.Collection(collection).Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err)
		return nil, UnexpectedMongoError{err}
	}
	var cards []CardWithComments
	err = cur.All(ctx, &cards)
	if err != nil {
		return nil, UnexpectedMongoDecodeError{err}
	}
	return cards, nil
}

func (wekan *Wekan) SelectCardsFromCustomTextField(ctx context.Context, name string, value string) ([]Card, error) {
	pipeline := wekan.BuildCardFromCustomTextFieldPipeline(name, value)
	return wekan.SelectCardsFromPipeline(ctx, "customFields", pipeline)
}

func (wekan *Wekan) GetCardFromID(ctx context.Context, cardID CardID) (Card, error) {
	cards, err := wekan.SelectCardsFromQuery(ctx, bson.M{"_id": cardID})
	if err != nil {
		return Card{}, UnexpectedMongoError{err}
	}
	if len(cards) == 0 {
		return Card{}, CardNotFoundError{cardID}
	}
	if len(cards) > 1 {
		return Card{}, UnexpectedMongoError{errors.New("erreur fatale, cette requête ne peut retourner qu'un objet")}
	}
	return cards[0], nil
}

func (wekan *Wekan) InsertCard(ctx context.Context, card Card) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	if err := wekan.CheckDocuments(
		ctx,
		card.BoardID,
		card.ListID,
		card.SwimlaneID,
	); err != nil {
		return err
	}
	if _, err := wekan.db.Collection("cards").InsertOne(ctx, card); err != nil {
		return UnexpectedMongoError{err}
	}

	activity, err := wekan.newActivityCreateCardFromCard(ctx, card)
	if err != nil {
		return err
	}
	_, err = wekan.insertActivity(ctx, activity)
	return err
}

func (wekan *Wekan) AddLabelToCard(ctx context.Context, cardID CardID, labelID BoardLabelID) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	card, err := cardID.GetDocument(ctx, wekan)
	if err != nil {
		return err
	}
	board, err := card.BoardID.GetDocument(ctx, wekan)
	if err != nil {
		return err
	}

	label := board.GetLabelByID(labelID)
	if label == (BoardLabel{}) {
		return BoardLabelNotFoundError{labelID, board}
	}
	stats, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{"_id": cardID},
		bson.M{
			"$addToSet": bson.M{
				"labelIds": labelID,
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

func (wekan *Wekan) BuildDomainCardsPipeline() Pipeline {
	matchBoardsStage := bson.M{
		"$match": bson.M{
			"slug": bson.M{"$regex": wekan.slugDomainRegexp},
		},
	}

	lookupCardsStage := bson.M{
		"$lookup": bson.M{
			"from": "cards",
			"let":  bson.M{"boardId": "$_id"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$boardId", "$$boardId"}}}},
			},
			"as": "card",
		},
	}

	unwindCardStage := bson.M{
		"$unwind": "$card",
	}

	replaceRootStage := bson.M{
		"$replaceRoot": bson.M{
			"newRoot": "$card",
		},
	}

	return Pipeline{
		matchBoardsStage,
		lookupCardsStage,
		unwindCardStage,
		replaceRootStage,
	}
}

func (wekan *Wekan) BuildCardFromCustomTextFieldPipeline(name string, value string) Pipeline {
	return wekan.BuildCardFromCustomTextFieldsPipeline(name, []string{value})
}

func (wekan *Wekan) BuildCardFromCustomTextFieldsPipeline(name string, values []string) Pipeline {
	matchNameStage := bson.M{
		"$match": bson.M{
			"name": name,
		}}
	unwindBoardIdsStage := bson.M{
		"$unwind": "$boardIds",
	}

	lookupBoardsPipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"$expr": bson.M{
					"$eq": bson.A{"$_id", "$$boardId"},
				},
			},
		}}

	lookupBoardsStage := bson.M{
		"$lookup": bson.M{
			"from": "boards",
			"let": bson.M{
				"boardId": "$boardIds",
			},
			"pipeline": lookupBoardsPipeline,
			"as":       "board",
		},
	}

	unwindBoardsStage := bson.M{
		"$unwind": "$board",
	}

	matchBoardsStage := bson.M{
		"$match": bson.M{
			"board.slug": bson.M{
				"$regex": primitive.Regex{
					Pattern: wekan.slugDomainRegexp,
					Options: "i",
				}}}}

	matchCardsBoardIds := bson.M{
		"$match": bson.M{
			"$expr": bson.M{
				"$eq": bson.A{"$boardId", "$$boardId"},
			},
		},
	}

	duplicateCardsCustomFields := bson.M{
		"$addFields": bson.M{"customField": "$customFields"},
	}

	unwindCardsCustomField := bson.M{
		"$unwind": bson.M{
			"path":                       "$customField",
			"preserveNullAndEmptyArrays": true,
		},
	}

	matchCardsCustomField := bson.M{
		"$match": bson.M{
			"$expr": bson.M{
				"$and": bson.A{
					bson.M{
						"$eq": bson.A{"$customField._id", "$$customFieldId"},
					},
					bson.M{
						"$in": bson.A{"$customField.value", values},
					},
				},
			},
		},
	}

	removeCardsCustomField := bson.M{
		"$project": bson.M{
			"customField": false,
		},
	}

	lookupCardsPipeline := Pipeline{
		matchCardsBoardIds,
		duplicateCardsCustomFields,
		unwindCardsCustomField,
		matchCardsCustomField,
		removeCardsCustomField,
	}

	lookupCardsStage := bson.M{
		"$lookup": bson.M{
			"let": bson.M{
				"boardId":       "$boardIds",
				"customFieldId": "$_id",
			},
			"from":     "cards",
			"pipeline": lookupCardsPipeline,
			"as":       "cards",
		},
	}

	unwindCardsStage := bson.M{
		"$unwind": "$cards",
	}

	replaceRootStage := bson.M{
		"$replaceRoot": bson.M{
			"newRoot": "$cards",
		},
	}

	pipeline := Pipeline{
		matchNameStage,
		unwindBoardIdsStage,
		lookupBoardsStage,
		unwindBoardsStage,
		matchBoardsStage,
		lookupCardsStage,
		unwindCardsStage,
		replaceRootStage,
	}

	return pipeline
}

func (wekan *Wekan) ArchiveCard(ctx context.Context, cardID CardID) error {
	update, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id":      cardID,
		"archived": false,
	}, bson.M{
		"$set": bson.M{
			"archived": true,
		},
		"$currentDate": bson.M{
			"modifiedAt":       true,
			"dateLastActivity": true,
		},
	})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	if update.MatchedCount == 0 {
		return CardNotFoundError{cardID}
	}
	if update.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	return nil
}

func (wekan *Wekan) UnarchiveCard(ctx context.Context, cardID CardID) error {
	update, err := wekan.db.Collection("cards").UpdateOne(ctx, bson.M{
		"_id":      cardID,
		"archived": true,
	}, bson.M{
		"$set": bson.M{
			"archived": false,
		},
		"$currentDate": bson.M{
			"modifiedAt":       true,
			"dateLastActivity": true,
		},
	})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	if update.MatchedCount == 0 {
		return CardNotFoundError{cardID}
	}
	if update.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	return nil
}

func (config *Config) CustomFieldWithName(card Card, name string) string {
	configBoard := config.Boards[card.BoardID]
	for _, customField := range card.CustomFields {
		if configBoard.CustomFields[customField.ID].Name == name {
			return customField.Value
		}
	}
	return ""
}

// TODO: écrire un test
func (wekan *Wekan) UpdateCardDescription(ctx context.Context, cardID CardID, description string) error {
	stats, err := wekan.db.Collection("cards").UpdateOne(ctx,
		bson.M{"_id": cardID},
		bson.M{
			"$set": bson.M{
				"description": description,
			},
			"$currentDate": bson.M{
				"modifiedAt":       true,
				"dateLastActivity": true,
			},
		},
	)
	if stats.MatchedCount == 0 {
		return CardNotFoundError{cardID: cardID}
	}
	if stats.ModifiedCount == 0 {
		return NothingDoneError{}
	}
	if err != nil {
		return UnexpectedMongoError{err: err}
	}
	return nil
}

func (wekan *Wekan) EnsureMoveCardList(ctx context.Context, cardID CardID, listID ListID, userID UserID) error {
	card, err := cardID.GetDocument(ctx, wekan)
	if err != nil {
		return err
	}
	// si la liste est déjà set, rien à faire
	if card.ListID == listID {
		return nil
	}

	// si la liste n'est pas dans cette board, on retourne une erreur
	lists, err := wekan.SelectListsFromBoardID(ctx, card.BoardID)
	listsIDs := mapSlice(lists, func(list List) ListID { return list.ID })
	if !contains(listsIDs, listID) {
		return ListNotFoundError{listID: listID}
	}

	// pas besoin de vérifier les stats, nous savons déjà que la liste est différente
	_, err = wekan.db.Collection("cards").UpdateOne(ctx, bson.M{"_id": cardID}, bson.M{"$set": bson.M{"listId": listID}})
	if err != nil {
		return UnexpectedMongoError{err}
	}

	// insertion de l'activité
	activity, err := wekan.newActivityMoveCardFromMovedCard(ctx, card, userID)
	if err != nil {
		return err
	}
	_, err = wekan.insertActivity(ctx, activity)
	return err
}
