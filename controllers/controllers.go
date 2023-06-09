package controllers

import (
	"context"
	"ecommerce/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}
//BindJSON binds the passed struct pointer using the specified binding engine. It will abort the request with HTTP 400 if any error occurs.
//JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		c.BindJSON(&user); err != nil { // USER model과 바인딩
			c.JSON{http.StatusBadRequest, gin.H{"error": err.Error()}}
			return
		}
//Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H("error": validationErr))
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError.gin.H{"error": err})
			return 
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return 
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone  number is alread in use"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		
		//parses a formatted string and returns the time value it represents.
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshtoekn, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)

		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.userCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}
		defer cancel()
		
		c.JSON(http.StatusCreated, "Successfully signed in!")

	}
}

func Login() gin.HandlerFunc {
	return func( c *gin.Context){
		var ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.Userc
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : err})
			return
		}

		err := Usercollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"login or password incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password))
		defer cancel()

		if !PasswordIsValid {
			c.JSON{http.StatusInternalServerError, gin.H{"error": msg}}
			fmt.Println(msg)
			return
		}
		token, refershToken, _ := generate.TokenGenerator(*founderuser.Email, *founduser.First_name, *founduser.Last_name, founduser.User_ID)
		defer cancel()
		
		generate.UpdateAllTokens(token, refreshToken, founderuser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}

}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

}

func SearchProductByQuery() gin.HandlerFunc {

}
