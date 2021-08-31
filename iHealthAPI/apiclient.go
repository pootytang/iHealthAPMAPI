package ihealthapi

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	BASEURL      = "https://ihealth-api.f5.com/qkview-analyzer/api/qkviews"
	AUTHURL      = "https://api.f5.com/auth/pub/sso/login/ihealth-api"
	ReturnJson   = ".json"
	UserAgent    = "f5.com apmclient"
	AcceptVND    = "application/vnd.f5.ihealth.api"
	ContentJSON  = "application/json"
	AcceptStream = "application/octet-stream"
	TempLogFile  = "./Temp/temp.log"
)

type ApiClient struct {
	BaseURL    string
	AuthURL    string
	HTTPClient *http.Client
}

// ApiRequestOptions are options that can be sent in a request
// 		Page = the page to request and is appended to the BaseURL
//		QKViewID = what do you think (maybe you're not sure. This is the id of the qkview generated after uploading it to iHealth)
// 		WaitSeconds = can be used if the server needs more time to process
type ApiRequestOptions struct {
	Page        string
	QKViewID    int
	WaitSeconds int64
}

// Status is a wrapper for an http status integer
type ApiResponseStatus struct {
	Code         int
	Message      string
	Error        error
	HttpResponse *http.Response
}

// Disable redirects according to https://pkg.go.dev/net/http:
/*
	ErrUseLastResponse can be returned by Client.CheckRedirect hooks to control how redirects are processed.
	If returned, the next request is not sent and the most recent response is returned with its body unclosed."
*/
func noRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

// MewClient returns a *ApiClient object preconfigured with some settings
func NewClient() *ApiClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error creating cookiejar: %s", err)
	}
	return &ApiClient{
		BaseURL: BASEURL,
		AuthURL: AUTHURL,
		HTTPClient: &http.Client{
			Jar:           jar,
			Timeout:       time.Minute,
			CheckRedirect: noRedirect,
		},
	}
}

// sendRequest does exactly that. It sends the request passed in to the BaseURL + Options page
func (c *ApiClient) sendRequest(req *http.Request, seconds int64) *ApiResponseStatus {
	// make the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Errorf("sendRequest() -> Problem making request: %s", err)
	}

	// Let set the default ApiResponseStatus. I'm assuming everything went well
	api_res := &ApiResponseStatus{
		Code:         res.StatusCode,
		Message:      "success",
		Error:        err,
		HttpResponse: res,
	}

	// Now lets check if something did go wrong
	if api_res.Error != nil {
		log.Errorf("sendRequest() -> Error communicating with api page: %s", api_res.Error)
		api_res.Message = "Failed sending request"
		return api_res
	}

	//Check the possible result codes
	// Anything marked with "**" the program has control over so it's unlikely to see these codes unless the backend changed
	// 200 - (All) Everything is good
	// 202 - (Metadata) This can happen if the server needs more processing time as documented at https://clouddocs.f5.com/api/ihealth/QKView_Metadata.html
	//
	// 400 - (All) For anything 400 and above
	// 		401 - (Auth) Authentication Failure (AuthURL rejected creds)
	// 		403 - (Metadata) not sure what scenario returns this value
	//		404 - (Metadata) possible problem with the params as the server responds with a Not Found
	// 		406 - (Auth, Metadata) Incorrect headers sent to the AuthURL (AuthURL may respond with this) **
	// 		415 - (Auth) Missing or incorrect headers such as Content-Type not being set to application/json **
	//
	// 500 - (All) Internal problems
	switch sc := api_res.Code; {
	case sc == http.StatusOK:
		log.Infof("sendRequest() -> Pretty flowers and happy unicorns")
		return api_res
	case sc == http.StatusAccepted:
		log.Warnf("202 received, Waiting %d seconds", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)
	case sc >= http.StatusBadRequest:
		log.Errorf("sendRequest() -> Bad request sent to backend: Response code = %d Url = %s", api_res.Code, req.URL)
		api_res.Message = "Bad request sent to backend"
		return api_res
	case sc >= http.StatusMultipleChoices && sc < http.StatusBadRequest:
		log.Warnf("Redirect to Login found. Need to authenticate: Response code = %d", api_res.Code)
		api_res.Message = "Need to authenticate"
		return api_res
	default:
		log.Errorf("sendRequest() -> An Unknown problem happened. Status code set to %d", api_res.Code)
		api_res.Message = "Unknown problem occurred"
	}

	return api_res
}
