package engine

import (
	"encoding/json"
	"github.com/project-flogo/core/data/schema"
	"io/ioutil"
	"os"

	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/engine/secret"
	"github.com/project-flogo/core/support"
)

var appName, appVersion string

func init() {
	//todo check and enable schema validation before loading the application
	if IsSchemaValidationEnabled() {
		schema.Enable()
	}
}

// Returns name of the application
func GetAppName() string {
	return appName
}

// Returns version of the application
func GetAppVersion() string {
	return appVersion
}

func LoadAppConfig(flogoJson string, compressed bool) (*app.Config, error) {

	var jsonBytes []byte

	if flogoJson == "" {

		// a json string wasn't provided, so lets lookup the file in path
		configPath := GetFlogoConfigPath()

		flogo, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}

		jsonBytes, err = ioutil.ReadAll(flogo)
		if err != nil {
			return nil, err
		}
	} else {

		if compressed {
			var err error
			jsonBytes, err = support.DecodeAndUnzip(flogoJson)
			if err != nil {
				return nil, err
			}
		} else {
			jsonBytes = []byte(flogoJson)
		}
	}

	updated, err := secret.PreProcessConfig(jsonBytes)
	if err != nil {
		return nil, err
	}

	appConfig := &app.Config{}
	err = json.Unmarshal(updated, &appConfig)
	if err != nil {
		return nil, err
	}

	appName = appConfig.Name
	appVersion = appConfig.Version

	return appConfig, nil
}
