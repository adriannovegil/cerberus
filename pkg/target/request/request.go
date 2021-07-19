package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	// RequestsList list of requests configuration
	RequestsList []Config
	throttle     chan int
)

const (
	// ContentType attribute
	ContentType = "Content-Type"
	// ContentLength attribute
	ContentLength = "Content-Length"
	// FormContentType attribute
	FormContentType = "application/x-www-form-urlencoded"
	// JSONContentType attribute
	JSONContentType = "application/json"
)

// Config data structure
type Config struct {
	ID           int
	URL          string            `yaml:"url"`
	RequestType  string            `yaml:"requestType"`
	Headers      map[string]string `yaml:"headers"`
	FormParams   map[string]string `yaml:"formParams"`
	URLParams    map[string]string `yaml:"urlParams"`
	ResponseCode int               `yaml:"responseCode"`
	ResponseTime int64             `yaml:"responseTime"`
	CheckEvery   time.Duration     `yaml:"checkEvery"`
	Fallbacks    []string          `yaml:"fallbacks"`
}

// SetID set Id for request
func (requestConfig *Config) SetID(id int) {
	requestConfig.ID = id
}

// Init Initialize data from config file and check all requests
//func Init(data []RequestConfig, concurrency int) {
//	RequestsList = data
//
//	//throttle channel is used to limit number of requests performed at a time
//	if concurrency == 0 {
//		throttle = make(chan int, DefaultConcurrency)
//	} else {
//		throttle = make(chan int, concurrency)
//	}
//
//	if len(data) == 0 {
//		println("\nNo requests to monitor.Please add requests to you config file")
//		os.Exit(3)
//	}
//	//send requests to make sure every every request is valid
//	println("\nSending requests to apis.....making sure everything is right before we start monitoring")
//	println("Api Count: ", len(data))
//
//	for i, requestConfig := range data {
//		println("Request #", i, " : ", requestConfig.RequestType, " ", requestConfig.URL)
//
//		//Perform request
//		reqErr := PerformRequest(requestConfig, nil)
//
//		if reqErr != nil {
//			//Request Failed
//			println("\nFailed !!!! Not able to perfome below request")
//			println("\n----Request Deatails---")
//			println("Url :", requestConfig.URL)
//			println("Type :", requestConfig.RequestType)
//			println("Error Reason :", reqErr.Error())
//			println("\nPlease check the config file and try again")
//			os.Exit(3)
//		}
//	}
//
//	println("All requests Successfull")
//}

// PerformRequest takes the date from requestConfig and creates http request and executes it
func PerformRequest(requestConfig Config, throttle chan int) error {
	//Remove value from throttel channel when request is completed
	defer func() {
		if throttle != nil {
			<-throttle
		}
	}()

	var request *http.Request
	var reqErr error

	if len(requestConfig.FormParams) == 0 {
		//formParams create a request
		request, reqErr = http.NewRequest(requestConfig.RequestType,
			requestConfig.URL,
			nil)

	} else {
		if requestConfig.Headers[ContentType] == JSONContentType {
			//create a request using using formParams

			jsonBody, jsonErr := GetJSONParamsBody(requestConfig.FormParams)
			if jsonErr != nil {
				return jsonErr
			}
			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.URL,
				jsonBody)

		} else {
			//create a request using formParams
			formParams := GetURLValues(requestConfig.FormParams)

			request, reqErr = http.NewRequest(requestConfig.RequestType,
				requestConfig.URL,
				bytes.NewBufferString(formParams.Encode()))

			request.Header.Add(ContentLength, strconv.Itoa(len(formParams.Encode())))

			if requestConfig.Headers[ContentType] != "" {
				//Add content type to header if user doesnt mention it config file
				//Default content type application/x-www-form-urlencoded
				request.Header.Add(ContentType, FormContentType)
			}
		}
	}

	if reqErr != nil {
		return reqErr
	}

	//add url parameters to query if present
	if len(requestConfig.URLParams) != 0 {
		urlParams := GetURLValues(requestConfig.URLParams)
		request.URL.RawQuery = urlParams.Encode()
	}

	//Add headers to the request
	AddHeaders(request, requestConfig.Headers)

	//TODO: put timeout ?
	/*
		timeout := 10 * requestConfig.ResponseTime

		client := &http.Client{
			Timeout: timeout,
		}
	*/

	client := &http.Client{}
	//start := time.Now()

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		//Request failed .
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != requestConfig.ResponseCode {
		return errResposeCode(getResponse.StatusCode, requestConfig.ResponseCode)
	}

	//elapsed := time.Since(start)

	return nil
}

// convertResponseToString convert response body to string
func convertResponseToString(resp *http.Response) string {
	if resp == nil {
		return " "
	}
	buf := new(bytes.Buffer)
	_, bufErr := buf.ReadFrom(resp.Body)

	if bufErr != nil {
		return " "
	}

	return buf.String()
}

// AddHeaders add header values from map to request
func AddHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

// GetURLValues convert params in map to url.Values
func GetURLValues(params map[string]string) url.Values {
	urlParams := url.Values{}
	i := 0
	for key, value := range params {
		if i == 0 {
			urlParams.Set(key, value)
		} else {
			urlParams.Add(key, value)
		}
	}

	return urlParams
}

// GetJSONParamsBody creates body for request of type application/json from map
func GetJSONParamsBody(params map[string]string) (io.Reader, error) {
	data, jsonErr := json.Marshal(params)

	if jsonErr != nil {

		jsonErr = errors.New("Invalid Parameters for Content-Type application/json : " + jsonErr.Error())

		return nil, jsonErr
	}

	return bytes.NewBuffer(data), nil
}

// errResposeCode creates an error when response code from server is not equal to response code mentioned in config file
func errResposeCode(status int, expectedStatus int) error {
	return fmt.Errorf("Got Response code %v. Expected Response Code %v ", status, expectedStatus)
}
