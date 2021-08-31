package apm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	ihealthapi "github.com/pootytang/iHealthAPMAPI/iHealthAPI"
	log "github.com/sirupsen/logrus"
)

type APMCommands []struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// Commands retrieves the various apm related commands qkview <id> contains
// Context should include the qkview id otherwise an error is thrown
// Returns a list of items:
//	If all is good, json containing a list of {"id" : "<id>", "value" : "<command>"}
//		id = the id of the command
//		value = the command to run
//	When something goes wrong, json containing a list of {"code" : "<status code>", "message" : "status message"}
//		code - message:
//			302 - Need authentication
//			400+ - Bad info provided (bad qkview id). Page couldn't be found. The message will indicate the problem
//			500 - Internal error. Message indicates the internal error
func Commands(ctx *gin.Context) {
	log.Debug("Commands() called")

	// these commands do not have apm in the name but should be checked as well because it affects apm
	// Some of these need additional processing to filter out apm related info. See NOTEs
	/* bigstart memstat, show /sys service all, list /sys software update-status, ls -alsR (critical files) */

	// Handle the path parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Errorf("Path param of %v is invalid", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "invalid qkview id. Should be an integer"})
	}

	// Configure the options we need
	opt := ihealthapi.ApiRequestOptions{
		Page:        fmt.Sprintf("%s/%d/commands%s", ihealthapi.BASEURL, id, ihealthapi.ReturnJson),
		QKViewID:    id,
		WaitSeconds: 0,
	}

	// Get the client from the gin Context middleware
	client, ok := ctx.MustGet("apiclient").(*ihealthapi.ApiClient)
	if !ok {
		log.Warn("Commands() -> Problem Getting api client connection")
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Server error getting connection"})
		return
	}

	// Make the request
	api_res, commands := client.GetCommands(&opt)
	if api_res.Error != nil {
		log.Errorf("Commands() -> Error when trying to get the Commands List: %s. Status set to: %d", api_res.Error, api_res.Code)
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Server error getting command list"})
		return
	}

	// Check if we're authenticated
	if api_res.Code == http.StatusFound {
		log.Warnf("Commands() -> Redirect Found with Code %d", api_res.Code)
		ctx.JSON(http.StatusFound, gin.H{"code": api_res.Code, "message": "Need to authenticate"})
		return
	}

	// Check if something else happened
	if api_res.Code >= http.StatusBadRequest {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": api_res.Code, "message": "Communication with commands endpoint problem"})
		return
	}

	// All is good so now we need to loop through the commands and get the things related to APM
	// if value contains APM add the value and id to the APM slice
	var apmcmds APMCommands
	for _, item := range commands {
		if strings.Contains(strings.ToLower(item.Value), "apm") {
			apmcmds = append(apmcmds, item)
		}
	}
	ctx.JSON(http.StatusOK, apmcmds)
}
