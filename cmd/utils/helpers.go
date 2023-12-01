package utils
import (
	"fmt"
	"os/exec"
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