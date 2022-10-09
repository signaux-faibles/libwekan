package libwekan

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Rules []Rule

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

func (ruleID RuleID) GetDocument(ctx context.Context, wekan *Wekan) (Rule, error) {
	return wekan.SelectRuleFromID(ctx, ruleID)
}

func (ruleID RuleID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.SelectRuleFromID(ctx, ruleID)
	return err
}

func (board Board) BuildTriggerAddedLabel(label BoardLabel) Trigger {
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

func (board Board) BuildTriggerRemovedLabel(label BoardLabel) Trigger {
	return Trigger{
		ID:           TriggerID(newId()),
		ActivityType: "removedLabel",
		BoardID:      board.ID,
		LabelID:      label.ID,
		Description:  fmt.Sprintf("quand l'étiquette %s est retirée de la carte par *", label.Name),
		UserID:       "*",
		CreatedAt:    time.Now(),
		ModifiedAt:   time.Now(),
	}
}

func (board Board) BuildActionAddMember(username Username) Action {
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

func (board Board) BuildActionRemoveMember(username Username) Action {
	return Action{
		ID:          ActionID(newId()),
		ActionType:  "removeMember",
		Username:    username,
		BoardID:     board.ID,
		Description: fmt.Sprintf("%s est exclu de la carte", username),
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}
}

func (board Board) BuildRuleAddMember(user User, labelName BoardLabelName) Rule {
	label := board.GetLabelByName(labelName)
	if label == (BoardLabel{}) {
		return Rule{}
	}
	if !board.UserIsActiveMember(user) {
		return Rule{}
	}

	action := board.BuildActionAddMember(user.Username)
	trigger := board.BuildTriggerAddedLabel(label)
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

func (board Board) BuildRuleRemoveMember(user User, labelName BoardLabelName) Rule {
	label := board.GetLabelByName(labelName)
	if label == (BoardLabel{}) {
		return Rule{}
	}
	if !board.UserIsActiveMember(user) {
		return Rule{}
	}

	action := board.BuildActionRemoveMember(user.Username)
	trigger := board.BuildTriggerRemovedLabel(label)
	rule := Rule{
		ID:         RuleID(newId()),
		Title:      fmt.Sprintf("Suppression %s (étiquette %s) ", user.Username, label.Name),
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

func (rules Rules) selectUser(username Username) Rules {
	return selectSlice(rules, func(rule Rule) bool { return rule.Action.Username == username })
}

func (rules Rules) SelectBoardLabelName(boardLabelID BoardLabelID) Rules {
	return selectSlice(rules, func(rule Rule) bool { return rule.Trigger.LabelID == boardLabelID })
}

func (rules Rules) SelectRemoveMemberFromTaskforceRule() Rules {
	return selectSlice(rules, func(rule Rule) bool {
		return rule.Action.ActionType == "removeMember" && rule.Trigger.ActivityType == "removedLabel"
	})
}

func (rules Rules) SelectAddMemberToTaskforceRule() Rules {
	return selectSlice(rules, func(rule Rule) bool {
		return rule.Action.ActionType == "addMember" && rule.Trigger.ActivityType == "addedLabel"
	})
}

func (wekan *Wekan) InsertRule(ctx context.Context, rule Rule) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
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
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	_, err := wekan.db.Collection("actions").InsertOne(ctx, action)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) InsertTrigger(ctx context.Context, trigger Trigger) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	_, err := wekan.db.Collection("triggers").InsertOne(ctx, trigger)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) SelectRulesFromBoardID(ctx context.Context, boardID BoardID) (Rules, error) {
	var malformedRules, rules Rules
	cur, err := wekan.db.Collection("rules").Find(ctx, bson.M{"boardId": boardID})
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	if err := cur.All(ctx, &malformedRules); err != nil {
		return nil, UnexpectedMongoError{err}
	}
	for _, malformedRule := range malformedRules {
		rule, err := wekan.SelectRuleFromID(ctx, malformedRule.ID)
		if err != nil {
			return Rules{}, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (wekan *Wekan) SelectRuleFromID(ctx context.Context, ruleID RuleID) (Rule, error) {
	var rule Rule
	err := wekan.db.Collection("rules").FindOne(ctx, bson.M{"_id": ruleID}).Decode(&rule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Rule{}, RuleNotFoundError{ruleID}
		}
		return Rule{}, UnexpectedMongoError{err}
	}
	if err := wekan.db.Collection("actions").FindOne(ctx, bson.M{"_id": rule.ActionID}).Decode(&rule.Action); err != nil {
		if err != mongo.ErrNoDocuments {
			return Rule{}, UnexpectedMongoError{err}
		}
		if rule.ActionID != nil {
			return Rule{}, ActionNotFoundError{*rule.ActionID}
		}
	}
	if err := wekan.db.Collection("triggers").FindOne(ctx, bson.M{"_id": rule.TriggerID}).Decode(&rule.Trigger); err != nil {
		if err != mongo.ErrNoDocuments {
			return Rule{}, UnexpectedMongoError{err}
		}
		if rule.TriggerID != nil {
			return Rule{}, TriggerNotFoundError{*rule.TriggerID}
		}
	}
	return rule, nil
}

func (wekan *Wekan) RemoveRuleWithID(ctx context.Context, ruleID RuleID) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	rule, err := ruleID.GetDocument(ctx, wekan)
	if err != nil {
		return err
	}
	_, err = wekan.db.Collection("actions").DeleteOne(ctx, bson.M{"_id": rule.Action.ID})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	_, err = wekan.db.Collection("triggers").DeleteOne(ctx, bson.M{"_id": rule.Trigger.ID})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	_, err = wekan.db.Collection("rules").DeleteOne(ctx, bson.M{"_id": rule.ID})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) EnsureRuleAddTaskforceMemberExists(ctx context.Context, user User, board Board, boardLabel BoardLabel) (bool, error) {
	boardRules, err := wekan.SelectRulesFromBoardID(ctx, board.ID)
	if err != nil {
		return false, err
	}
	existingRules := boardRules.SelectBoardLabelName(boardLabel.ID).selectUser(user.Username).SelectAddMemberToTaskforceRule()
	if len(existingRules) == 0 {
		rule := board.BuildRuleAddMember(user, boardLabel.Name)
		err = wekan.InsertRule(ctx, rule)
		return err == nil, err
	}
	return false, nil
}

func (wekan *Wekan) EnsureRuleRemoveTaskforceMemberExists(ctx context.Context, user User, board Board, boardLabel BoardLabel) (bool, error) {
	boardRules, err := wekan.SelectRulesFromBoardID(ctx, board.ID)
	if err != nil {
		return false, err
	}
	existingRules := boardRules.SelectBoardLabelName(boardLabel.ID).selectUser(user.Username).SelectRemoveMemberFromTaskforceRule()
	if len(existingRules) == 0 {
		rule := board.BuildRuleRemoveMember(user, boardLabel.Name)
		err = wekan.InsertRule(ctx, rule)
		return err == nil, err
	}
	return false, nil
}
