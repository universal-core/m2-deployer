package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/pflag"
	"universal-core.com/m2-deployer/reader"
)

const (
	DEPLOY_PATH         = "deploy"
	CONFIGURATIONS_PATH = "configurations"
)

var OUT_PATH string

type OutFile struct {
	Data []byte
	Path string
}

// createPath creates necessary directories up to the specified file path
func createPath(file string) string {
	// Extract the directory portion of the file path
	dir := filepath.Dir(file)

	// Create necessary directories
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("Error creating directories for file (%s): %s\n", dir, err))
	}

	return dir
}

func main() {
	// Parse command-line argument
	var files []string
	pflag.StringSliceVar(&files, "config", []string{}, "Paths to the configuration files")
	out := pflag.String("out", DEPLOY_PATH, "Defines the output path for the deployment")
	channels := pflag.Int("channels", 0, "Defines amount of channels we want to generate")
	pflag.Parse()

	// Check if channels flag is set
	if *channels == 0 {
		fmt.Println("Error: --channels flag is mandatory")
		pflag.Usage()
		os.Exit(1)
	}

	// Split the out path into a list
	outPaths := filepath.SplitList(*out)

	// Reassemble the path using the correct separator for the OS
	outputFolder := filepath.Join(outPaths...)

	prepVars := map[string]interface{}{
		"channels": *channels,
	}
	reader := reader.NewConfigReader(files)

	if err := reader.Preprocess(prepVars); err != nil {
		log.Fatalf(err.Error())
	}

	if err := reader.Generate(); err != nil {
		log.Fatalf(err.Error())
	}

	outFiles := make([]OutFile, 0)

	globals := reader.GetGlobals()
	globals["paths"] = make([]string, 0)
	globals["__out__cores"] = []interface{}{}

	// Process each core entry
	for _, entry := range reader.GetConfigurations() {
		globals["paths"] = append(globals["paths"].([]string), entry.GetPath())
		globals["__out__cores"] = append(globals["__out__cores"].([]interface{}), entry.Values)
		files, err := buildCommandFiles(entry)
		if err != nil {
			log.Fatalf(err.Error())
		}

		outFiles = append(outFiles, files...)
	}

	// process root
	res, err := processTemplates(filepath.Join(CONFIGURATIONS_PATH, "root_command"), "", globals)
	if err != nil {
		log.Fatalf(err.Error())
	}
	outFiles = append(outFiles, res...)

	for folder, entry := range reader.GetAdditionalConfigs() {
		entry.Values["globals"] = globals
		res, err := processTemplates(filepath.Join(CONFIGURATIONS_PATH, folder), entry.OutPath, entry.Values)
		if err != nil {
			log.Fatalf(err.Error())
		}
		outFiles = append(outFiles, res...)
	}

	// flush
	for _, f := range outFiles {
		// Write to file in deploy directory
		filePath := filepath.Join(outputFolder, f.Path)
		createPath(filePath)
		if err := os.WriteFile(filePath, f.Data, 0644); err != nil {
			log.Fatalf("!> Error writing output file: %w", err)
		}
	}

	// postcreate
	for _, fi := range reader.GetTouch() {
		filePath := filepath.Join(outputFolder, fi)
		//create file if not exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("> Touching [%s]\n", filePath)
			createPath(filePath)
			if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
				log.Fatalf("!> Error touching file [%s]: %w", fi, err)
			}
		}
	}
	for _, d := range reader.GetCreateDirs() {
		dirPath := filepath.Join(outputFolder, d)
		//create dir if not exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Printf("> Touching [%s]\n", dirPath)
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				log.Fatalf("!> Error creating dir [%s]: %w", d, err)
			}
		}
	}

	// lastly, launch install
	if len(reader.GetMetadata().InstallScript) != 0 && runtime.GOOS == "linux" {
		// First command: Change permissions of .sh files
		chmodCmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && find . -type f -name \"*.sh\" -exec chmod +x {} \\;", outputFolder))
		var chmodOutBuf, chmodErrBuf bytes.Buffer
		chmodCmd.Stdout = &chmodOutBuf
		chmodCmd.Stderr = &chmodErrBuf

		err := chmodCmd.Run()
		chmodStderr := chmodErrBuf.String()

		if err != nil {
			fmt.Printf("chmod command execution failed: %v\n", err)
			fmt.Printf("stderr: %s\n", chmodStderr)
			return
		} else {
			fmt.Println("chmod command executed successfully")
		}

		// Second command: Execute the install script
		installCmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && ./%s;", outputFolder, reader.GetMetadata().InstallScript))
		var installOutBuf, installErrBuf bytes.Buffer
		installCmd.Stdout = &installOutBuf
		installCmd.Stderr = &installErrBuf

		err = installCmd.Run()
		installStdout := installOutBuf.String()
		installStderr := installErrBuf.String()

		if err != nil {
			fmt.Printf("InstallScript execution failed: %v\n", err)
			fmt.Printf("stderr: %s\n", installStderr)
		} else {
			fmt.Println("InstallScript executed successfully")
			fmt.Printf("stdout: %s\n", installStdout)
			fmt.Printf("stderr: %s\n", installStderr)
		}
	}
}
