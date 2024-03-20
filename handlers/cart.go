package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"online_shopping_mongodb/db"
	"online_shopping_mongodb/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cartDB = db.ConnectMongoDB()

// AddOrUpdateCart - adds a user's cart with product ids. If user_uuid exists, updates only quantity. If user_uuid not exists, creates new user
func AddOrUpdateCart(c *gin.Context) {
	collection := cartDB.Collection("carts")
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existingCart models.Cart
	if err := collection.FindOne(context.Background(), bson.M{"user_uuid": cart.UserUUID}).Decode(&existingCart); err == nil {
		for _, newProduct := range cart.Products {
			found := false
			for i, existingProduct := range existingCart.Products {
				if existingProduct.UUID == newProduct.UUID {
					existingCart.Products[i].Quantity += newProduct.Quantity
					found = true
					break
				}
			}
			if !found {
				existingCart.Products = append(existingCart.Products, newProduct)
			}
		}

		if _, err := collection.ReplaceOne(context.Background(), bson.M{"user_uuid": cart.UserUUID}, existingCart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update cart",
			})
			return
		}
		var sum float64

		for _, existproduct := range existingCart.Products {
			var product models.Product
			_ = cartDB.Collection("products").FindOne(context.TODO(), bson.M{"uuid": existproduct.UUID}).Decode(&product)
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{
			// 		"error": err.Error(),
			// 	})
			// 	return
			// }
			sum += (float64(existproduct.Quantity) * product.Price)
		}

		existingCart.TotalAmount = sum
		// pp.Println(existingCart)
		c.JSON(http.StatusOK, existingCart)

	} else if errors.Is(err, mongo.ErrNoDocuments) {
		if _, err := collection.InsertOne(context.Background(), cart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}

		var sum float64

		for _, existproduct := range cart.Products {
			var product models.Product
			err := cartDB.Collection("products").FindOne(context.TODO(), bson.M{"uuid": existproduct.UUID}).Decode(&product)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			sum += (float64(existproduct.Quantity) * product.Price)
		}
		cart.TotalAmount = sum
		c.JSON(http.StatusOK, cart)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check cart"})
		return
	}
}

// DeleteFromCart deletes a given product ID and quantity from the user's cart
func DeleteFromCart(c *gin.Context) {
	collection := cartDB.Collection("carts")
	var deleteReq models.Cart
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existingCart models.Cart
	if err := collection.FindOne(context.Background(), bson.M{"user_uuid": deleteReq.UserUUID}).Decode(&existingCart); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	for _, deleteProduct := range deleteReq.Products {
		for i, existingProduct := range existingCart.Products {
			if existingProduct.UUID == deleteProduct.UUID {
				if existingProduct.Quantity > deleteProduct.Quantity {
					existingCart.Products[i].Quantity -= deleteProduct.Quantity
				} else {
					existingCart.Products = append(existingCart.Products[:i], existingCart.Products[i+1:]...)
				}
				break
			}
		}
	}

	if _, err := collection.ReplaceOne(context.Background(), bson.M{"user_uuid": deleteReq.UserUUID}, existingCart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}
	var sum float64
	for _, existproduct := range existingCart.Products {
		var product models.Product
		err := cartDB.Collection("products").FindOne(context.TODO(), bson.M{"uuid": existproduct.UUID}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		sum += (float64(existproduct.Quantity) * product.Price)
	}
	existingCart.TotalAmount = sum

	c.JSON(http.StatusOK, existingCart)
}

func GetAllCarts(c *gin.Context) {
	stringpage := c.Param("page")
	stringlimit := c.Param("limit")
	page, err := strconv.Atoi(stringpage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	limit, err := strconv.Atoi(stringlimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	offset := (page - 1) * limit
	collection := cartDB.Collection("carts")
	options := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit))
	cursor, err := collection.Find(context.TODO(), bson.M{}, options)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var carts []models.Cart
	for cursor.Next(context.TODO()) {
		var cart models.Cart
		if err := cursor.Decode(&cart); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		carts = append(carts, cart)
	}

	c.JSON(http.StatusOK, carts)
}

func ClearCartByUserUUID(c *gin.Context) {
	userUUID := c.Param("uuid")
	fmt.Println(userUUID)
	collection := cartDB.Collection("carts")
	count, err := collection.CountDocuments(context.TODO(), bson.M{"user_uuid": userUUID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User UUID does not exist"})
		return
	}

	_, err = collection.DeleteMany(context.TODO(), bson.M{"user_uuid": userUUID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.
		H{"message": "Cart cleared successfully"})
}

func GetAllCartsWithoutUUID(c *gin.Context) {
	stringPage := c.Param("page")
	stringLimit := c.Param("limit")
	page, err := strconv.Atoi(stringPage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page number",
		})
		return
	}

	limit, err := strconv.Atoi(stringLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit",
		})
		return
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10 // default limit
	}

	offset := (page - 1) * limit

	collection := cartDB.Collection("carts")
	options := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit))

	cursor, err := collection.Find(context.Background(), bson.M{}, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch carts",
		})
		return
	}
	defer cursor.Close(context.Background())

	var carts []models.Cart
	for cursor.Next(context.Background()) {
		var cart models.Cart
		if err := cursor.Decode(&cart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to decode cart",
			})
			return
		}
		carts = append(carts, cart)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cursor error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"carts": carts,
	})
}
