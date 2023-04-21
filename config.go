package libwekan

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfigBoard struct {
	Board        Board                             `bson:"board"`
	Swimlanes    map[SwimlaneID]Swimlane           `bson:"swimlanes"`
	Lists        map[ListID]List                   `bson:"lists"`
	CustomFields map[CardCustomFieldID]CustomField `bson:"customFields"`
}

type Config struct {
	Boards map[BoardID]ConfigBoard `bson:"boards"`
	Users  map[UserID]User         `bson:"users"`
}

func (config *Config) Copy() Config {
	return *config
}

func (config *Config) GetUserByUsername(username Username) (User, bool) {
	for _, user := range config.Users {
		if user.Username == username {
			return user, true
		}
	}
	return User{}, false
}

func buildConfigPipeline(slugDomainRegexp string) []bson.M {
	matchBoards := bson.M{
		"$match": bson.M{
			"slug": bson.M{
				"$regex": primitive.Regex{
					Pattern: slugDomainRegexp,
					Options: "i",
				}}}}

	projectBoards := bson.M{
		"$project": bson.M{
			"board": "$$ROOT",
		},
	}

	lookupSwimlanes := bson.M{
		"$lookup": bson.M{
			"from": "swimlanes",
			"let":  bson.M{"boardId": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$eq": bson.A{"$boardId", "$$boardId"},
					}}},
				{"$project": bson.M{
					"_id": 0,
					"k":   "$_id",
					"v":   "$$ROOT",
				}},
			},
			"as": "swimlanes",
		},
	}

	lookupLists := bson.M{
		"$lookup": bson.M{
			"from": "lists",
			"let":  bson.M{"boardId": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"archived": false,
				}},
				{"$match": bson.M{
					"$expr": bson.M{
						"$eq": bson.A{"$boardId", "$$boardId"},
					}}},
				{"$project": bson.M{
					"_id": 0,
					"k":   "$_id",
					"v":   "$$ROOT",
				}},
			},
			"as": "lists",
		},
	}

	lookupCustomFields := bson.M{
		"$lookup": bson.M{
			"from": "customFields",
			"let":  bson.M{"boardId": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$in": bson.A{"$$boardId", "$boardIds"},
					}}},
				{"$project": bson.M{
					"_id": 0,
					"k":   "$_id",
					"v":   "$$ROOT",
				}},
			},
			"as": "customFields",
		},
	}

	buildBoardConfig := bson.M{
		"$addFields": bson.M{
			"swimlanes": bson.M{
				"$arrayToObject": "$swimlanes",
			},
			"lists": bson.M{
				"$arrayToObject": "$lists",
			},
			"customFields": bson.M{
				"$arrayToObject": "$customFields",
			},
		},
	}

	prepareBoardsFoArrayToObject := bson.M{
		"$project": bson.M{
			"_id": 0,
			"k":   "$_id",
			"v":   "$$ROOT",
		},
	}

	groupByBoard := bson.M{
		"$group": bson.M{
			"_id":    0,
			"boards": bson.M{"$push": "$$ROOT"},
		},
	}

	formatBoards := bson.M{
		"$project": bson.M{
			"_id":    0,
			"boards": bson.M{"$arrayToObject": "$boards"},
		}}

	lookupUsers := bson.M{
		"$lookup": bson.M{
			"from": "users",
			"let":  bson.M{},
			"pipeline": []bson.M{
				{"$match": bson.M{"username": bson.M{"$exists": true}}},
				{"$project": bson.M{
					"_id": 0,
					"k":   "$_id",
					"v":   "$$ROOT"}},
			},
			"as": "users",
		}}

	formatUsers := bson.M{
		"$project": bson.M{
			"_id":    0,
			"boards": 1,
			"users":  bson.M{"$arrayToObject": "$users"},
		}}

	return []bson.M{
		matchBoards,
		projectBoards,
		lookupSwimlanes,
		lookupLists,
		lookupCustomFields,
		buildBoardConfig,
		prepareBoardsFoArrayToObject,
		groupByBoard,
		formatBoards,
		lookupUsers,
		formatUsers,
	}
}

func (wekan *Wekan) SelectConfig(ctx context.Context) (Config, error) {
	var config Config
	if wekan.db == nil {
		return Config{}, fmt.Errorf("impossible de sélectionner la configuration car la base de données n'est pas définie")
	}
	pipeline := buildConfigPipeline(wekan.slugDomainRegexp)

	cur, err := wekan.db.Collection("boards").Aggregate(ctx, pipeline, nil)
	if err != nil {
		return config, UnexpectedMongoError{err}
	}
	cur.Next(ctx)
	err = cur.Decode(&config)
	if err != nil {
		fmt.Println(err)
		return Config{}, UnexpectedMongoDecodeError{err}
	}
	err = cur.Close(ctx)
	if err != nil {
		return Config{}, UnexpectedMongoDecodeError{err}
	}

	return config, nil
}
