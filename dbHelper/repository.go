package dbHelper

import (
	"GO-07_mongoDB_RMS/database"
	"GO-07_mongoDB_RMS/models"
	"GO-07_mongoDB_RMS/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func VerifyLoginCredentials(credentials models.LoginCredentials) (models.UserInfo, error) {
	var userInfo models.UserInfo
	err := database.UserCollection.FindOne(database.MongoCtx, bson.M{"email": credentials.Email}).Decode(&userInfo)
	if err != nil {
		return userInfo, err
	}

	checkErr := utils.CheckPassword(credentials.Password, userInfo.Pwd)
	return userInfo, checkErr
}
