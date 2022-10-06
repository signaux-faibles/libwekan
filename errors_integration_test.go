//go:build integration

// nolint:errcheck
package libwekan

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestErrors_UpstreamDeadlineExceeded(t *testing.T) {
	badWekan := newTestBadWekan("notAWekanDB")
	errs := []error{
		badWekan.AddMemberToBoard(ctx, "", BoardMember{}),
		badWekan.AddMemberToCard(ctx, "", ""),
		badWekan.AssertPrivileged(ctx),
		badWekan.CheckDocuments(ctx, UserID("")),
		badWekan.DisableBoardMember(ctx, "", ""),
		badWekan.DisableUser(ctx, User{}),
		badWekan.DisableUsers(ctx, Users{User{}}),
		badWekan.EnableBoardMember(ctx, "", ""),
		badWekan.EnableUser(ctx, User{}),
		badWekan.EnableUsers(ctx, Users{User{}}),
		badWekan.EnsureMemberInCard(ctx, "", ""),
		badWekan.EnsureMemberOutOfCard(ctx, "", ""),
		badWekan.EnsureRuleExists(ctx, User{}, Board{}, BoardLabel{}),
		badWekan.EnsureUserIsActiveBoardMember(ctx, "", ""),
		badWekan.EnsureUserIsInactiveBoardMember(ctx, "", ""),
		badWekan.EnsureUserIsInactiveBoardMember(ctx, "", ""),
		badWekan.InsertAction(ctx, Action{}),
		badWekan.InsertBoard(ctx, Board{}),
		badWekan.InsertBoardLabel(ctx, Board{}, BoardLabel{}),
		badWekan.InsertCard(ctx, Card{}),
		badWekan.InsertList(ctx, List{}),
		badWekan.InsertRule(ctx, Rule{}),
		badWekan.InsertSwimlane(ctx, Swimlane{}),
		badWekan.InsertTemplates(ctx, UserTemplates{}),
		badWekan.InsertTrigger(ctx, Trigger{}),
		badWekan.InsertUser(ctx, User{}),
		badWekan.InsertUsers(ctx, Users{User{}}),
		badWekan.RemoveMemberFromCard(ctx, "", ""),
		badWekan.RemoveRuleWithID(ctx, ""),
		ActivityID("").Check(ctx, &badWekan),
		BoardID("").Check(ctx, &badWekan),
		CardID("").Check(ctx, &badWekan),
		ListID("").Check(ctx, &badWekan),
		RuleID("").Check(ctx, &badWekan),
		SwimlaneID("").Check(ctx, &badWekan),
	}
	var err error
	_, err = ActivityID("").GetDocument(ctx, &badWekan)
	errs = append(errs, err)
	_, err = BoardID("").GetDocument(ctx, &badWekan)
	errs = append(errs, err)
	_, err = CardID("").GetDocument(ctx, &badWekan)
	errs = append(errs, err)
	_, err = ListID("").GetDocument(ctx, &badWekan)
	errs = append(errs, err)
	_, err = RuleID("").GetDocument(ctx, &badWekan)
	errs = append(errs, err)
	_, err = SwimlaneID("").GetDocument(ctx, &badWekan)
	errs = append(errs, err)
	_, err = badWekan.GetActivityFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetBoardFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetBoardFromSlug(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetBoardFromTitle(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetCardFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetListFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetSwimlaneFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetSwimlanesFromBoardID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetUserFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetUserFromUsername(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.GetUsers(ctx)
	errs = append(errs, err)
	_, err = badWekan.GetUsersFromIDs(ctx, []UserID{""})
	errs = append(errs, err)
	_, err = badWekan.GetUsersFromUsernames(ctx, []Username{""})
	errs = append(errs, err)
	_, err = badWekan.SelectActivitiesFromBoardID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectActivitiesFromQuery(ctx, bson.M{})
	errs = append(errs, err)
	_, err = badWekan.SelectActivityFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectBoardsFromMemberID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectCardsFromBoardID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectCardsFromListID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectCardsFromQuery(ctx, bson.M{})
	errs = append(errs, err)
	_, err = badWekan.SelectCardsFromUserID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectDomainBoards(ctx)
	errs = append(errs, err)
	_, err = badWekan.SelectListsFromBoardID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectRuleFromID(ctx, "")
	errs = append(errs, err)
	_, err = badWekan.SelectRulesFromBoardID(ctx, "")
	for i, err := range errs {
		assert.ErrorAs(t, err, &UnexpectedMongoError{}, "l'étape %d a échoué", i)
		assert.ErrorIs(t, err, context.DeadlineExceeded, "l'étape %d a échoué", i)
	}
}
