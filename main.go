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

		dist, distErr := getDist(configName)
		if distErr != nil {
			fmt.Fprintln(os.Stderr, distErr)
			return
		}

		fmt.Println(dist.config.Name)

	case "clean":

		configPath, configPathErr := getConfigPath()
		if configPathErr != nil {
			fmt.Fprintln(os.Stderr, configPathErr)
			return
		}

        removeErr := RemoveConfig(configPath)
		if removeErr != nil {
			fmt.Fprintln(os.Stderr, removeErr)
            return
		}
	}
}

func getDist(name string) (Distribution, error) {
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
			distPath := distsPath + "/" + dist.Name()
			manifestUri := distPath + "/manifest.json"
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

func getConfigPath() (string, error) {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		return "", homeDirErr
	}

	if runtime.GOOS == "Windows" {
		return filepath.Join(homeDir, "/AppData/Local/nvim"), nil
	}

	return filepath.Join(homeDir, "/.config/nvim"), nil
}

func RemoveConfig(rootPath string) error {
	err := filepath.WalkDir(rootPath, func(path string, dir fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if rootPath != path {
			if dir.IsDir() {
                fmt.Println(path)
				os.RemoveAll(path)
                fmt.Println("after")
			} else {
				os.Remove(path)
			}
		}

		return nil
	})

	return err
}

func copyConfig(target string, src string) error {
	return nil
}
