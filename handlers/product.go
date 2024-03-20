package handlers

import (
	"context"
	"errors"
	"net/http"
	"online_shopping_mongodb/db"
	"online_shopping_mongodb/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productDB = db.ConnectMongoDB()

// CreateProduct - creates product
func CreateProduct(c *gin.Context) {
	collection := productDB.Collection("products")
	var product models.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existingProduct models.Product
	err = collection.FindOne(context.TODO(), bson.M{"name": product.Name}).Decode(&existingProduct)
	if err == nil {
		existingProduct.Quantity += product.Quantity
		_, err = collection.ReplaceOne(context.TODO(), bson.M{"name": product.Name}, existingProduct)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, existingProduct)
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		genUUID := uuid.NewString()
		product.UUID = genUUID
		_, err = collection.InsertOne(context.TODO(), product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, product)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
}

// UpdateProduct - updates product
func UpdateProduct(c *gin.Context) {
	collection := productDB.Collection("products")
	var product models.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	filter := bson.M{"uuid": product.UUID}

	update := bson.M{
		"$set": bson.M{
			"name":     product.Name,
			"price":    product.Price,
			"quantity": product.Quantity,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetProductByUUID - gets product by uuid
func GetProductByUUID(c *gin.Context) {
	collection := productDB.Collection("products")
	reqUUID := c.Param("uuid")
	var product models.Product
	filter := bson.M{"uuid": reqUUID}
	err := collection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProductByUUID - deletes product by uuid
func DeleteProductByUUID(c *gin.Context) {
	reqUUID := c.Param("uuid")
	var product models.Product
	collection := productDB.Collection("products")
	filter := bson.M{"uuid": reqUUID}
	err := collection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetAllProducts - gets all products with pagination
func GetAllProducts(c *gin.Context) {
	collection := productDB.Collection("products")
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
	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
