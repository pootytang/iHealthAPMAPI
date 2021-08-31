package ihealthapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type FilesList []struct {
	ID           string `json:"id"`
	Permission   string `json:"permissions"`
	Size         int    `json:"size"`
	LastModified string `json:"lastModified"`
	Value        string `json:"value"`
}

/* EXAMPLE
[
    {
        "id": "bW9kX3NoYXJlZC54bWw",
        "lastModified": "Aug 14 2021 08:25",
        "permissions": "0644",
        "size": 126,
        "value": "/mod_shared.xml"
    },
    {
        "id": "c3RhdF9tb2R1bGUueG1s",
        "lastModified": "Aug 14 2021 08:25",
        "permissions": "0644",
        "size": 7482862,
        "value": "/stat_module.xml"
    },
    {
        "id": "bW9kX2Y1b3B0aWNzLnhtbA",
        "lastModified": "Aug 14 2021 08:25",
        "permissions": "0644",
        "size": 128,
        "value": "/mod_f5optics.xml"
    },
    {
        "id": "cG9zdGdyZXNfbW9kdWxlLnhtbA",
        "lastModified": "Aug 14 2021 08:25",
        "permissions": "0644",
        "size": 1979,
        "value": "/postgres_module.xml"
    }
*/

// GetFiles retrieves the various files available in the qkview
// The options should contain the page to request from
// Returns:
// 	*ApiResponseStatus
//	FilessList - This is a slice object pointing to a list of Files
func (api_c *ApiClient) GetFiles(options *ApiRequestOptions) (*ApiResponseStatus, FilesList) {
	log.Info("GetFiles called")
	// Setup some vars
	flNil := FilesList{}
	req, err := http.NewRequest("GET", options.Page, nil)
	if err != nil {
		log.Errorf("GetFiles() -> Error creating request object: %s", err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "GetFiles failed to build request object", Error: err, HttpResponse: nil}, flNil
	}

	// Set the headers
	req.Header.Set("Accept", AcceptVND)
	req.Header.Set("Content-Type", ContentJSON)
	req.Header.Set("Uses-Agent", UserAgent)

	// Make the request
	res := api_c.sendRequest(req, 0)
	if res.Error != nil {
		log.Errorf("GetFiles() -> Problem with the request to the files endpoint for qkview id %s: %s", options.QKViewID, err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "failed to send files list request", Error: err, HttpResponse: nil}, flNil
	}

	defer res.HttpResponse.Body.Close()

	// We eitehr get a 200 response or a 404 (in the case where the qkview id is wrong)
	var fl FilesList
	if res.Code == http.StatusOK {
		log.Info("GetFiles() -> Successfully retrieved the files")
		if err := json.NewDecoder(res.HttpResponse.Body).Decode(&fl); err != nil {
			log.Errorf("GetFiles() -> Problem decoding files list response: Status = %d, Error = %s", res.Code, err)
			return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Failed to parse files list response", Error: err, HttpResponse: nil}, flNil
		}

		// AHHHHHH YEAH!!!
		return res, fl
	} else {
		log.Warnf("GetFiles() -> There was a problem retrieving files: Status %d", res.Code)
		ars := &ApiResponseStatus{Code: res.Code, Message: "negative response for files", Error: err, HttpResponse: nil}
		return ars, flNil
	}
}
