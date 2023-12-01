package cmd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/CeoFred/nturu/cmd/utils"
)

//go:embed templates/*
var embededTemplates embed.FS

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Creates a new microservice using the default boilerplate.",
	Long:  `This generates a new microservice using the default boilerplate.`,
	Run: func(cmd *cobra.Command, args []string) {
		// ask for app name
		var ModulePath string
		var AppName string

		fmt.Println("\033[1;31m What is your application name?\033[0m")
		fmt.Scanln(&AppName)
		clearScreen()

		fmt.Println("\033[1;31m What is your preferred module path?\033[0m")
		fmt.Scanln(&ModulePath)
		clearScreen()

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Current working directory:", currentDir)

		err = createFolder(currentDir, AppName)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		sourceFolder := filepath.Join(currentDir, "templates")
		destinationFolder := filepath.Join(currentDir, AppName)

		err = copyTemplates(sourceFolder, destinationFolder)
		if err != nil {
			fmt.Println("Error:", err)
			deleteFolder(destinationFolder)
			return
		}

		fmt.Println("----------------------------------------------------------------")

		// err = utils.CreateGoModFile(destinationFolder, ModulePath)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	deleteFolder(destinationFolder)
		// 	return
		// }

		zippedTemplate := filepath.Join(destinationFolder, "default.zip")

		err = utils.UnzipFile(zippedTemplate)
		if err != nil {
			fmt.Println("Error:", err)
			deleteFolder(destinationFolder)
			return
		}

		fmt.Println("Extracting template..")


	
		err = os.Remove(zippedTemplate)
		if err != nil {
			fmt.Println("Error:", err)
			deleteFolder(destinationFolder)
			return
		}
		fmt.Println("\033[1;31mDone! Template generated successfully. Say Hi to @codemon_")

		
	},
}

func deleteFolder(destinationFolder string) {
	err := os.RemoveAll(destinationFolder)
	if err != nil {
		fmt.Println("Error deleting folder:", err)
		return
	}
}

func createFolder(dir, newFolderName string) error {
	// Join the current directory and the new folder name
	newFolderPath := filepath.Join(dir, newFolderName)
	err := os.Mkdir(newFolderPath, 0755) // 0755 is the default permission mode
	return err
}

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func copyTemplates(src, dst string) error {

	templates, err := embededTemplates.ReadDir("templates")
	if err != nil {
		fmt.Println("failed to read templates")
		return err
	}

	for _, entry := range templates {
		srcPath := filepath.Join("templates", entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			data, err := embededTemplates.ReadFile(srcPath)
			if err != nil {
				return err
			}

			err = os.WriteFile(destPath, data, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
