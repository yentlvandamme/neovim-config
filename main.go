package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	configs "github.com/yentlvandamme/neovim-config/libs"
)

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

        dist, distErr := configs.FindDist(configName)
		if distErr != nil {
			fmt.Fprintln(os.Stderr, distErr)
			return
		}

        fmt.Printf("Found distributions %s\n", dist.Config.Name)

        removeErr := RemoveConfig(configPath)
        if removeErr != nil {
            fmt.Fprintln(os.Stderr, removeErr)
        }

        copyConfigErr := CopyConfig(configPath, dist.Path)
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
