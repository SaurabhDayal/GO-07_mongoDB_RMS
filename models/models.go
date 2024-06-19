package models

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleType string

const (
	AdminRole    RoleType = "admin"
	SubAdminRole RoleType = "sub-admin"
	CustomerRole RoleType = "customer"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password,omitempty" bson:"password"`
	Credit    int                `json:"credit" bson:"credit"`
	Roles     []RoleType         `json:"roles" bson:"roles"`
	Addresses []Address          `json:"addresses,omitempty" bson:"addresses,omitempty"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name" bson:"name,omitempty"`
	Email string             `json:"email" bson:"email"`
	Pwd   string             `json:"-" bson:"password"`
	Roles []RoleType         `json:"roles"`
}

type UserClaims struct {
	UserID primitive.ObjectID `json:"userId"`
	Name   string             `json:"name"`
	Email  string             `json:"email"`
	Roles  []RoleType         `json:"roles"`
	jwt.StandardClaims
}

type LoginResponse struct {
	IsValid      bool   `json:"isValid"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type Address struct {
	Name      string  `json:"name" bson:"name"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type AddressDistance struct {
	UserAddId int     `json:"userAddId"`
	RestId    int     `json:"restId"`
	Distance  float64 `json:"distance"`
}

type Restaurant struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	Address       *Address           `json:"address" bson:"address"`
	OwnedByUserID primitive.ObjectID `json:"ownedByUserId" bson:"ownedByUserId"`
}

type Dish struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Cost            int                `json:"cost" bson:"cost"`
	PreparingTime   string             `json:"preparingTime" bson:"preparingTime"`
	RestaurantId    primitive.ObjectID `json:"restaurantId" bson:"restaurantId"`
	CreatedByUserID primitive.ObjectID `json:"createdByUserId" bson:"createdByUserId"`
}

type Order struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeliveryTime string             `json:"deliveryTime" bson:"deliveryTime"`
	IsCancelled  bool               `json:"isCancelled" bson:"isCancelled"`
	IsDelivered  bool               `json:"isDelivered" bson:"isDelivered"`
	DishId       primitive.ObjectID `json:"dishId" bson:"dishId"`
	CustomerID   primitive.ObjectID `json:"customerId" bson:"customerId"`
}

type Distance struct {
	DistanceInKM float64 `json:"distanceInKM"`
}
