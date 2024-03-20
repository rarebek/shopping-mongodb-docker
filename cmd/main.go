package main

import (
	"log"
	"online_shopping_mongodb/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	api := router.Group("/v1")
	//User routes
	api.POST("/user/create", handlers.CreateUser)
	api.PUT("/user/update", handlers.UpdateUser)
	api.GET("/user/get/:uuid", handlers.GetUserByUUID)
	api.GET("/users/:page/:pageSize", handlers.GetAllUsers)
	api.DELETE("/user/delete/:uuid", handlers.DeleteUserByUUID)

	//Product routes
	api.POST("/product/create", handlers.CreateProduct)
	api.PUT("/product/update", handlers.UpdateProduct)
	api.GET("/product/get/:uuid", handlers.GetProductByUUID)
	api.GET("/products/:page/:pageSize", handlers.GetAllProducts)
	api.DELETE("/product/delete/:uuid", handlers.DeleteProductByUUID)

	//Shopping Cart routes
	api.POST("/cart/add", handlers.AddOrUpdateCart)
	api.DELETE("/cart/delete/", handlers.DeleteFromCart)
	api.GET("/cart/getall/:page/:limit", handlers.GetAllCarts)
	api.GET("/cart/getallcarts/:page/:limit", handlers.GetAllCartsWithoutUUID)
	api.DELETE("/cart/clear/:uuid", handlers.ClearCartByUserUUID)

	err := router.Run(":9999")
	if err != nil {
		log.Fatal(err)
	}
}
