package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

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

    os.MkdirAll(rootDirPath+"/.config/utils/theme/", 0755)
    themeFile, err := os.Create(rootDirPath+"/.config/utils/theme/colors.lua")
    if err != nil {
        fmt.Println("Could not create file")
        panic(err)
    }
    defer themeFile.Close()

	return rootDirPath
}

func TestRemoveCurrentConfig(t *testing.T) {
	// Arrange
	rootDirPath := createConfigFileStructure(t)

	// Act
	RemoveConfig(rootDirPath)

	// Assert
	// => The root file should still exist
	rootDir, err := os.Stat(rootDirPath)
	if err != nil {
		t.Fatalf("Root directory %s does not exist", rootDirPath)
	}
	if !rootDir.IsDir() {
		t.Fatalf("Root path %s is not a directory", rootDirPath)
	}

	// => Everything in the root file should be gone
	filepath.WalkDir(rootDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatalf(err.Error())
		}

		if rootDirPath != path {
			t.Fatalf("Directory %s is not empty", rootDirPath)
		}

		return err
	})

	// Clean-up
	os.RemoveAll(rootDirPath)
}
