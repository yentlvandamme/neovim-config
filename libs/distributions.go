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

func NewDistManager(distsPath string) (*DistributionManager, error) {
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

func (mngr *DistributionManager) FindDist(name string) (Distribution, error) {
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

