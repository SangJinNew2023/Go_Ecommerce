package controllers

import (
	"context"
	"ecommerce/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

//GenerateFromPassword returns the bcrypt hash of the password at the given cost
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) //byte타입의 array로 password를 입력받아 cost(14)로 암호화후 반환
	if err != nil {
		log.Panic(err)
	}
	return string(bytes) //bytes를 string으로 변환 후 반환
}

//password 검증 
func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.compareHashAndPassword([]byte(givenPassowrd), []byte(userPassowrd))
	valid := true
	msg := ""

	if err != nil {
		msg = "Lgin or Password is incorrect"
		valid = false
	}
	return valid, msg
}


//BindJSON binds the passed struct pointer using the specified binding engine. It will abort the request with HTTP 400 if any error occurs.
//JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
func Signup() gin.HandlerFunc {

	return func(c *gin.Context) { //함수 생성기

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		
		var user models.User

		/// USER model과 gin.context 바인딩
		if err := c.BindJSON(&user); err != nil { 
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}


		//Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.	
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		//CountDocuments returns the number of documents in the collection with filter elements
		//D is a slice and M is a map. 
		//bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}, bson.M{"foo": "bar", "hello": "world", "pi": 3.14159}
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return 
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user already exists"})
		}

		//CountDocuments returns the number of documents in the collection with filter elements
		//D is a slice and M is a map. 
		//bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}, bson.M{"foo": "bar", "hello": "world", "pi": 3.14159}
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

		//user.Password를 hash 암호화 후 다시 저장
		password := HashPassword(*user.Password)
		user.Password = &password
		
		//time.Parse() parses a formatted string and returns the time value it represents.
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		//create token form user data
		token, refreshtoekn, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)

		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.userCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		//InsertOne() insert one document, 필요한 정보를 user에 더한 후 UserCollection DB에 추가
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

		var user models.User

		/// USER model과 gin.context 바인딩
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : err})
			return
		}

//FindOne()은 bson.M{"email": user.Email}에 해당하는 데이터를 Usercollection으로 부터 찾고 decode후 &founduser에 저장
		err := Usercollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"login or password incorrect"})
			return
		}

//VerifyPassword로 유저에게 입력받은 *founduser.Password 와 *user.Password 유효성 검사
		passwordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password))
		defer cancel()

		//password가 이상있으면 에러 msg반환
		if !PasswordIsValid { 
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		//password가 이상 없으면 token 생성
		token, refershToken, _ := generate.TokenGenerator(*founderuser.Email, *founduser.First_name, *founduser.Last_name, founduser.User_ID)
		defer cancel()

		// renews the user tokens when they login
		generate.UpdateAllTokens(token, refreshToken, founderuser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}

}


func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()

        //Reading all data from a collection consists of making the request, then working with the results cursor
		//Find()함수 사용시 element로 빈 bson.D{{}}를 사용하면 모든 data를 읽어 온다는 의미
		//D is a slice and M is a map.
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Someting went wrong, please try after some time")
			return
		}

		//All() iterates the cursor and decodes each document into results(&productlist). 
		//Find()를 통해 searchquery에 담은 결과를 decode method를 이용해 Go type으로 decode
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		defer cursor.Close()
		if err := cursor.err(); err != nil { 
			// Don't forget to log errors. I log them really simple here just
			// to get the point across.
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}

		defer cancel()
		
		c.IndentedJSON(200, productlist)
	}

}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

		var searchProducts []models.Product

		//Query returns the keyed url query value if it exists, otherwise it returns an empty string `("")`.
		queryParam := c.Query("name")

		// you want to check if it's empty
		// c.Header writes a header in the response. If value == "", this method removes the header `c.Writer.Header().Del(key)`
		// c.JSON creates a properly formatted JSON
		// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//Find() executes a find command and returns a Cursor over the matching multiple documents in the collection.
		//The filter parameter in Find() must be a document containing query operators and can be used to select which documents are included in the result.
		//D is a slice and M is a map. 
		//bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}, bson.M{"foo": "bar", "hello": "world", "pi": 3.14159}
		//$regex operator is used to search a specific string in the documents.
		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}})

		if err != nil {
			c.IndentedJSON(404, "something wen wrong while fetching the data")
			return
		}

		//All() iterates the cursor and decodes each document into results(&searchproducts). 
		//Find()를 통해 searchquery에 담은 결과를 decode method를 이용해 Go type으로 decode
		err = searchquerydb.All(ctx, &searchproducts)

		//IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		//Err() returns the last error seen
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchproducts)
	}
}
