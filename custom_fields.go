package libwekan

import "time"

type CustomField struct {
	ID                  CustomFieldID `bson:"_id"`
	Name                string
	Type                string
	Settings            struct{}
	ShowOnCard          bool
	ShowLabelOnMiniCard bool
	AutomaticallyOnCard bool
	BoardIDs            []BoardID
	CreatedAt           time.Time
	ModifiedAt          time.Time
	ShowSumAtTopOfList  bool
}

type CustomFieldSettings struct {
	DropdownItems []struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	} `bson:"dropdownItems"`
}
