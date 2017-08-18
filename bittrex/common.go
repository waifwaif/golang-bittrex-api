package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func init() {
	defaultConnOpts = ConnectOptions{}
}

// BaseURL represents the base URL for all requests
const (
	BaseURL    = "https://bittrex.com/api"
	Private    = "auth"
	Public     = "pub"
	APIVersion = "2.0"
)

// Auth represents the auth credentials to authenticate to the Bittrex API:
//
// It consists of a set of a private and a public key.
type Auth struct {
	PublicKey  string // The public key to connect to bittrex API.
	PrivateKey string // The private key to connect to bittrex API.
}

// publicParams represents the possible public parameters that can
// be passed for public Api calls.
type publicParams struct {
	MarketName   *string
	TickInterval *string
}

// AddToQueryString adds the non empty fields of the publicParams struct
// to the specified query string.
func (pp publicParams) AddToQueryString(queryString url.Values) {
	if queryString != nil {
		if pp.MarketName != nil {
			queryString.Add("marketName", *pp.MarketName)
		}
		if pp.TickInterval != nil {
			queryString.Add("tickInterval", *pp.TickInterval)
		}
	}
}

// privateParams represents the possible private parameters that can
// be passed for auth API calls.
type privateParams struct {
}

// AddToPostForm adds the non empty fields of the publicParams struct
// to the specified post form.
func (pp privateParams) AddToPostForm(postForm url.Values) {
	if postForm != nil {

	}
}

// defaultConnOpts represents the default configuration for ConnectOptions.
var defaultConnOpts ConnectOptions

// ConnectOptions represents custom Connect
// Configurations for HTTP requests.
type ConnectOptions struct {
	hmacSignature interface{}
	Auth          Auth
	ConnTimeout   time.Duration
}

// checkOptions checks the specified options and sets them to the default values.
func checkOptions() *ConnectOptions {
	return nil
}

// apiCall performs a generic API call.
func apiCall(Version, Visibility, Entity, Feature string, GetParameters *publicParams, PostParameters *privateParams, options *ConnectOptions) (*json.RawMessage, error) {
	client := http.Client{
		Timeout: options.ConnTimeout,
	}
	URL := fmt.Sprintf("%s/v%s/%s/%s/%s", BaseURL, Version, Visibility, Entity, Feature)
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Add("Accept", "application/json")

	if err != nil {
		return nil, err
	}

	if Visibility == Public && GetParameters != nil { // Add them to query string
		GetParameters.AddToQueryString(req.URL.Query())
	} else if Visibility == Private {
		addSecurityHeaders(req.Header)
		if PostParameters != nil {
			PostParameters.AddToPostForm(req.PostForm)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status Code: %s", err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret response
	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Success == false {
		return nil, fmt.Errorf("Error Response: %s", ret.Message)
	}

	return ret.Result, nil
}

// publicCall performs a call to the public bittrex API.
//
// It does not need API Keys.
func publicCall(Entity, Feature string, GetParameters *publicParams, options *ConnectOptions) (*json.RawMessage, error) {
	options = checkOptions(options)
	return apiCall(APIVersion, Public, Entity, Feature, GetParameters, nil, options)
}

// authCall performs a call to the private bittrex API.
//
// It needs an Auth struct to be passed with valid Keys.
func authCall(Entity, Feature string, PostParams *privateParams, options *ConnectOptions) (*json.RawMessage, error) {
	//options = checkOptions()
	if options.Auth.PublicKey == "" || options.Auth.PrivateKey == "" {
		return nil, errors.New("Cannot perform private api requst without authentication keys")
	}
	//createHMAC signature
	return apiCall(APIVersion, Private, Entity, Feature, nil, PostParams, options)
}

// addSecurityHeaders adds security headers, required for bittrex private API calls.
//
// Example of this headers which need to be added are HMAC signature.
func addSecurityHeaders(header http.Header) {

}
