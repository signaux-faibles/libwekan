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

func TestActivities_newActivityJoinMember(t *testing.T) {
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
	activity := newActivityJoinMember(
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

func TestActivitie_newActivityAddedLabel(t *testing.T) {
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
