package repo

import (
	"io/ioutil"
	"strings"

	builder "github.com/autamus/builder/repo"
	"gopkg.in/yaml.v3"
)

type Instruction struct {
	BuildSI BuildData     `yaml:"buildsi"`
	Spack   builder.Spack `yaml:"spack"`
}

type BuildData struct {
	Release  int                `yaml:"release"`
	Versions map[string]Version `yaml:"versions"`
}

type Version struct {
	VariantOnly bool     `yaml:"variant_only"`
	Variants    []string `yaml:"variants"`
}

// GetChangedInstructions returns a list of the changed instructions
// from a list of changed filepaths.
func GetChangedInstructions(InstructionsPath string, filepaths []string) (result []Instruction, err error) {

	for _, path := range filepaths {
		if strings.Contains(path, InstructionsPath) && strings.HasSuffix(path, ".yaml") {
			instruction := Instruction{}
			in, err := ioutil.ReadFile(path)
			if err != nil {
				return result, err
			}
			err = yaml.Unmarshal(in, &instruction)
			if err != nil {
				return result, err
			}
			result = append(result, instruction)
		}
	}
	return result, nil
}
