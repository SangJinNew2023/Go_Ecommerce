package controllers

import (
	"context"
	"ecommerce/database"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

// create Application struct instance
func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}

}

// Query returns the keyed url query value if it exists, otherwise it returns an empty string `("")`
// Package primitive contains types similar to Go primitives for BSON types that do not have direct Go primitive representations.
func (app *Application) AddToCart() gin.Handler {
	return func(c *gin.Context) {
		//create roductQueryID
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		//createa userQueryID
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		//create productID
		//bjectIDFromHex creates a new ObjectID from a hex string. It returns an error if the hex string is not a valid ObjectID.
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		//create context.WithTimeout
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		//call database.AddProductToCar() with element created above
		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		//IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully added to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		//create roductQueryID
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		//createa userQueryID
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		//create productID
		//bjectIDFromHex creates a new ObjectID from a hex string. It returns an error if the hex string is not a valid ObjectID.
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		//create context.WithTimeout
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//call database.RemoveCartItem() with element created above
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

		//IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully removed item from cart")
	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {

}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("id")

		if userQueryID == "" {
			log.Panicln("user id is empty")
			c.AbortWithError(http.StatusBadRequest, error.New("UserID is empty"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON("Successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		//bjectIDFromHex creates a new ObjectID from a hex string. It returns an error if the hex string is not a valid ObjectID.
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

		//IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully removed item from cart")
	}

}
