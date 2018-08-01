package mongo

import (
	"gopkg.in/mgo.v2"
	"arby-user-api/pkg/configuration"
	"arby-user-api/pkg/models"
	"gopkg.in/mgo.v2/bson"
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

func InsertUser(user models.User) (err error) {
	return Users().Insert(user)
}

func FindUserByEmail(email string) (user *models.User, err error) {
	err = Users().Find(bson.M{"email": email}).One(&user)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	return user, nil
}
