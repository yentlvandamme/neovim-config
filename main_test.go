package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func initConfigFileStructure() vfs.FileSystem {
	fs := mapfs.New(map[string]string{
		".config/":                      "",
		".config/dist/":                 "",
		".config/dist/main.lua":         "main file",
		".config/dist/utils/":           "",
		".config/dist/utils/format.lua": "utils file",
	})

	return fs
}

func createConfigFileStructure(t *testing.T) string {
	rootDirPath := t.TempDir()
	os.MkdirAll(rootDirPath+"/.config/", 0755)
	os.MkdirAll(rootDirPath+"/.config/dist/", 0755)
	mainFile, err := os.Create(rootDirPath + "/.config/main.lua")
	if err != nil {
		fmt.Println("Could not create main file")
		panic(err)
	}
	defer mainFile.Close()

	os.MkdirAll(rootDirPath+"/.config/utils/", 0755)
	utilsFile, err := os.Create(rootDirPath + "/.config/utils/format.lua")
	if err != nil {
		fmt.Println("Could not create file")
		panic(err)
	}
	defer utilsFile.Close()

	return rootDirPath
}

func TestRemoveCurrentConfig(t *testing.T) {
	rootDirPath := createConfigFileStructure(t)
	filepath.WalkDir(rootDirPath, func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		if err != nil {
			return err
		}
		return nil
	})
	os.RemoveAll(rootDirPath)
}
