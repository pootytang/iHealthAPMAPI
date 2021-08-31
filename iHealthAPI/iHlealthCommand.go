package ihealthapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CommandList []struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
	Output string `json:"output"`
}

// Retrieve all the commands available in the qkview
// Options should contain a qkview id (https://ihealth-api.f5.com/qkview-analyzer/api/qkviews/16642310/commands)
// Return
// 	*ApiResponseStatus
//	CommandList - This is a slice object already which is a pointer inherently so no need to return a pointer to it
func (api_c *ApiClient) GetCommand(options *ApiRequestOptions) (*ApiResponseStatus, CommandList) {
	log.Info("GetCommand called")
	// Setup some vars
	clNil := CommandList{}
	req, err := http.NewRequest("GET", options.Page, nil)
	if err != nil {
		log.Errorf("GetCommand() -> Error creating request object: %s", err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Failed to build request object", Error: err, HttpResponse: nil}, clNil
	}

	// Set the headers
	req.Header.Set("Accept", AcceptVND)
	req.Header.Set("Content-Type", ContentJSON)
	req.Header.Set("Uses-Agent", UserAgent)

	// Make the request
	res := api_c.sendRequest(req, 0)
	if res.Error != nil {
		log.Errorf("GetCommand() -> Problem with the request to the command endpoint for page %s: %s", options.Page, err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "failed to send command list request", Error: err, HttpResponse: nil}, clNil
	}

	defer res.HttpResponse.Body.Close()

	// We eitehr get a 200 response or a 404 (in the case where the qkview id is wrong)
	var cl CommandList
	if res.Code == http.StatusOK {
		log.Info("GetCommand() -> Successfully retrieved the command list")
		if err := json.NewDecoder(res.HttpResponse.Body).Decode(&cl); err != nil {
			log.Errorf("GetCommand() -> Problem decoding the command list response: Status = %d, Error = %s", res.Code, err)
			return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Failed to parse command list response", Error: err, HttpResponse: nil}, clNil
		}

		// AHHHHHH YEAH!!!
		return res, cl
	} else {
		log.Warnf("GetCommand() -> There was a problem retrieving commands: Status %d", res.Code)
		ars := &ApiResponseStatus{Code: res.Code, Message: "negative response for command list", Error: err, HttpResponse: nil}
		return ars, clNil
	}
}
