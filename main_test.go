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

func createEmptyConfigFolder(rootDirPath string) string {
    configPath := rootDirPath+".config/"
    os.MkdirAll(configPath, 0755)
    return configPath
}

func createDistribution(rootDirPath string) string {
    distPath := rootDirPath+"/distributions/"
    os.MkdirAll(distPath, 0755)
    os.MkdirAll(distPath+"src/", 0755)
    os.MkdirAll(distPath+"src/themes/", 0755)

    mainFile, err := os.Create(distPath+"src/main.lua")
	if err != nil {
		fmt.Println("Could not create main file")
		panic(err)
	}
	defer mainFile.Close()

    colorsFile, err := os.Create(distPath+"src/themes/colors.lua")
	if err != nil {
		fmt.Println("Could not create colors file")
		panic(err)
	}
	defer colorsFile.Close()

    return distPath
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

func TestCopyConfig(t *testing.T) {
    // Arrange
    rootDirPath := t.TempDir()
    configPath := createEmptyConfigFolder(rootDirPath)
    distributionPath := createDistribution(rootDirPath)

    // Act
    CopyConfig(configPath, distributionPath)

    // Assert
    pathsToCheck := []string{configPath+"src", configPath+"src/themes", configPath+"src/main.lua", configPath+"src/themes/colors.lua"}
    for _, path := range pathsToCheck {
        if _, err := os.Stat(path); os.IsNotExist(err) {
            t.Fatalf("The directory does not contain the path: %s", path)
        }
    }

    // Clean-up
    os.RemoveAll(rootDirPath)
}
