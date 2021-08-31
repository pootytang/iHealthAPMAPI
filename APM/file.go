package apm

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ihealthapi "github.com/pootytang/iHealthAPMAPI/iHealthAPI"
	log "github.com/sirupsen/logrus"
)

// APMFile represents the file the caller wants to retrieve
// Caller sends json in the form of:
// 	{"id" : "<qkviewid>", "file_id" : "<file_id>", "file_name" : "<file_name>"}
type APMFile struct {
	ID       string `json:"id"`
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
}

// File retrieves the content of the file requested
// Context should include the qkview id and file id otherwise an error is returned
// Returns:
//	If all is good, json containing {"id" : "<fileid>", "name" : "name of file", "content": "<base64>"}
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
func File(ctx *gin.Context) {
	log.Info("File called")

	// Handle the json posted
	var af APMFile
	err := ctx.BindJSON(&af)
	if err != nil {
		log.Warnf("File() -> Unable to bind the requested apm file: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Error binding the requested apm file"})
		return
	}

	// Handle the path parameters
	qkviewid, err := strconv.Atoi(af.ID)
	if err != nil {
		log.Errorf("QKview ID of %v is invalid", qkviewid)
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "invalid qkview id. Should be an integer"})
		return
	}

	// Make sure the File id and File name are not empty
	if af.FileID == "" || af.FileName == "" {
		log.Error("file_id or file_name is invalid. Both must be populated")
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid parameters"})
		return
	}

	log.Debugf("Files() -> Retrieved QKView id: %d, File id: %s, FileName: %s", qkviewid, af.FileID, af.FileName)

	// Configure the options we need
	opt := ihealthapi.ApiRequestOptions{
		Page:        fmt.Sprintf("%s/%d/files/%s", ihealthapi.BASEURL, qkviewid, ctx.Param("fileId")),
		QKViewID:    qkviewid,
		WaitSeconds: 0,
	}

	// Get the client from the gin Context middleware
	client, ok := ctx.MustGet("apiclient").(*ihealthapi.ApiClient)
	if !ok {
		log.Warn("FIle() -> Problem Getting api client connection")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Server error getting connection to the file endpoint"})
		return
	}

	// Make the request
	api_res, fl := client.GetFile(&opt)
	if api_res.Error != nil {
		log.Errorf("File() -> Error when trying to get the File: %s. Status set to: %d", api_res.Error, api_res.Code)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Server error getting the file: %s", ctx.Param("fileId"))})
		return
	}

	// Check if we're authenticated
	if api_res.Code == http.StatusFound {
		log.Warnf("File() -> Redirect Found with Code %d", api_res.Code)
		ctx.JSON(http.StatusFound, gin.H{"code": api_res.Code, "message": "Need to authenticate"})
		return
	}

	// Check if something else happened
	if api_res.Code >= http.StatusBadRequest {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": api_res.Code, "message": "Problem communicating with the File endpoint"})
		return
	}

	// All is good return the list to the client as json
	// TODO: Maybe return the session id's or from another function, return the path a session took. No session id found return a notfound status code
	ctx.JSON(http.StatusOK, fl)
}
