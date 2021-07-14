package config

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/requests"
)

// Checks whether each request in config file has valid data
// Creates unique ids for each request using math/rand
func validateAndCreateIdsForRequests(reqs []requests.RequestConfig) []requests.RequestConfig {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	//an array of ids used by database pacakge to calculate mean response time and send notifications
	ids := make(map[int]int64, 0)

	//an array of new requests data after updating the ids
	newreqs := make([]requests.RequestConfig, 0)

	for i, requestConfig := range reqs {
		err := requestConfig.Validate()
		if err != nil {
			log.Fatal().Err(err).Msgf("Invalid Request data in config file for Request #%d %s", i, requestConfig.URL)
		}

		//Set a random value as id
		randInt := random.Intn(1000000)
		ids[randInt] = requestConfig.ResponseTime
		requestConfig.SetID(randInt)
		newreqs = append(newreqs, requestConfig)
	}

	return newreqs
}
