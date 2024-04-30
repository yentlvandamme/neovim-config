package distribution

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)
type Config struct {
	Name        string
	Version     string
	Description string
}

type Distribution struct {
	path   string
	config Config
}

func FindDist(name string) (Distribution, error) {
	distribution := Distribution{}
	config := Config{}

	currentDir, getwdErr := os.Getwd()
	if getwdErr != nil {
		return distribution, getwdErr
	}

	distsPath := currentDir + "/distributions"
	dists, readDirErr := os.ReadDir(distsPath)
	if readDirErr != nil {
		return distribution, readDirErr
	}

	for _, dist := range dists {
		if dist.IsDir() {
			distPath := distsPath + string(os.PathSeparator) + dist.Name()
			manifestUri := distPath + string(os.PathSeparator) + "manifest.json"
			if _, fileNotExistErr := os.Stat(manifestUri); errors.Is(fileNotExistErr, os.ErrNotExist) {
				continue
			}

			data, readFileErr := os.ReadFile(manifestUri)
			if readFileErr != nil {
				return distribution, readFileErr
			}

			if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
				return distribution, unmarshalErr
			}

			if config.Name == name {
				distribution.config = config
				distribution.path = distPath
				return distribution, nil
			}
		}
	}

	return distribution, fmt.Errorf("Could not find matching configuration")
}


