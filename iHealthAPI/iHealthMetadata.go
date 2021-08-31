package ihealthapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

/*
	Ignoring these attributes:
		Files           string      `json:"files"`
		Commands        string      `json:"commands"`
		Bigip           string      `json:"bigip"`
		PrimaryBladeURI interface{} `json:"primary_blade_uri"`
		Upload         struct {
			PerformedBy struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"performed_by"`
			Date int64 `json:"date"`
		} `json:"upload"`
		ProcessingStatus   string      `json:"processing_status"`
		ProcessingMessages string      `json:"processing_messages"`
		Diagnostics        string      `json:"diagnostics"`
		SecondaryBladeURI  interface{} `json:"secondary_blade_uri"`
*/
type Metadata struct {
	GuiURI         string              `json:"gui_uri"`
	ChassisSerial  string              `json:"chassis_serial"`
	Hostname       string              `json:"hostname"`
	VisibleInGui   string              `json:"visible_in_gui"`
	Description    string              `json:"description"`
	F5SupportCase  string              `json:"f5_support_case"`
	Entitlement    MetadataEntitlement `json:"entitlement"`
	GenerationDate int64               `json:"generation_date"`
	ExpirationDate int64               `json:"expiration_date"`
	FileSize       int                 `json:"file_size"`
	FileName       string              `json:"file_name"`
}

type MetadataEntitlement struct {
	ExpirationDate interface{} `json:"expiration_date"`
	DaysLeft       interface{} `json:"days_left"`
}

// GetMetadata grabs Metadata for a given qkview id at the iHealth API's BaseURL
// I prepare a json response that eventually goes to the client here because this
// is where the body of the response is closed
// Returns:
//	*ApiResponseStatus - containing info about the response
//	*AuthResult - containing info about the result of binding and data including json response
func (api_c *ApiClient) GetMetadata(options *ApiRequestOptions) (*ApiResponseStatus, *Metadata) {
	// Prepare some stuff for later use
	mdNIL := &Metadata{
		GuiURI:         "",
		ChassisSerial:  "",
		Hostname:       "",
		VisibleInGui:   "",
		Description:    "",
		F5SupportCase:  "",
		Entitlement:    MetadataEntitlement{ExpirationDate: nil, DaysLeft: nil},
		GenerationDate: 0,
		ExpirationDate: 0,
		FileSize:       0,
		FileName:       "",
	}

	req, err := http.NewRequest("GET", options.Page, nil)
	if err != nil {
		log.Errorf("GetMetadata() -> Error creating request object: %s", err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "GetMetadata failed to build request object", Error: err, HttpResponse: nil}, mdNIL
	}

	req.Header.Set("Accept", AcceptVND)

	// make the request to the backend api. I pass WaitSeconds here because Metadata request may take time processing the response
	res := api_c.sendRequest(req, options.WaitSeconds)
	if res.Error != nil {
		log.Errorf("GetMetadata() -> Problem communicating with the metadata page for qkview id %s: %s", options.QKViewID, err)
		return &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "GetMetadata failed to send request", Error: err, HttpResponse: nil}, mdNIL
	}

	defer res.HttpResponse.Body.Close()

	// prepare Metadata object for success case
	var mdSuccess Metadata

	// sendRequest checked many errors already so just want to check for good or bad only here
	if res.Code == http.StatusOK {
		// KIND OF HAPPY PATH
		log.Info("GetMetadata() -> GetMetadata was successful")
		if err := json.NewDecoder(res.HttpResponse.Body).Decode(&mdSuccess); err != nil {
			log.Errorf("GetMetadata() -> Problem decoding success response body into json object: %s", err)
			ars := &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Problem decoding success response body", Error: err, HttpResponse: nil}
			return ars, mdNIL
		}
		// "look at all the pretty colored hats"
		return res, &mdSuccess
	} else {
		// KIND OF SAD PATH
		// the api in the case of a 302 to the login page does not return json so no need to create a json object here
		log.Warnf("GetMetadata() -> There was a problem retrieving metadata: Status %d", res.Code)
		ars := &ApiResponseStatus{Code: res.Code, Message: "Problem decoding failed response body", Error: err, HttpResponse: nil}
		return ars, mdNIL
	}
}
