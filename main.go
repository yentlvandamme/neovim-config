package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
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

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "load", "use":
		configName := os.Args[2]
		var err error

		dist, err := getDist(configName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(dist.config.Name)

	case "clean":
		_, err := getConfigPath()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		// RemoveConfig(configPath) // TODO: Don't execute this yet. Write a passing unit test first
	}
}

func getDist(name string) (Distribution, error) {
	var err error
	distribution := Distribution{}
	config := Config{}

	currentDir, err := os.Getwd()
	if err != nil {
		return distribution, err
	}

	distsPath := currentDir + "/distributions"
	dists, err := os.ReadDir(distsPath)
	if err != nil {
		return distribution, err
	}

	for _, dist := range dists {
		if dist.IsDir() {
			distPath := distsPath + "/" + dist.Name()
			manifestUri := distPath + "/manifest.json"
			if _, err := os.Stat(manifestUri); errors.Is(err, os.ErrNotExist) {
				continue
			}

			data, err := os.ReadFile(manifestUri)
			if err != nil {
				return distribution, err
			}

			if err := json.Unmarshal(data, &config); err != nil {
				return distribution, err
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

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "Windows" {
		return filepath.Join(homeDir, "/AppData/Local/nvim"), nil
	}

	return filepath.Join(homeDir, "/.config/nvim"), nil
}

func RemoveConfig(rootPath string) {
	filepath.WalkDir(rootPath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if rootPath != path {
			if dir.IsDir() {
				os.RemoveAll(path)
			} else {
				os.Remove(path)
			}
		}

		return nil
	})
}

func copyConfig(target string, src string) error {
	return nil
}
