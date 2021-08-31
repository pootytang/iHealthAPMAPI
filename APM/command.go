package apm

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ihealthapi "github.com/pootytang/iHealthAPMAPI/iHealthAPI"
	log "github.com/sirupsen/logrus"
)

// Command retrieves the output of the apm command that command <id> contains
// Context should include the qkview id and command id otherwise an error is returned
// Returns:
//	If all is good, json containing {"id" : "<id>", "name" : "name of commmand", "status" : "int", "output": "<base64>"}
//		id = the id of the command
//		name = the command to run (show /apm profile access all)
//		status = 0 for success, 1 for fail
//		output = base64 string of the output the command returned
//
//	When something goes wrong, json containing {"code" : "<status code>", "message" : "status message"}
//		code - message:
//			302 - Need authentication
//			400+ - Bad info provided (bad qkview or command id). Page couldn't be found. The message will indicate the problem
//			500 - Internal error. Message indicates the internal error
func Command(ctx *gin.Context) {
	log.Info("Command called")

	// Handle the path parameters
	qkviewid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Errorf("QKview ID of %v is invalid", qkviewid)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "invalid qkview id. Should be an integer"})
	}

	// Configure the options we need
	opt := ihealthapi.ApiRequestOptions{
		Page:        fmt.Sprintf("%s/%d/commands/%s%s", ihealthapi.BASEURL, qkviewid, ctx.Param("commandId"), ihealthapi.ReturnJson),
		QKViewID:    qkviewid,
		WaitSeconds: 0,
	}

	// Get the client from the gin Context middleware
	client, ok := ctx.MustGet("apiclient").(*ihealthapi.ApiClient)
	if !ok {
		log.Warn("Command() -> Problem Getting api client connection")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Server error getting connection to command endpoint"})
		return
	}

	// Make the request
	api_res, commands := client.GetCommand(&opt)
	if api_res.Error != nil {
		log.Errorf("Command() -> Error when trying to get the Commands List: %s. Status set to: %d", api_res.Error, api_res.Code)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Server error getting command list"})
		return
	}

	// Check if we're authenticated
	if api_res.Code == http.StatusFound {
		log.Warnf("Command() -> Redirect Found with Code %d", api_res.Code)
		ctx.JSON(http.StatusFound, gin.H{"code": api_res.Code, "message": "Need to authenticate"})
		return
	}

	// Check if something else happened
	if api_res.Code >= http.StatusBadRequest {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": api_res.Code, "message": "Communication with command endpoint problem"})
		return
	}

	// All is good return the list to the client as json
	ctx.JSON(http.StatusOK, commands)
}
