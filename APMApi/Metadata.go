package apmapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ihealthapi "github.com/pootytang/iHealthAPMAPI/iHealthAPI"
	log "github.com/sirupsen/logrus"
)

// getMetadata retrieves metadata from the qkview with the provided ID
// Returns:
//		200 (http.StatusOK) - All good
//		400 (http.StatusBadRequest) - for any 400 code (API lists 403, 404, 406)
//		500 (http.StatusInternalServerError) - something bad happened internally
func Metadata(ctx *gin.Context) {
	log.Debug("Metadata called")

	// Handle the path parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Errorf("Path param of %v is invalid", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "invalid qkview id. Should be an integer"})
	}

	opt := ihealthapi.ApiRequestOptions{
		Page:        fmt.Sprintf("%s/%d/commands%s", ihealthapi.BASEURL, id, ihealthapi.ReturnJson),
		QKViewID:    id,
		WaitSeconds: 10,
	}

	// Get the client from the gin Context middleware
	client, ok := ctx.MustGet("apiclient").(*ihealthapi.ApiClient)
	if !ok {
		log.Warn("Metadata() -> Problem Getting api client connection")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Server error getting connection"})
		return
	}

	// Make the request
	api_res, md := client.GetMetadata(&opt)
	if api_res.Error != nil {
		log.Errorf("Metadata() -> Error calling GetMetadata: %s. Status set to: %d", api_res.Error, api_res.Code)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Server error getting metadata"})
		return
	}

	if api_res.Code == http.StatusFound {
		log.Warnf("%d response code found", api_res.Code)
		ctx.JSON(http.StatusFound, gin.H{"code": api_res.Code, "message": "Need to authenticate"})
		return
	}

	if api_res.Code >= http.StatusBadRequest {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": api_res.Code, "message": "Communication with metadata endpoint problem"})
		return
	}

	// If we made it this far all is good so return the Metadata
	ctx.JSON(http.StatusOK, md)
	return
}
