package mongo

import (
	"gopkg.in/mgo.v2"
	"arby-user-api/pkg/configuration"
)

const UsersCollectionName = "users"

var db *mgo.Database

func InitializeDatabase(config *configuration.Config) (err error) {
	session, err := mgo.Dial(config.Mongo.ConnectionString)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Strong, true)
	db = session.DB(config.Mongo.DatabaseName)

	return nil
}

func Users() (users *mgo.Collection) {
	return db.C(UsersCollectionName)
}
