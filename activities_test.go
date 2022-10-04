package libwekan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivities_NewActivityCreateBoard(t *testing.T) {
	ass := assert.New(t)
	expected := Activity{
		BoardID:        "boardID",
		ActivityTypeID: "boardID",
		UserID:         "userID",
		ActivityType:   "createBoard",
		Type:           "board",
	}
	activity := newActivityCreateBoard("userID", "boardID")
	ass.Equal(expected, activity)
}

func TestActivities_newActivityCreateSwimlane(t *testing.T) {
	ass := assert.New(t)
	expected := Activity{
		SwimlaneID:   "swimlaneID",
		UserID:       "userID",
		BoardID:      "boardID",
		ActivityType: "createSwimlane",
		Type:         "list",
	}
	activity := newActivityCreateSwimlane("userID", "boardID", "swimlaneID")
	ass.Equal(expected, activity)
}

func TestActivities_newActivityAddBoardMember(t *testing.T) {
	ass := assert.New(t)
	expected := Activity{
		UserID:       "userID",
		MemberID:     "memberID",
		BoardID:      "boardID",
		ActivityType: "addBoardMember",
		Type:         "member",
	}
	activity := newActivityAddBoardMember("userID", "memberID", "boardID")
	ass.Equal(expected, activity)
}

func TestActivities_newActivityRemoveBoardMember(t *testing.T) {
	ass := assert.New(t)
	expected := Activity{
		UserID:       "userID",
		MemberID:     "memberID",
		BoardID:      "boardID",
		ActivityType: "removeBoardMember",
		Type:         "member",
	}
	activity := newActivityRemoveBoardMember("userID", "memberID", "boardID")
	ass.Equal(expected, activity)
}

func TestActivities_newActivityCardJoinMember(t *testing.T) {
	ass := assert.New(t)
	expected := Activity{
		UserID:       "userID",
		Username:     "username",
		MemberID:     "memberID",
		BoardID:      "boardID",
		CardID:       "cardID",
		ListID:       "listID",
		SwimlaneID:   "swimlaneID",
		ActivityType: "joinMember",
	}
	activity := newActivityCardJoinMember(
		"userID", "username", "memberID", "boardID",
		"listID", "cardID", "swimlaneID",
	)
	ass.Equal(expected, activity)
}

func TestActivities_newActivityAddComment(t *testing.T) {
	ass := assert.New(t)
	expected := Activity{
		UserID:       "userID",
		BoardID:      "boardID",
		CardID:       "cardID",
		CommentID:    "commentID",
		ListID:       "listID",
		SwimlaneID:   "swimlaneID",
		ActivityType: "addComment",
	}
	activity := newActivityAddComment(
		"userID", "boardID", "cardID",
		"commentID", "listID", "swimlaneID",
	)
	ass.Equal(expected, activity)
}

func TestActivities_newActivityAddedLabel(t *testing.T) {
	ass := assert.New(t)
	expected := Activity{
		UserID:       "userID",
		BoardLabelID: "boardLabelID",
		ActivityType: "addedLabel",
		BoardID:      "boardID",
		SwimlaneID:   "swimlaneID",
	}
	activity := newActivityAddedLabel("userID", "boardLabelID", "boardID", "swimlaneID")
	ass.Equal(expected, activity)
}

func TestActivities_newActivityCreateCard(t *testing.T) {
	expected := Activity{
		UserID:       "userID",
		ActivityType: "createCard",
		BoardID:      "card.BoardID",
		ListName:     "list.Title",
		ListID:       "list.ID",
		CardID:       "card.ID",
		CardTitle:    "card.Title",
		SwimlaneName: "swimlane.Title",
		SwimlaneID:   "swimlane.ID",
	}
	list := List{
		ID:    "list.ID",
		Title: "list.Title",
	}
	swimlane := Swimlane{
		ID:    "swimlane.ID",
		Title: "swimlane.Title",
	}
	card := Card{
		ID:      "card.ID",
		Title:   "card.Title",
		BoardID: "card.BoardID",
	}
	activity := newActivityCreateCard("userID", list, card, swimlane)
	assert.Equal(t, expected, activity)
}
