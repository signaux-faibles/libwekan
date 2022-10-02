package libwekan

import (
	"context"
	"fmt"
	"time"
)

type RuleID string
type Rule struct {
	ID         RuleID     `bson:"_id"`
	Title      string     `bson:"title"`
	TriggerID  *TriggerID `bson:"triggerId"`
	ActionID   *ActionID  `bson:"actionId"`
	BoardID    BoardID    `bson:"boardId"`
	CreatedAt  time.Time  `bson:"createdAt"`
	ModifiedAt time.Time  `bson:"modifiedAt"`
	Action     Action     `bson:"-"`
	Trigger    Trigger    `bson:"-"`
}

type TriggerID string
type Trigger struct {
	ID           TriggerID    `bson:"_id"`
	ActivityType string       `bson:"activityType"`
	BoardID      BoardID      `bson:"boardId"`
	LabelID      BoardLabelID `bson:"labelId"`
	Description  string       `bson:"desc"`
	UserID       UserID       `bson:"userId"`
	CreatedAt    time.Time    `bson:"createdAt"`
	ModifiedAt   time.Time    `bson:"modifiedAt"`
}

type ActionID string
type Action struct {
	ID          ActionID  `bson:"_id"`
	ActionType  string    `bson:"actionType"`
	Username    Username  `bson:"username"`
	BoardID     BoardID   `bson:"boardId"`
	Description string    `bson:"desc"`
	CreatedAt   time.Time `bson:"createdAt"`
	ModifiedAt  time.Time `bson:"modifiedAt"`
}

func (board Board) BuildTrigger(label BoardLabel) Trigger {
	return Trigger{
		ID:           TriggerID(newId()),
		ActivityType: "addedLabel",
		BoardID:      board.ID,
		LabelID:      label.ID,
		Description:  fmt.Sprintf("quand l'étiquette %s est ajoutée à la carte par *", label.Name),
		UserID:       "*",
		CreatedAt:    time.Now(),
		ModifiedAt:   time.Now(),
	}
}

func (board Board) BuildAction(username Username) Action {
	return Action{
		ID:          ActionID(newId()),
		ActionType:  "addMember",
		Username:    username,
		BoardID:     board.ID,
		Description: fmt.Sprintf("%s devient membre de la carte", username),
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}
}

func (board Board) BuildRule(user User, labelName BoardLabelName) Rule {
	label := board.GetLabelByName(labelName)
	if label == (BoardLabel{}) {
		return Rule{}
	}
	if !board.UserIsActiveMember(user) {
		return Rule{}
	}

	action := board.BuildAction(user.Username)
	trigger := board.BuildTrigger(label)
	rule := Rule{
		ID:         RuleID(newId()),
		Title:      fmt.Sprintf("Ajout %s (étiquette %s) ", user.Username, label.Name),
		TriggerID:  &trigger.ID,
		ActionID:   &action.ID,
		BoardID:    board.ID,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Action:     action,
		Trigger:    trigger,
	}
	return rule
}

func (wekan *Wekan) InsertRule(ctx context.Context, rule Rule) error {
	if err := wekan.AssertHasAdmin(ctx); err != nil {
		return err
	}

	if rule == (Rule{}) {
		return InsertEmptyRuleError{}
	}
	err := wekan.InsertAction(ctx, rule.Action)
	if err != nil {
		return err
	}

	err = wekan.InsertTrigger(ctx, rule.Trigger)
	if err != nil {
		return err
	}

	_, err = wekan.db.Collection("rules").InsertOne(ctx, rule)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) InsertAction(ctx context.Context, action Action) error {
	if err := wekan.AssertHasAdmin(ctx); err != nil {
		return err
	}

	_, err := wekan.db.Collection("actions").InsertOne(ctx, action)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) InsertTrigger(ctx context.Context, trigger Trigger) error {
	if err := wekan.AssertHasAdmin(ctx); err != nil {
		return err
	}

	_, err := wekan.db.Collection("triggers").InsertOne(ctx, trigger)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) SelectRulesFromBoardID(ctx context.Context, boardID BoardID) ([]Rule, error) {
	return nil, NotImplemented{}
}

func (wekan *Wekan) SelectRuleFromID(ctx context.Context, ruleID RuleID) (Rule, error) {
	return Rule{}, NotImplemented{}
}

func (wekan *Wekan) RemoveRuleWithID(ctx context.Context, ruleID RuleID) error {
	return NotImplemented{}
}
