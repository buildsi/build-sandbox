package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	buildconfig "github.com/autamus/buildconfig/repo"
	"github.com/autamus/go-parspack"
	"github.com/buildsi/build-sandbox/verConf/config"
	"github.com/buildsi/build-sandbox/verConf/repo"
)

func main() {
	// Set initial values for Repository.
	path := config.Global.Repository.Path
	instructionsPath := config.Global.Instructions.Path
	packagesPath := config.Global.Packages.Path
	mainBranch := config.Global.Repository.DefaultBranch

	// Get the name of the current branch.
	currentBranch, err := buildconfig.GetBranchName(path)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all of the changed files in the commit.
	filepaths, err := buildconfig.GetChangedFiles(path, currentBranch, mainBranch)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all of the changed instructions.
	instructs, err := repo.GetChangedInstructions(instructionsPath, filepaths)
	if err != nil {
		log.Fatal(err)
	}

	// Check if no instructions were changed.
	if len(instructs) <= 0 {
		fmt.Println("No Instructions Changed Not Proceeding with Build.")
		return
	}

	// Check for multiple changed instructions
	if len(instructs) > 1 {
		log.Fatal(
			"More than one instruction changed. Unknown how to proceed with a build.",
		)
	}

	// Start processing current build for information.
	// Get the filepath of the build's spec.
	currentBuild := instructs[0]
	// Check if more than one spec is defined in the instruction.
	if len(currentBuild.Spack.Specs) > 1 {
		log.Fatal("More than one Spec is defined in the instruction. " +
			"Because GitHub actions current doesn't support a matrix within a matrix, " +
			"we can only support building many versions of one package at the moment.")
	}

	currentSpec := currentBuild.Spack.Specs[0]
	specPath := filepath.Join(packagesPath, currentSpec, "package.py")
	content, err := ioutil.ReadFile(specPath)
	if err != nil {
		log.Fatal(err)
	}
	// Parse the package's spec.
	pkg, err := parspack.Decode(string(content))
	if err != nil {
		log.Fatal(err)
	}

	// Construct Versions -> Variant Map.
	out := make(map[string][]string)

	// Check if all variants key is defined.
	_, isAll := currentBuild.BuildSI.Versions["all"]

	// Construct the base of all of the versions.
	for _, version := range pkg.Versions {
		if isAll && !currentBuild.BuildSI.Versions[version.Value.String()].VariantOnly {
			out[version.Value.String()] = currentBuild.BuildSI.Versions["all"].Variants
		} else {
			out[version.Value.String()] = []string{}
		}
	}

	// Apply special variants to defined versions
	for version, variants := range currentBuild.BuildSI.Versions {
		if version == "all" {
			continue
		}
		out[version] = append(out[version], variants.Variants...)
	}

	// Write output list of versions/variants to JSON list
	versionsList := []string{}
	for version, variants := range out {
		versionsList = append(versionsList, currentSpec+"@"+version+strings.Join(variants, ""))
	}
	result, err := json.Marshal(versionsList)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("::set-output name=matrix::%s\n", string(result))
}
