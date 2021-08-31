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

// Files retrieves a list of apm related files
//	When something goes wrong, json containing a list of {"code" : "<status code>", "message" : "status message"}
//		code - message:
//			302 - Need authentication
//			400+ - Bad info provided (bad qkview id). Page couldn't be found. The message will indicate the problem
//			500 - Internal error. Message indicates the internal error
func Files(ctx *gin.Context) {
	log.Info("Files called")

	// Handle the path parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Errorf("Files() -> Received QKView ID of %v is invalid", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "invalid qkview id"})
		return
	}

	log.Debugf("Files() -> Retrieved QKView id: %d", id)

	// Configure the options we need
	opt := ihealthapi.ApiRequestOptions{
		Page:        fmt.Sprintf("%s/%d/files%s", ihealthapi.BASEURL, id, ihealthapi.ReturnJson),
		QKViewID:    id,
		WaitSeconds: 0,
	}

	// Get the client from the gin Context middleware
	client, ok := ctx.MustGet("apiclient").(*ihealthapi.ApiClient)
	if !ok {
		log.Warn("Files() -> Problem Getting api client connection")
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Server error getting connection"})
		return
	}

	// Make the request
	log.Debugf("Files() -> Requesting %s", opt.Page)
	api_res, flist := client.GetFiles(&opt)
	if api_res.Error != nil {
		log.Errorf("Files() -> Error when trying to get the List of files: %s. Status set to: %d", api_res.Error, api_res.Code)
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Server error getting files list"})
		return
	}

	// Check if we're authenticated
	if api_res.Code == http.StatusFound {
		log.Warnf("Files() -> Redirect Found with Code %d", api_res.Code)
		ctx.JSON(http.StatusFound, gin.H{"code": api_res.Code, "message": "Need to authenticate"})
		return
	}

	// Check if something else happened
	if api_res.Code >= http.StatusBadRequest {
		log.Warnf("Files() -> Bad request response returned: Status = %s", api_res.Code)
		ctx.JSON(http.StatusBadRequest, gin.H{"code": api_res.Code, "message": "Communication with files endpoint problem"})
		return
	}

	// All is good so now we need to loop through the files and get the things related to APM
	var apmfiles ihealthapi.FilesList
	for _, item := range flist {
		if strings.Contains(strings.ToLower(item.Value), "apm") {
			apmfiles = append(apmfiles, item)
		}
	}
	ctx.JSON(http.StatusOK, apmfiles)
}
