package ihealthapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type AuthResult struct {
	Success AuthSuccess
	Failure AuthFailed
}

// Mimmic a successful authentication response from the backend
type AuthSuccess struct {
	Expires int64 `json:"expires"`
}

// Mimmic a failed authentication response from the backend
type AuthFailed struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Code2   string `json:"code2"`
	Message string `json:"message"`
}

type Credentials struct {
	Username string `json:"user_id"`
	Password string `json:"user_secret"`
}

// Authenticate method uses the credentials passed in to authenticate against the AuthURL
// Responases are:
// 200 (StatusOK) - Everything is happy
// 400 (StatusBadRequest) - For anything 400 and above
// 		401 - Authentication Failure (AuthURL rejected creds)
// 		406 - Incorrect headers sent to the AuthURL (AuthURL may respond with this)
// 		415 - Missing or incorrect headers such as Content-Type not being set to application/json (I have control of this)
// 500 (StatusInternalServerError) - Internal problems
func (api_c *ApiClient) Authenticate(creds *Credentials) (*ApiResponseStatus, *AuthResult) {
	ar := &AuthResult{
		Success: AuthSuccess{Expires: 0},
		Failure: AuthFailed{Errors: nil},
	}
	jsonStr := fmt.Sprintf(`{"user_id":"%s", "user_secret":"%s"}`, creds.Username, creds.Password)
	jsonBytes := []byte(jsonStr)
	req, err := http.NewRequest("POST", api_c.AuthURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("Authenticate() -> Error creating request object: %s", err)
		ars := &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Authenticate failed to build request object", Error: err, HttpResponse: nil}
		return ars, ar
	}

	// Need to make sure Content-Type header is set to avoid a 415 response and the Accept header is set to avoid a 406 response
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", ContentJSON)
	req.Header.Set("Uses-Agent", UserAgent)

	// Make the request
	res := api_c.sendRequest(req, 0)
	if res.Error != nil {
		log.Errorf("Authenticate() -> Problem communicating with the auth page: %s", err)
		ars := &ApiResponseStatus{Code: http.StatusBadRequest, Message: "Problm communicating with the auth page", Error: err, HttpResponse: nil}
		return ars, ar
	}

	defer res.HttpResponse.Body.Close()

	// sendRequest checked many errors already so just want to check for good or bad only here
	var authsuccess AuthSuccess
	var authfailed AuthFailed

	if res.Code == http.StatusOK {
		// The world is beautiful
		log.Info("Authenticate() -> Authentication was successful")
		if err := json.NewDecoder(res.HttpResponse.Body).Decode(&authsuccess); err != nil {
			log.Errorf("Authenticate() -> Problem decoding success response body into json object: %s", err)
			ars := &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Problem decoding response body", Error: err, HttpResponse: nil}
			return ars, ar
		}
		// all is well
		ar.Success = authsuccess
		return res, ar
	} else {
		// The world sucks
		log.Warnf("Authenticate() -> There was a problem authenticating: Status %ds", res.Code)
		if err := json.NewDecoder(res.HttpResponse.Body).Decode(&authfailed); err != nil {
			log.Errorf("Authenticate() -> Problem decoding failure response body into json object: %s", err)
			ars := &ApiResponseStatus{Code: http.StatusInternalServerError, Message: "Problem decoding failure response body", Error: err, HttpResponse: nil}
			return ars, ar
		}
		// all is well except we failed authentication
		ar.Failure = authfailed
		return res, ar
	}
}
