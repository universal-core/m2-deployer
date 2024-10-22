package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"universal-core.com/m2-deployer/reader"
)

func buildCommandFiles(entry reader.SingleConfig) ([]OutFile, error) {
	out := make([]OutFile, 0)
	// fmt.Printf("building: %s with values %+v", outputPath, entry)
	commandDir := filepath.Join(CONFIGURATIONS_PATH, "command")
	files, err := os.ReadDir(commandDir)
	if err != nil {
		return nil, fmt.Errorf("Error reading command directory: %s\n", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			file, err := buildCommandFile(file, commandDir, entry)
			if err != nil {
				return nil, err
			}

			out = append(out, file)
		}
	}
	return out, nil
}

func buildCommandFile(file fs.DirEntry, commandDir string, entry reader.SingleConfig) (OutFile, error) {
	// Read template file
	f := file.Name()
	p := entry.GetPath()
	filePath := filepath.Join(commandDir, file.Name())
	tmplData, err := os.ReadFile(filePath)
	if err != nil {
		return OutFile{}, fmt.Errorf("Error reading command file (%s): %s\n", filePath, err)
	}

	tmpl, err := template.New(file.Name()).Parse(string(tmplData))
	if err != nil {
		return OutFile{}, fmt.Errorf("Error parsing template: %s\n", err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, entry.Values)
	if err != nil {
		return OutFile{}, fmt.Errorf("Error executing template: %s\n", err)
	}

	// perform renaming
	outName := file.Name()
	if rename, ok := entry.Renames[outName]; ok {
		outName = rename
	}
	// build the output file
	outFilePath := filepath.Join(entry.GetPath(), outName)

	fmt.Printf("> Core template [%s] for core [%s] processed succesfully\n", f, p)
	return OutFile{
		Data: result.Bytes(),
		Path: outFilePath,
	}, nil
}
