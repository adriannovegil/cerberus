package config

import (
	"errors"
	"net/url"
	"time"

	"devcircus.com/cerberus/pkg/target/request"
	"devcircus.com/cerberus/pkg/util/constant"
)

// Validate whether all requestConfig fields are valid
func validate(reqs []request.Config) error {

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
			requestConfig.ResponseCode = constant.DefaultResponseCode
		}

		if requestConfig.CheckEvery == 0 {
			defTime, _ := time.ParseDuration(constant.DefaultTime)
			requestConfig.CheckEvery = defTime
		}
	}

	return nil
}
