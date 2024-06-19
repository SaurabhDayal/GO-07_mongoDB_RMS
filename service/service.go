package service

import (
	"GO-07_mongoDB_RMS/database"
	"GO-07_mongoDB_RMS/dbHelper"
	"GO-07_mongoDB_RMS/models"
	"GO-07_mongoDB_RMS/utils"
	"context"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterCustomer(user models.User) (models.User, error) {

	// Check if user email already exists
	var customerRoleExists bool
	existingUser := models.User{}
	err := database.UserCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// User with this email already exists, append CustomerRole to existing roles
		for _, role := range existingUser.Roles {
			if role == models.CustomerRole {
				customerRoleExists = true
			}
		}
		if !customerRoleExists {
			existingUser.Roles = append(existingUser.Roles, models.CustomerRole)
			existingUser.Credit = 1000
		}
		// Update existing user with CustomerRole
		update := bson.M{
			"$addToSet": bson.M{
				"roles":     models.CustomerRole,
				"addresses": bson.M{"$each": user.Addresses},
			},
			"$set": bson.M{
				"credit": existingUser.Credit,
			},
		}
		_, err := database.UserCollection.UpdateOne(context.Background(), bson.M{"_id": existingUser.ID}, update)
		if err != nil {
			return existingUser, err
		}

		err = database.UserCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
		existingUser.Password = ""
		return existingUser, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		// User with this email does not exist, proceed with adding new user
		user.Roles = []models.RoleType{models.CustomerRole}

		password, err := utils.HashPassword(user.Password)
		if err != nil {
			return user, errors.Wrap(err, "setPassword")
		}
		user.Password = password
		user.Credit = 1000

		result, err := database.UserCollection.InsertOne(context.Background(), user)
		if err != nil {
			return user, err
		}
		user.ID = result.InsertedID.(primitive.ObjectID)
		user.Password = ""

		return user, nil
	} else {
		// Handle other errors
		return user, err
	}

}

func LoginUser(credentials models.LoginCredentials) (models.LoginResponse, error) {

	loginResponse := models.LoginResponse{}
	info, err := dbHelper.VerifyLoginCredentials(credentials)
	if err != nil {
		return loginResponse, err
	}

	return utils.GenerateToken(info)
}

func LogoutUser(tokenString string) error {

	err := utils.InvalidateJWT(tokenString)
	if err != nil {
		return err
	}
	return nil
}

func CreateSubAdmin(user models.User) (models.User, error) {

	// Check if user email already exists
	var subAdminRoleExists bool
	existingUser := models.User{}
	err := database.UserCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		for _, role := range existingUser.Roles {
			if role == models.CustomerRole {
				subAdminRoleExists = true
			}
		}
		if !subAdminRoleExists {
			existingUser.Roles = append(existingUser.Roles, models.SubAdminRole)
		}
		update := bson.M{
			"$addToSet": bson.M{
				"roles": models.SubAdminRole,
			},
		}
		_, err := database.UserCollection.UpdateOne(context.Background(), bson.M{"_id": existingUser.ID}, update)
		if err != nil {
			return existingUser, err
		}

		err = database.UserCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
		existingUser.Password = ""
		return existingUser, nil

	} else if errors.Is(err, mongo.ErrNoDocuments) {

		user.Roles = []models.RoleType{models.SubAdminRole}
		user.Addresses = nil

		password, err := utils.HashPassword(user.Password)
		if err != nil {
			return user, errors.Wrap(err, "setPassword")
		}
		user.Password = password

		result, err := database.UserCollection.InsertOne(database.MongoCtx, user)
		if err != nil {
			return user, err
		}
		user.ID = result.InsertedID.(primitive.ObjectID)
		user.Password = ""

		return user, nil
	} else {
		return user, err
	}

}

func GetSubAdminList(limit, offset int) ([]models.User, error) {

	var subAdmins []models.User

	filter := bson.M{
		"roles": models.SubAdminRole,
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := database.UserCollection.Find(database.MongoCtx, filter, findOptions)
	if err != nil {
		return subAdmins, err
	}
	defer cursor.Close(database.MongoCtx)

	for cursor.Next(database.MongoCtx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return subAdmins, err
		}
		user.Password = ""
		subAdmins = append(subAdmins, user)
	}

	if err := cursor.Err(); err != nil {
		return subAdmins, err
	}

	return subAdmins, nil
}

func CreateRestaurant(restaurant models.Restaurant) (models.Restaurant, error) {

	result, err := database.RestaurantCollection.InsertOne(database.MongoCtx, restaurant)
	if err != nil {
		return restaurant, err
	}
	restaurant.ID = result.InsertedID.(primitive.ObjectID)
	return restaurant, nil
}

func CreateDish(dish models.Dish) (models.Dish, error) {

	result, err := database.DishCollection.InsertOne(database.MongoCtx, dish)
	if err != nil {
		return dish, err
	}
	dish.ID = result.InsertedID.(primitive.ObjectID)
	return dish, nil
}

func GetRestaurantByOwnerId(limit, offset int, uc models.UserClaims) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	userHighestRole := utils.GetUserHighestRole(uc)

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	var filter bson.M
	if userHighestRole == models.SubAdminRole {
		filter = bson.M{
			"ownedByUserId": uc.UserID,
		}
	}

	cursor, err := database.RestaurantCollection.Find(database.MongoCtx, filter, findOptions)
	if err != nil {
		return restaurants, err
	}
	defer cursor.Close(database.MongoCtx)

	for cursor.Next(database.MongoCtx) {
		var restaurant models.Restaurant
		if err := cursor.Decode(&restaurant); err != nil {
			return restaurants, err
		}
		restaurants = append(restaurants, restaurant)
	}

	if err := cursor.Err(); err != nil {
		return restaurants, err
	}

	return restaurants, nil
}

func GetDishByCreatedUserId(limit, offset int, uc models.UserClaims) ([]models.Dish, error) {
	var dishes []models.Dish
	userHighestRole := utils.GetUserHighestRole(uc)

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	var filter bson.M
	if userHighestRole == models.SubAdminRole {
		filter = bson.M{
			"createdByUserId": uc.UserID,
		}
	}

	cursor, err := database.DishCollection.Find(database.MongoCtx, filter, findOptions)
	if err != nil {
		return dishes, err
	}
	defer cursor.Close(database.MongoCtx)

	for cursor.Next(database.MongoCtx) {
		var dish models.Dish
		if err := cursor.Decode(&dish); err != nil {
			return dishes, err
		}
		dishes = append(dishes, dish)
	}

	if err := cursor.Err(); err != nil {
		return dishes, err
	}

	return dishes, nil
}

func GetUserList(limit, offset int) ([]models.User, error) {

	var subAdmins []models.User

	filter := bson.M{
		"roles": models.CustomerRole,
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := database.UserCollection.Find(database.MongoCtx, filter, findOptions)
	if err != nil {
		return subAdmins, err
	}
	defer cursor.Close(database.MongoCtx)

	for cursor.Next(database.MongoCtx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return subAdmins, err
		}
		user.Password = ""
		subAdmins = append(subAdmins, user)
	}

	if err := cursor.Err(); err != nil {
		return subAdmins, err
	}

	return subAdmins, nil
}

func GetRestaurantList(limit, offset int) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	var filter bson.M

	cursor, err := database.RestaurantCollection.Find(database.MongoCtx, filter, findOptions)
	if err != nil {
		return restaurants, err
	}
	defer cursor.Close(database.MongoCtx)

	for cursor.Next(database.MongoCtx) {
		var restaurant models.Restaurant
		if err := cursor.Decode(&restaurant); err != nil {
			return restaurants, err
		}
		restaurants = append(restaurants, restaurant)
	}

	if err := cursor.Err(); err != nil {
		return restaurants, err
	}

	return restaurants, nil
}

func GetRestaurantDishList(limit, offset int, resID primitive.ObjectID) ([]models.Dish, error) {
	var dishes []models.Dish

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	var filter bson.M
	filter = bson.M{
		"restaurantId": resID,
	}

	cursor, err := database.DishCollection.Find(database.MongoCtx, filter, findOptions)
	if err != nil {
		return dishes, err
	}
	defer cursor.Close(database.MongoCtx)

	for cursor.Next(database.MongoCtx) {
		var dish models.Dish
		if err := cursor.Decode(&dish); err != nil {
			return dishes, err
		}
		dishes = append(dishes, dish)
	}

	if err := cursor.Err(); err != nil {
		return dishes, err
	}

	return dishes, nil
}

func GetDistance(latitude, longitude float64, resID primitive.ObjectID) (models.Distance, error) {
	var distance float64
	var restaurant models.Restaurant
	var distanceBetween models.Distance

	filter := bson.M{
		"_id": resID,
	}

	err := database.RestaurantCollection.FindOne(context.Background(), filter).Decode(&restaurant)
	if err != nil {
		return distanceBetween, err
	}

	if restaurant.Address == nil {
		return distanceBetween, errors.New("restaurant location not found")
	}

	distance = utils.Haversine(latitude, longitude, restaurant.Address.Latitude, restaurant.Address.Longitude)
	distanceBetween.DistanceInKM = distance

	return distanceBetween, nil
}
