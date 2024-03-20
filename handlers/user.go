package handlers

import (
	"context"
	"net/http"
	"online_shopping_mongodb/db"
	"online_shopping_mongodb/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userDB = db.ConnectMongoDB()

// CreateUser - creates user
func CreateUser(c *gin.Context) {
	collection := userDB.Collection("users")
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	genUUID := uuid.NewString()
	user.UUID = genUUID
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	filter := bson.M{"uuid": user.UUID}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser - updates user
func UpdateUser(c *gin.Context) {
	collection := userDB.Collection("users")
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	filter := bson.M{"uuid": user.UUID}

	update := bson.M{
		"$set": bson.M{
			"name": user.Name,
			"age":  user.Age,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByUUID - gets user by uuid
func GetUserByUUID(c *gin.Context) {
	collection := userDB.Collection("users")
	reqUUID := c.Param("uuid")
	var user models.User
	filter := bson.M{"uuid": reqUUID}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUserByUUID - deletes user by uuid
func DeleteUserByUUID(c *gin.Context) {
	reqUUID := c.Param("uuid")
	var user models.User
	collection := userDB.Collection("users")
	filter := bson.M{"uuid": reqUUID}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetAllUsers - gets all users with pagination
func GetAllUsers(c *gin.Context) {
	collection := userDB.Collection("users")

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.Param("pageSize"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	skip := (page - 1) * pageSize

	options := options.Find().SetLimit(int64(pageSize)).SetSkip(int64(skip))
	cursor, err := collection.Find(context.TODO(), bson.M{}, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
