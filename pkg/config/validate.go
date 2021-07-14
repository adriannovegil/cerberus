package config

import (
	"errors"
	"math/rand"
	"net/url"
	"time"

	"devcircus.com/cerberus/pkg/requests"
	"devcircus.com/cerberus/pkg/util/constants"
)

// Checks whether each request in config file has valid data
// Creates unique ids for each request using math/rand
func generateAndAssignIdsForRequests(reqs []requests.RequestConfig) []requests.RequestConfig {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// an array of ids used by database pacakge to calculate mean response time
	// and send notifications
	ids := make(map[int]int64, 0)

	// an array of new requests data after updating the ids
	newreqs := make([]requests.RequestConfig, 0)

	for _, requestConfig := range reqs {
		//Set a random value as id
		randInt := random.Intn(1000000)
		ids[randInt] = requestConfig.ResponseTime
		requestConfig.SetID(randInt)
		newreqs = append(newreqs, requestConfig)
	}

	return newreqs
}

// Validate whether all requestConfig fields are valid
func validate(reqs []requests.RequestConfig) error {

	for _, requestConfig := range reqs {

		if len(requestConfig.URL) == 0 {
			return errors.New("Invalid Url")
		}

		if _, err := url.Parse(requestConfig.URL); err != nil {
			return errors.New("Invalid Url")
		}

		if len(requestConfig.RequestType) == 0 {
			return errors.New("RequestType cannot be empty")
		}

		if requestConfig.ResponseTime == 0 {
			return errors.New("ResponseTime cannot be empty")
		}

		if requestConfig.ResponseCode == 0 {
			requestConfig.ResponseCode = constants.DefaultResponseCode
		}

		if requestConfig.CheckEvery == 0 {
			defTime, _ := time.ParseDuration(constants.DefaultTime)
			requestConfig.CheckEvery = defTime
		}
	}

	return nil
}
