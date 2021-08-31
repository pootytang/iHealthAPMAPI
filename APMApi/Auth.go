package apmapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ihealthapi "github.com/pootytang/iHealthAPMAPI/iHealthAPI"
	log "github.com/sirupsen/logrus"
)

// Auth endpoint where you need to authenticate before makimg "real" requests for data
// Post credentials in json format example: {"user_id": "foo@f5.com", "user_secret": "foobarpass"}
// Responses:
//   {"Message": "Success"} - indicates successfully authenticated and now you can make real requests
//   {"message" : "user_id or user_secret missing"} - self explanatory
func Auth(ctx *gin.Context) {
	var creds ihealthapi.Credentials
	err := ctx.BindJSON(&creds)
	if err != nil {
		log.Warnf("Unable to bind credentials: %s", err)
		ctx.JSON(500, gin.H{"message": "Something happend, the server was unable to parse the credentials"})
		return
	}

	// Check password and username is not empty
	if creds.Password == "" || creds.Username == "" {
		if creds.Password == "" {
			log.Warn("Password is empty")
		}
		if creds.Username == "" {
			log.Warn("Username is empty")
		}
		log.Debug("Calling ctx.JSON NOW!!!")
		ctx.JSON(http.StatusOK, gin.H{"message": "user_id or user_secret missing"})
		return
	}

	log.Infof("Credentials are populated for %s", creds.Username)

	// Grab the client from the gin context middleware
	client, ok := ctx.MustGet("apiclient").(*ihealthapi.ApiClient)
	if !ok {
		log.Error("getMetadata() -> Problem Getting api client connection")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Server error getting connection"})
		return
	}

	// Make the request and get the ApiResultStatus object
	api_res, auth_res := client.Authenticate(&creds)
	if api_res.Error != nil || api_res.Code > http.StatusOK {
		log.Error("Auth() -> error calling Authenticate: Status set to: %d", api_res.Code)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Auth Failure"})
		return
	} else { // Successful Authentication
		t := time.Unix(auth_res.Success.Expires, 0)
		log.Infof("Auth() -> Authentication was successful. Credentials expire at: %s", t)
		ctx.JSON(http.StatusOK, gin.H{"message": "Authentication succeeded", "utcsessionend": t})
		return
	}
}
