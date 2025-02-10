package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

var (
	waypointDirs = []string{".git", "cmd", "feature", "foundation", "infra", "modmake"}
	importDirs   = []string{"cmd/yourapp", "feature", "foundation"}
	pathDirs     = []string{"modmake", "cmd/yourapp", "README.md", "docker-compose.yaml"}
	explicitAdd  = []string{"cmd/yourapp/templates/util.go"}
)

func main() {
	gitExec, goExec := checkReadiness()
	var (
		scanner = bufio.NewScanner(os.Stdin)
	)
	fmt.Print("Enter new Go module name: ")
	scanner.Scan()
	newModuleName, err := validateModuleName(scanner.Text())
	if err != nil {
		fmt.Println("Invalid module name:", err)
		os.Exit(1)
	}
	fmt.Println("Using", newModuleName, "as module name")
	fmt.Print("Enter new app name: ")
	scanner.Scan()
	newAppName, err := validateAppName(scanner.Text())
	if err != nil {
		fmt.Println("Invalid app name:", err)
		os.Exit(1)
	}
	fmt.Println("Using", newAppName, "as app name")
	fmt.Println("Applying customization changes to repo, please wait...")
	generate := exec.Command(goExec, "run", "./modmake", "generate")
	generate.Stdout = os.Stdout
	generate.Stderr = os.Stderr
	if err := generate.Run(); err != nil {
		fmt.Println("Failed to run initial generation")
		os.Exit(1)
	}
	if err := doTransform(newModuleName, newAppName); err != nil {
		fmt.Println("Error running customization logic:", err)
		fmt.Println("Rolling back changes, please wait...")
		gitRollback(gitExec)
		os.Exit(1)
	}
	fmt.Println("Customizations complete, validating Go code...")
	vet := exec.Command(goExec, "vet", "./...")
	vet.Stdout = os.Stdout
	vet.Stderr = os.Stdout
	if err := vet.Run(); err != nil {
		fmt.Println("Detected validation errors, rolling back changes...")
		gitRollback(gitExec)
		os.Exit(1)
	}
	//fmt.Println("Re-initializing git repository...")
	//if err := doFinalize(gitExec, "yourapp", newAppName); err != nil {
	//	fmt.Println("Failed to re-initialize git repo:", err)
	//	fmt.Println("This will need to be done manually")
	//}
	fmt.Println("Customization complete. Enjoy!")
}

func doTransform(newModuleName, newAppName string) error {
	for _, dir := range importDirs {
		if err := replaceDirImports(dir, "yourapp", newModuleName); err != nil {
			return err
		}
	}
	if err := changeModuleName(newModuleName); err != nil {
		return err
	}

	return nil
}

func doFinalize(gitExec, oldAppName, newAppName string) error {
	if err := os.RemoveAll(".git"); err != nil {
		return fmt.Errorf("unable to remove .git directory: %w", err)
	}
	if err := exec.Command(gitExec, "init").Run(); err != nil {
		return fmt.Errorf("unable to initialize fresh git repository: %w", err)
	}
	_ = exec.Command(gitExec, "add", ".").Run()
	for _, explicit := range explicitAdd {
		explicit = strings.ReplaceAll(explicit, oldAppName, newAppName)
		_ = exec.Command(gitExec, "add", explicit).Run()
	}
	_ = exec.Command(gitExec, "commit", "-m", "Customized template from github.com/saylorsolutions/goth-stack").Run()
	return nil
}

func checkReadiness() (gitExecutable string, goExecutable string) {
	// Make sure that we're in the root dir
	for _, dir := range waypointDirs {
		fi, err := os.Stat(dir)
		if err != nil || !fi.IsDir() {
			fmt.Println("Must be run from the root of the repository: missing directory", dir)
			os.Exit(1)
		}
	}

	git, err := exec.LookPath("git")
	if err != nil {
		fmt.Println("git is required to finalize/rollback setup")
		os.Exit(1)
	}

	goExec, err := exec.LookPath("go")
	if err != nil {
		fmt.Println("go is required to validate setup")
		os.Exit(1)
	}
	return git, goExec
}

func validateModuleName(modName string) (string, error) {
	modName = strings.TrimSpace(modName)
	if len(modName) == 0 {
		return "", errors.New("module name is empty")
	}
	return modName, nil
}

func validateAppName(appName string) (string, error) {
	appName = strings.TrimSpace(appName)
	if len(appName) == 0 {
		return "", errors.New("app name is empty")
	}
	var (
		nameRunes    = []rune(appName)
		numWritten   int
		writePos     int
		writtenRunes = make([]rune, len(nameRunes))
	)
	for _, r := range nameRunes {
		if numWritten > 0 && (unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			writtenRunes[writePos] = r
			writePos++
			numWritten++
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			writtenRunes[writePos] = r
			writePos++
			numWritten++
		}
	}
	if numWritten == 0 {
		return "", fmt.Errorf("no valid characters in name '%s'", appName)
	}
	appName = string(writtenRunes[:numWritten])
	return appName, nil
}
