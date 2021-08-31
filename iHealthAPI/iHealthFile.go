package ihealthapi

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// Retrieve the contents of an individual file
// Options should contain a qkview id (https://ihealth-api.f5.com/qkview-analyzer/api/qkviews/16642310/commands)
// Return
// 	*ApiResponseStatus
func (api_c *ApiClient) GetFile(options *ApiRequestOptions) (*ApiResponseStatus, *os.File) {
	log.Info("GetFile called")
	// Setup some vars

	req, err := http.NewRequest("GET", options.Page, nil)
	if err != nil {
		log.Errorf("GetFile() -> Error creating request object: %s", err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Failed to build request object", Error: err, HttpResponse: nil}, nil
	}

	// Set the headers
	req.Header.Set("Accept", AcceptStream)
	req.Header.Set("Content-Type", ContentJSON)
	req.Header.Set("Uses-Agent", UserAgent)

	// Make the request
	res := api_c.sendRequest(req, 0)
	if res.Error != nil {
		log.Errorf("GetFile() -> Problem with the request to the file endpoint for page %s: %s", options.Page, err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "failed to send file list request", Error: err, HttpResponse: nil}, nil
	}

	defer res.HttpResponse.Body.Close()

	// We eitehr get a 200 response or a 404 (in the case where the qkview id is wrong)
	if res.Code == http.StatusOK {
		log.Info("GetFile() -> Successfully retrieved the file")
		log.Debugf("creating or overwriting tempfile: %s", TempLogFile)

		// Write to the temp file
		f, ok := writeToFile(TempLogFile, res.HttpResponse.Body)
		if !ok {
			log.Error("GetFile() -> Problem writing to temp file")
			return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "failed to write to temp file", Error: err, HttpResponse: nil}, nil
		}

		// AHHHHHH YEAH!!!
		return res, f
	} else {
		log.Warnf("GerFile() -> There was a problem retrieving the file: Status %d", res.Code)
		ars := &ApiResponseStatus{Code: res.Code, Message: "negative response retrieving the file", Error: err, HttpResponse: nil}
		return ars, nil
	}
}
