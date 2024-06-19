package utils

import (
	"GO-07_mongoDB_RMS/database"
	"GO-07_mongoDB_RMS/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"math"
	"os"
	"time"
)

var tokenBlacklist = make(map[string]bool)

func GenerateToken(info models.UserInfo) (models.LoginResponse, error) {
	res := models.LoginResponse{}

	token, err := GenerateTokenPair(info.Id)
	if err != nil {
		return res, err
	}

	res.IsValid = true
	res.Token = token["token"]
	res.RefreshToken = token["refresh_token"]
	return res, nil
}

func GenerateTokenPair(userID primitive.ObjectID) (map[string]string, error) {

	jwtExpirationTime := time.Now().Add(time.Minute * 60).Unix()

	var userInfo models.UserInfo
	err := database.UserCollection.FindOne(database.MongoCtx, bson.M{"_id": userID}).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	if userInfo.Email == "" {
		return nil, errors.New("email not found for admin")
	}

	claims := &models.UserClaims{
		UserID: userID,
		Roles:  userInfo.Roles,
		Email:  userInfo.Email,
		Name:   userInfo.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtExpirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("jwtSecret")))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["userId"] = userID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	rt, err := refreshToken.SignedString([]byte(os.Getenv("jwtSecret")))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"token":         t,
		"refresh_token": rt,
	}, nil
}

func InvalidateJWT(tokenString string) error {
	// Check if the token is valid
	token, err := jwt.ParseWithClaims(tokenString, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		// Return the secret key for validation
		return []byte(os.Getenv("jwtSecret")), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	// Check if the token already exists in the blacklist
	if _, found := tokenBlacklist[tokenString]; found {
		return errors.New("token already invalidated")
	}

	// Add the token to the blacklist
	tokenBlacklist[tokenString] = true
	return nil
}

func IsTokenValid(tokenString string) bool {
	if _, found := tokenBlacklist[tokenString]; found {
		return false
	}
	return true
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserHighestRole(uc models.UserClaims) models.RoleType {

	for _, role := range uc.Roles {
		if role == models.AdminRole {
			return models.AdminRole
		}
	}

	return models.SubAdminRole
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of the Earth in kilometers
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
