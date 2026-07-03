package controllers

import (
	"net/http"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/GdAyo19/FundsFlow/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// bind the request body to the body struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}
	// hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(body.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "password encryption failed",
		})
		return
	}
	// create a new user using the data from the request body and the hashed password
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	// save the user to the database using GORM
	register := config.DB.Create(&user)

	if register.Error != nil {
		c.JSON(400, gin.H{
			"error": "User already exist",
		})
		return
	}

	// return a success response with the user data without the password
	c.JSON(201, gin.H{
		"message": "user created successfully",
		"user": models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})

}

func GetUsers(c *gin.Context) {
	// retrieve all users from the database using GORM
	var users []models.User

	// check if there was an error retrieving the users from the database
	result := config.DB.Find(&users)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	var responses []models.UserResponse

	// not to print the password in the response, we create a new struct to hold the user data without the password
	for _, user := range users {
		response := models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		responses = append(responses, response)
	}

	c.JSON(200, responses)
}

func Login(c *gin.Context) {
	// create a body struct to hold the email and password from the request body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// bind the request body to the body struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"Error": "invalid data",
		})

		return
	}

	// get the user from the database using the email from the request body
	var user models.User
	
	// check if the user exists in the database
	result := config.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "invalid email or password",
		})
		return
	}
	// compare the password from the request body with the hashed password from the database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "invalid email or password",
		})
		return
	}

	// generate a JWT token for the user using the user ID
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// return the token in the response
	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func Profile(c *gin.Context) {
	// getting the userID from the context set by the AuthMiddleware
	userID := c.GetUint("userID")

	var user models.User

	// searching for the user in the database using the userID
	result := config.DB.First(&user, userID)

	if result.Error != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})

		return
	}

	// returning the user data without the password in the response
	c.JSON(http.StatusOK, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
