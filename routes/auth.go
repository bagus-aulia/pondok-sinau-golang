package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/danilopolani/gocialite/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// RedirectHandler to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GITHUB"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GITHUB"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_GOOGLE"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GOOGLE"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// CallbackHandler callback of provider
func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// if provider == "google" {
	// 	newMember := getOrRegisterMember(provider, user)
	// 	jwtToken := createTokenMember(&newMember)

	// 	c.JSON(200, gin.H{
	// 		"data":    newMember,
	// 		"token":   jwtToken,
	// 		"user":    user,
	// 		"message": "login success",
	// 	})
	// } else {
	newAdmin := getOrRegisterAdmin(provider, user)
	jwtToken := createTokenAdmin(&newAdmin)

	c.JSON(200, gin.H{
		"data":    newAdmin,
		"token":   jwtToken,
		"user":    user,
		"message": "login success",
	})
	// }

	// Print in terminal user information
	// fmt.Printf("%#v", user)
}

func getOrRegisterAdmin(provider string, admin *structs.User) models.Admin {
	var userData models.Admin
	username := admin.Username

	config.DB.Where("provider = ? AND social_id = ?", provider, admin.ID).First(&userData)

	if username == "" {
		username = admin.ID
	}

	if userData.ID == 0 {
		newAdmin := models.Admin{
			FullName: admin.FullName,
			Username: username,
			Email:    admin.Email,
			SocialID: admin.ID,
			Provider: provider,
			Avatar:   admin.Avatar,
		}

		config.DB.Create(&newAdmin)
		return newAdmin
	}

	return userData
}

func getOrRegisterMember(provider string, member *structs.User) models.Member {
	var userData models.Member
	username := member.Username

	config.DB.Where("provider = ? AND social_id = ?", provider, member.ID).First(&userData)

	if username == "" {
		username = member.ID
	}

	if userData.ID == 0 {
		newMember := models.Member{
			FullName: member.FullName,
			Username: username,
			Email:    member.Email,
			SocialID: member.ID,
			Provider: provider,
			Avatar:   member.Avatar,
		}

		config.DB.Create(&newMember)

		return newMember
	}

	return userData
}

func createTokenAdmin(user *models.Admin) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": user.IsAdmin,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func createTokenMember(user *models.Member) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": "member",
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}
