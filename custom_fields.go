package libwekan

import "time"

type CustomField struct {
	ID                  CardCustomFieldID `bson:"_id" json:"_id,omitempty"`
	Name                string
	Type                string
	Settings            CustomFieldSettings
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
		ID   string `bson:"_id" json:"_id,omitempty"`
		Name string `bson:"name" json:"name,omitempty"`
	} `bson:"dropdownItems" json:"dropdownItems,omitempty"`
}
