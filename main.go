package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
)

type Config struct {
	Name        string
	Version     string
	Description string
}

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "load", "use":
		configName := os.Args[2]
		var err error

		config, err := getConfig(configName)
		if err != nil {
            fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(config.Name)
	}
}

func getConfig(name string) (Config, error) {
	var err error
	config := Config{}

	currentDir, err := os.Getwd()
	if err != nil {
		return config, err
	}

	distsPath := currentDir + "/distributions"
	dists, err := os.ReadDir(distsPath)
	if err != nil {
		return config, err
	}

	for _, dist := range dists {
		if dist.IsDir() {
			manifestUri := distsPath + "/" + dist.Name() + "/manifest.json"
			if _, err := os.Stat(manifestUri); errors.Is(err, os.ErrNotExist) {
				continue
			}

			data, err := os.ReadFile(manifestUri)
			if err != nil {
				return config, err
			}

			if err := json.Unmarshal(data, &config); err != nil {
				return config, err
			}

			if config.Name == name {
				return config, nil
			}
		}
	}

	return config, fmt.Errorf("Could not find matching configuration")
}

func getConfigPath() string {
	if runtime.GOOS == "Windows" {
		return "~/AppData/Local/nvim"
	}
	return "~/.config/nvim"
}

func copyConfig(target string, src string) error {
	return nil
}
