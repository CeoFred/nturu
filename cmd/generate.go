package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/CeoFred/nturu/utils"
)

//go:embed templates/*
var embededTemplates embed.FS
var Framework string
var Verbose bool

func init() {
	generateCmd.Flags().StringVarP(&Framework, "framework", "f", "default", "Go lang Framework to use")
	rootCmd.AddCommand(generateCmd)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "generate", "g", false, "verbose output")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Creates a new microservice using the default boilerplate.",
	Long:  `This generates a new microservice using the default boilerplate.`,
	Run: func(cmd *cobra.Command, args []string) {

		var framework string

		if len(args) > 0 {
			framework = args[0]
		}

		if framework == "" {
			framework = "default"
		}

		var availableTemplates []string

		availableTemplates = append(availableTemplates, "fiber", "default")

		if !isTemplateAvailable(framework, availableTemplates) {
			fmt.Println("Error:", "Template not available")
			os.Exit(1)
		}
		clearScreen()

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

		sourceFolder := filepath.Join(currentDir, "templates/"+framework)
		destinationFolder := filepath.Join(currentDir, AppName)

		err = copyTemplates(sourceFolder, destinationFolder, framework)
		if err != nil {
			fmt.Println("Error:", err)
			deleteFolder(destinationFolder)
			return
		}

		fmt.Println("----------------------------------------------------------------")

		zippedTemplate := filepath.Join(destinationFolder, framework+".zip")

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

		utils.ReplaceInDirectory(destinationFolder, "github.com/nturu/microservice-template", ModulePath)
		fmt.Println("\033[1;31mDone! Template generated successfully. Say Hi to @codemon_")

	},
}

func isTemplateAvailable(template string, availableTemplates []string) bool {
	for _, t := range availableTemplates {
		if t == template {
			return true
		}
	}
	return false
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

func copyTemplates(src, dst, framework string) error {
	templates, err := embededTemplates.ReadDir(filepath.Join("templates", framework))
	if err != nil {
		fmt.Printf("failed to read templates for framework %s\n", framework)
		return err
	}

	for _, entry := range templates {
		srcPath := filepath.Join("templates", framework, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		// Ensure that the entry is within the specified framework directory
		if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(dst)) {
			return errors.New("attempted to copy files outside of the specified framework directory")
		}

		if entry.IsDir() {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}

			// Recursively copy contents of subdirectories
			err = copyTemplates(srcPath, destPath, framework)
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
