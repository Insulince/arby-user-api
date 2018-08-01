package models

import "gopkg.in/mgo.v2/bson"

// TODO: Implement username support.
type User struct {
	Id                bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Email             string        `json:"email" bson:"email"`
	PasswordHash      []byte        `json:"password-hash" bson:"password-hash"`
	CreationTimestamp int64         `json:"creation-timestamp" bson:"creation-timestamp"`
}
