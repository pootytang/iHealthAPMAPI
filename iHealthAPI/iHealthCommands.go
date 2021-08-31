package ihealthapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CommandsList []struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// Retrieve all the commands available in the qkview
// Options should contain a the proper page to request from (https://ihealth-api.f5.com/qkview-analyzer/api/qkviews/16642310/commands)
// Return
// 	*ApiResponseStatus
//	CommandsList - This is a slice object already which is a pointer inherently so no need to return a pointer to it
func (api_c *ApiClient) GetCommands(options *ApiRequestOptions) (*ApiResponseStatus, CommandsList) {
	log.Info("GetCommands called")
	// Setup some vars
	clNil := CommandsList{}
	req, err := http.NewRequest("GET", options.Page, nil)
	if err != nil {
		log.Errorf("GetCommands() -> Error creating request object: %s", err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "GetCommands failed to build request object", Error: err, HttpResponse: nil}, clNil
	}

	// Set the headers
	req.Header.Set("Accept", AcceptVND)
	req.Header.Set("Content-Type", ContentJSON)
	req.Header.Set("Uses-Agent", UserAgent)

	// Make the request
	res := api_c.sendRequest(req, 0)
	if res.Error != nil {
		log.Errorf("GetCommands() -> Problem with the request to the commands endpoint for qkview id %s: %s", options.QKViewID, err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "failed to send commands list request", Error: err, HttpResponse: nil}, clNil
	}

	defer res.HttpResponse.Body.Close()

	// We eitehr get a 200 response or a 404 (in the case where the qkview id is wrong)
	var cl CommandsList
	if res.Code == http.StatusOK {
		log.Info("GetCommands() -> Successfully retrieved the commands")
		if err := json.NewDecoder(res.HttpResponse.Body).Decode(&cl); err != nil {
			log.Errorf("GetCommands() -> Problem decoding commands list response: Status = %d, Error = %s", res.Code, err)
			return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Failed to parse command list response", Error: err, HttpResponse: nil}, clNil
		}

		// AHHHHHH YEAH!!!
		return res, cl
	} else {
		log.Warnf("GetCommands() -> There was a problem retrieving commands: Status %d", res.Code)
		ars := &ApiResponseStatus{Code: res.Code, Message: "negative response for commands", Error: err, HttpResponse: nil}
		return ars, clNil
	}
}
