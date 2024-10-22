package reader

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"gopkg.in/yaml.v2"
)

// Config represents the structure of the YAML file.
type SingleConfig struct {
	Renames map[string]string      `yaml:"renames"`
	Values  map[string]interface{} `yaml:"values"`
}

func (sg SingleConfig) GetPath() string {
	return sg.Values["path"].(string)
}

type AdditionalConfigs map[string]AdditionalConfig

type AdditionalConfig struct {
	OutPath string                 `yaml:"out_path"`
	Values  map[string]interface{} `yaml:"values"`
}

type Config struct {
	Cores             map[string]SingleConfig `yaml:"cores"`
	Share             map[string]interface{}  `yaml:"share"`
	Globals           map[string]interface{}  `yaml:"globals"`
	Metadata          ConfigMetadata          `yaml:"metadata"`
	AdditionalConfigs AdditionalConfigs       `yaml:"additional_folders"`
}

type ConfigMetadata struct {
	CreateDirs    []string `yaml:"create_dirs"`
	Touch         []string `yaml:"touch"`
	InstallScript string   `yaml:"install_script"`
}

type ConfigReader struct {
	files          []loadedFile
	configurations []SingleConfig
	base           Config
}

type loadedFile struct {
	path   string
	data   []byte
	parsed map[string]interface{}
}

// Custom template functions
func add(x, y int) int {
	return x + y
}

// Custom template functions
func sub(x, y int) int {
	return x - y
}

// Helper functions for the template
func seq(start, end int) []int {
	var s []int
	for i := start; i < end; i++ {
		s = append(s, i)
	}
	return s
}

func mul(x, y int) int {
	return x * y
}

func NewConfigReader(files []string) *ConfigReader {
	dataFiles := []loadedFile{}
	for _, file := range files {
		// Read and parse each YAML file
		data, err := os.ReadFile(file)
		if err != nil {
			erro := ""
			if os.IsNotExist(err) {
				erro = fmt.Sprintf("file [%s] do not exist", file)
			} else {
				erro = fmt.Sprintf("error loading file [%s]: %s", file, err.Error())
			}
			panic(erro)
		}
		dataFiles = append(dataFiles, loadedFile{
			path: file,
			data: data,
		})
	}
	reader := &ConfigReader{
		files:          dataFiles,
		configurations: make([]SingleConfig, 0),
	}
	return reader
}

func (cr *ConfigReader) Preprocess(preprocessorVars map[string]interface{}) error {
	for i, f := range cr.files {
		// Parse and execute template
		tmpl, err := template.New("config").Funcs(template.FuncMap{
			"sub": sub,
			"add": add,
			"mul": mul,
			"seq": seq,
		}).Parse(string(f.data))
		if err != nil {
			return fmt.Errorf("> Preprocess | could not parse configuration file [%s]: %s", f.path, err.Error())
		}

		var output bytes.Buffer
		if err := tmpl.Execute(&output, preprocessorVars); err != nil {
			return fmt.Errorf("> Preprocess | could not execute configuration file [%s]: %s", f.path, err.Error())
		}

		var parsed map[interface{}]interface{}
		err = yaml.Unmarshal(output.Bytes(), &parsed)
		if err != nil {
			return fmt.Errorf("> Preprocess | could not deserialize configuration file [%s]: %s", f.path, err.Error())
		}

		// Convert and expand dotted keys
		converted := convertMapInterface(parsed)
		cr.files[i].parsed = expandDots(converted)
	}
	return nil
}

func (cr *ConfigReader) Generate() error {
	if err := cr.mergeYAMLFiles(&cr.base); err != nil {
		return err
	}

	cr.base.Globals = convertMapInterface(cr.base.Globals)
	for _, coreConfig := range cr.base.Cores {
		if err := injectGlobals(&coreConfig, cr.base.Globals); err != nil {
			log.Fatal(err)
		}
		cr.configurations = append(cr.configurations, coreConfig)
	}
	return nil
}

func (cr *ConfigReader) GetConfigurations() []SingleConfig {
	return cr.configurations
}

func (cr *ConfigReader) GetMetadata() ConfigMetadata {
	return cr.base.Metadata
}

func (cr *ConfigReader) GetGlobals() map[string]interface{} {
	return cr.base.Globals
}

func (cr *ConfigReader) GetAdditionalConfigs() AdditionalConfigs {
	return cr.base.AdditionalConfigs
}

func (cr *ConfigReader) GetTouch() []string {
	return cr.GetMetadata().Touch
}

func (cr *ConfigReader) GetCreateDirs() []string {
	return cr.GetMetadata().CreateDirs
}

func injectGlobals(coreConfig *SingleConfig, globals map[string]interface{}) error {
	coreConfig.Values = convertMapInterface(coreConfig.Values)
	if err := checkTypeConflicts(coreConfig.Values, globals, ""); err != nil {
		return fmt.Errorf("> Generate | Type missmatch while injecting globals in [%s]: %s", coreConfig.Values["hostname"].(string), err.Error())
	}

	// Merge the current file's data into the mergedData
	if err := mergo.Merge(&coreConfig.Values, globals, mergo.WithOverride); err != nil {
		return fmt.Errorf("> Generate | Could not merge globals [%s]: %s", coreConfig.Values["hostname"].(string), err.Error())
	}
	return nil
}

// mergeYAMLFiles reads multiple YAML files, merges their content, and unmarshals into the provided struct
func (cr *ConfigReader) mergeYAMLFiles(result interface{}) error {
	var mergedData map[string]interface{}

	for _, file := range cr.files {
		if err := checkTypeConflicts(mergedData, file.parsed, ""); err != nil {
			return fmt.Errorf("> Generate | Type missmatch in file [%s]: %s", file.path, err.Error())
		}

		// Merge the current file's data into the mergedData
		if err := mergo.Merge(&mergedData, file.parsed, mergo.WithOverride); err != nil {
			return fmt.Errorf("> Generate | Could not merge file [%s]: %s", file.path, err.Error())
		}
	}

	// Unmarshal the merged data into the result struct
	mergedYAML, err := yaml.Marshal(mergedData)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(mergedYAML, result); err != nil {
		return err
	}

	return nil
}

// flattenKeys takes a nested map and flattens it into a single map with dotted keys
func flattenKeys(prefix string, m map[string]interface{}) map[string]interface{} {
	flattened := make(map[string]interface{})
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		switch nested := v.(type) {
		case map[string]interface{}:
			for nestedKey, nestedValue := range flattenKeys(fullKey, nested) {
				flattened[nestedKey] = nestedValue
			}

		default:
			flattened[fullKey] = v
		}
	}
	return flattened
}

// convertMapInterface converts map[interface{}]interface{} to map[string]interface{}
func convertMapInterface[T comparable](input map[T]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for k, v := range input {
		strKey := fmt.Sprintf("%v", k)
		switch val := v.(type) {
		case map[interface{}]interface{}:
			converted[strKey] = convertMapInterface(val)
		default:
			converted[strKey] = val
		}
	}
	return converted
}

// expandDots expands dotted keys into nested maps
func expandDots(m map[string]interface{}) map[string]interface{} {
	expanded := make(map[string]interface{})

	for k, v := range m {
		keys := strings.Split(k, ".")
		currentMap := expanded

		for i, key := range keys {
			if i == len(keys)-1 {
				currentMap[key] = v
			} else {
				if _, exists := currentMap[key]; !exists {
					currentMap[key] = make(map[string]interface{})
				}
				currentMap = currentMap[key].(map[string]interface{})
			}
		}
	}

	// Recursively expand nested maps
	for k, v := range expanded {
		if nestedMap, ok := v.(map[string]interface{}); ok {
			expanded[k] = expandDots(nestedMap)
		}
	}

	return expanded
}

func checkTypeConflicts(map1, map2 map[string]interface{}, parentKey string) error {
	for key, val1 := range map1 {
		fullKey := key
		if parentKey != "" {
			fullKey = parentKey + "." + key
		}

		if val2, exists := map2[key]; exists {
			type1 := fmt.Sprintf("%T", val1)
			type2 := fmt.Sprintf("%T", val2)

			if type1 != type2 {
				return fmt.Errorf("type conflict for key %s: %s vs %s", fullKey, type1, type2)
			}

			// Recursively check nested maps
			if nestedMap1, ok1 := val1.(map[string]interface{}); ok1 {
				if nestedMap2, ok2 := val2.(map[string]interface{}); ok2 {
					if err := checkTypeConflicts(nestedMap1, nestedMap2, fullKey); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
