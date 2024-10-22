package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Custom function to filter out strings not containing a certain word.
func excludeStrings(slice []string, word string) []string {
	var result []string
	for _, str := range slice {
		if !strings.Contains(str, word) {
			result = append(result, str)
		}
	}
	return result
}

// Custom function to filter strings containing a certain word.
func filterStrings(slice []string, word string) []string {
	var result []string
	for _, str := range slice {
		if strings.Contains(str, word) {
			result = append(result, str)
		}
	}
	return result
}

func processTemplates(fileDir string, outPath string, args map[string]interface{}) ([]OutFile, error) {
	out := make([]OutFile, 0)

	// Read files from fileDir
	files, err := os.ReadDir(fileDir)
	if err != nil {
		return nil, fmt.Errorf("reading directory [%s]: %s", fileDir, err.Error())
	}

	// Process each file as a template
	for _, file := range files {
		if !file.IsDir() {
			// fmt.Printf("processing %s\n", file.Name())
			filePath := filepath.Join(fileDir, file.Name())
			outFile, err := executeTemplate(filePath, outPath, args, file.Name())
			if err != nil {
				return nil, err
			}

			out = append(out, outFile)
		}
	}

	return out, nil
}

func executeTemplate(filePath, outPath string, vars map[string]interface{}, fileName string) (OutFile, error) {
	// Read template file
	tmplData, err := os.ReadFile(filePath)
	if err != nil {
		return OutFile{}, fmt.Errorf("reading template file: %s", err.Error())
	}

	// Parse and execute template
	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"excludeStrings": excludeStrings,
		"filterStrings":  filterStrings,
	}).Parse(string(tmplData))
	if err != nil {
		return OutFile{}, fmt.Errorf("parsing template: %s", err.Error())
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, vars); err != nil {
		return OutFile{}, fmt.Errorf("executing template: %s", err.Error())
	}

	fmt.Printf("> Template [%s] processed succesfully\n", fileName)
	return OutFile{
		Data: output.Bytes(),
		Path: filepath.Join(outPath, fileName),
	}, nil
}
