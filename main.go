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


func main() {
	cmd := os.Args[1]

	switch cmd {
	case "load", "use":
		configName := os.Args[2]

        configPath, configPathErr := getConfigPath()
        if configPathErr != nil {
            fmt.Fprintln(os.Stderr, configPathErr)
            return
        }

        dist, distErr := FindDist(configName)
		if distErr != nil {
			fmt.Fprintln(os.Stderr, distErr)
			return
		}

        fmt.Printf("Found distributions %s\n", dist.config.Name)

        removeErr := RemoveConfig(configPath)
        if removeErr != nil {
            fmt.Fprintln(os.Stderr, removeErr)
        }

        copyConfigErr := CopyConfig(configPath, dist.path)
        if copyConfigErr != nil {
            fmt.Fprintln(os.Stderr, copyConfigErr)
        }
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

func getConfigPath() (string, error) {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		return "", homeDirErr
	}

    var configPath string
	if runtime.GOOS == "Windows" {
		configPath = filepath.Join(homeDir, "/AppData/Local/nvim")
	} else {
        configPath = filepath.Join(homeDir, "/.config/nvim")
    }

    if _,configPathErr := os.Stat(configPath); os.IsNotExist(configPathErr) {
        mkDirErr := os.MkdirAll(configPath, 0755)
        if mkDirErr != nil {
            return "", mkDirErr
        }
    }

    return configPath, nil
}

// TODO: This function isn't properly returning errors. The returned error is hard-coded as nil
func RemoveConfig(rootPath string) error {
	err := filepath.WalkDir(rootPath, func(path string, dir fs.DirEntry, walkErr error) error {
		if path != rootPath {
			if dir.IsDir() {
				os.RemoveAll(path)
			} else {
				os.Remove(path)
			}
		}

		return nil
	})

	return err
}

func CopyConfig(configPath string, distPath string) error {
    entries, err := os.ReadDir(distPath)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        currentEntryPath := distPath + string(os.PathSeparator) + entry.Name()
        entryTargetPath := configPath + string(os.PathSeparator) + entry.Name()

        if entry.IsDir() {
            os.MkdirAll(entryTargetPath, 0755)
            err := CopyConfig(entryTargetPath, currentEntryPath)
            if err != nil {
                return err
            }
        } else {
            fileContents, err := os.ReadFile(currentEntryPath)
            if err != nil {
                return err
            }
            os.WriteFile(entryTargetPath, fileContents, 0666)
        }
    }

	return nil
}
