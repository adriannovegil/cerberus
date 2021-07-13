package config

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/requests"
)

// Checks whether each request in config file has valid data
// Creates unique ids for each request using math/rand
func validateAndCreateIdsForRequests(reqs []requests.RequestConfig) ([]requests.RequestConfig, map[int]int64) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	//an array of ids used by database pacakge to calculate mean response time and send notifications
	ids := make(map[int]int64, 0)

	//an array of new requests data after updating the ids
	newreqs := make([]requests.RequestConfig, 0)

	for i, requestConfig := range reqs {
		validateErr := requestConfig.Validate()
		if validateErr != nil {
			log.Info().Msg(fmt.Sprintf("Invalid Request data in config file for Request #%d %s", i, requestConfig.URL))
			log.Info().Msg(fmt.Sprintf("Error: %s", validateErr.Error()))
			os.Exit(3)
		}

		//Set a random value as id
		randInt := random.Intn(1000000)
		ids[randInt] = requestConfig.ResponseTime
		requestConfig.SetID(randInt)
		newreqs = append(newreqs, requestConfig)
	}

	return newreqs, ids
}
