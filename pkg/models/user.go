package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id                bson.ObjectId   `json:"_id" bson:"_id,omitempty"`
	Username          string          `json:"username" bson:"username"`
	PasswordHash      []byte          `json:"password-hash" bson:"password-hash"`
	CreationTimestamp int64           `json:"creation-timestamp" bson:"creation-timestamp"`
}
