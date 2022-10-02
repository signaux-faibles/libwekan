package libwekan

import (
	"context"
	"time"
)

type CardID string

type Poker struct {
	Question             bool     `bson:"question"`
	One                  []string `bson:"one"`
	Two                  []string `bson:"two"`
	Three                []string `bson:"three"`
	Five                 []string `bson:"five"`
	Eight                []string `bson:"eight"`
	Thirteen             []string `bson:"thirteen"`
	Twenty               []string `bson:"twenty"`
	Forty                []string `bson:"forty"`
	OneHundred           []string `bson:"oneHundred"`
	Unsure               []string `bson:"unsure"`
	End                  *bool    `bson:"end"`
	AllowNonBoardMembers bool     `bson:"allowNonBoardMembers"`
}
type Vote struct {
	Question             string   `bson:"question"`
	Positive             []string `bson:"positive"`
	Negative             []string `bson:"negative"`
	End                  *bool    `bson:"end"`
	Public               bool     `bson:"public"`
	AllowNonBoardMembers bool     `bson:"allowNonBoardMembers"`
}

type CustomField struct {
	ID    string `json:"_id" bson:"_id"`
	Value string `json:"value" bson:"value"`
}

type Card struct {
	ID               CardID         `bson:"_id"`
	Title            string         `bson:"title"`
	Members          []UserID       `bson:"members"`
	LabelIDs         []BoardLabelID `bson:"labelIds"`
	CustomFields     []CustomField  `bson:"customFields"`
	ListID           ListID         `bson:"listId"`
	BoardID          BoardID        `bson:"boardId"`
	Sort             float64        `bson:"sort"`
	SwimlaneID       SwimlaneID     `bson:"swimlaneId"`
	Type             string         `bson:"type"`
	Archived         bool           `bson:"archived"`
	ParentID         CardID         `bson:"parentId"`
	CoverID          string         `bson:"coverId"`
	CreatedAt        time.Time      `bson:"createdAt"`
	ModifiedAt       time.Time      `bson:"modifiedAt"`
	DateLastActivity time.Time      `bson:"dateLastActivity"`
	Description      string         `bson:"description"`
	RequestedBy      UserID         `bson:"requestedBy"`
	AssignedBy       UserID         `bson:"assignedBy"`
	Assignees        []UserID       `bson:"assignees"`
	SpentTime        int            `bson:"spentTime"`
	IsOverTime       bool           `bson:"isOvertime"`
	UserID           UserID         `bson:"userId"`
	SubtaskSort      int            `bson:"subtaskSort"`
	LinkedID         CardID         `bson:"linkedId"`
	Vote             Vote           `bson:"vote"`
	Poker            Poker          `bson:"poker"`
	TargetIDGantt    []string       `bson:"targetId_gantt"`
	LinkTypeGantt    []string       `bson:"linkType_gantt"`
	LinkIDGantt      []string       `bson:"linkId_gantt"`
	StartAt          time.Time      `bson:"startAt"`
}

func (wekan *Wekan) SelectCardsFromUserID(ctx context.Context, userID UserID) ([]Card, error) {
	return nil, NotImplemented{}
}

func (wekan *Wekan) SelectCardsFromMemberID(ctx context.Context, userID UserID) ([]Card, error) {
	return nil, NotImplemented{}
}

func (wekan *Wekan) SelectCardsFromBoardID(ctx context.Context, boardID BoardID) ([]Card, error) {
	return nil, NotImplemented{}
}

func (wekan *Wekan) SelectCardsFromSwimlaneID(ctx context.Context, boardID BoardID) ([]Card, error) {
	return nil, NotImplemented{}
}

func (wekan *Wekan) SelectCardsFromListID(ctx context.Context, boardID BoardID) ([]Card, error) {
	return nil, NotImplemented{}
}

func (wekan *Wekan) GetCardFromID(ctx context.Context, cardID CardID) ([]Card, error) {
	return nil, NotImplemented{}
}

func (wekan *Wekan) InsertCard(ctx context.Context, card Card) error {
	return NotImplemented{}
}
