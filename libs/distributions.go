package configs

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
	Path   string
	Config Config
}

type DistributionManager struct {
    Distributions map[string]Distribution
}

// TODO: Should we consider creating a DistributionManager, which takes care of clearing configs,
// As well as getting the current config path? It'd extend or wrap around Distribution-methods.
// ie. when loading a distribution, we pass the requested distribution name to the manager. That manager
// then looks up the distribution, and executes the distribution's load method

// Should create a distribution manager
func NewDistManager() (*DistributionManager, error) {
    currentDir, err := os.Getwd()
    if err != nil {
        return &DistributionManager{}, err
    }

    distsPath := currentDir + "/distributions"
    dists, err := os.ReadDir(distsPath)
	if err != nil {
		return &DistributionManager{}, err
	}
    if err != nil {
        return &DistributionManager{}, err
    }

    distributions := make(map[string]Distribution)
    for _, dist := range dists {
		if dist.IsDir() {
			distPath := distsPath + string(os.PathSeparator) + dist.Name()
			manifestUri := distPath + string(os.PathSeparator) + "manifest.json"
			if _, fileNotExistErr := os.Stat(manifestUri); errors.Is(fileNotExistErr, os.ErrNotExist) {
				continue
            }

			data, err := os.ReadFile(manifestUri)
			if err != nil {
				return &DistributionManager{}, err
			}

            config := Config{}
			if err := json.Unmarshal(data, &config); err != nil {
				return &DistributionManager{}, err
			}

            distributions[config.Name] = Distribution{
                Config: config,
                Path: distPath,
            }
		}
	}

    return &DistributionManager{
        Distributions: distributions,
    }, nil
}

// Gets a distribution from the collection of distributions
func (mngr *DistributionManager) FindDistV2(name string) (Distribution, error) {
    dist, ok := mngr.Distributions[name]
    if !ok {
        return dist, errors.New("Could not find distribution")
    }
    return dist, nil
}

// Should loop through all the distributions for testing purposes
func (mngr *DistributionManager) Debug() {
    for key, val := range mngr.Distributions {
        fmt.Printf("%s: %s\n", key, val)
    }
}

// TODO: After all the functions above have been tested and the implementation has been fixed,
// we can remove this function.
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
				distribution.Config = config
				distribution.Path = distPath
				return distribution, nil
			}
		}
	}

	return distribution, fmt.Errorf("could not find matching configuration")
}

