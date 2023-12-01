package utils
import (
	"fmt"
	"os/exec"

	"archive/zip"
	"io"
	"os"
	"path/filepath"
)


func CreateGoModFile(dst string, moduleName string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dst

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running 'go mod init': %v\n%s", err, output)
	}

	fmt.Printf("go.mod file created successfully in %s with module name '%s'\n", dst, moduleName)
	return nil
}

func UnzipFile(source string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Get the parent directory of the zip file
	parentDir := filepath.Dir(source)

	for _, file := range zipReader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		// Modify the target path to the parent directory
		targetPath := filepath.Join(parentDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
			extractedFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer extractedFile.Close()

			_, err = io.Copy(extractedFile, zippedFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFolder(src, dst string) error {
	// Create the destination folder if it doesn't exist
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	// Walk through the source folder
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root folder
		if path == src {
			return nil
		}

		// Construct the destination path
		destPath := filepath.Join(dst, path[len(src):])

		// If it's a directory, create it in the destination
		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		// If it's a file, copy it to the destination
		return copyFile(path, destPath)
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}