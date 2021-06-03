package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	buildconfig "github.com/autamus/buildconfig/repo"
	builder "github.com/autamus/builder/repo"
	"github.com/autamus/builder/spack"
	"github.com/buildsi/build-sandbox/verBuilder/config"
)

func main() {
	// Set initial values for Repository
	path := config.Global.Repository.Path
	pubKeyURL := config.Global.Packages.PublicKeyURL
	instructionsPath := filepath.Join(path, config.Global.Instructions.Path)
	defaultEnvPath := filepath.Join(path, config.Global.Instructions.DefaultEnVPath)
	// Declare instruction values
	currentBuild := config.Global.Instructions.Current
	currentDockerfile := ""

	mainBranch := config.Global.Repository.DefaultBranch
	// Get the name of the current branch.
	currentBranch, err := buildconfig.GetBranchName(path)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the current run is a PR
	prVal, prExists := os.LookupEnv("GITHUB_EVENT_NAME")
	isPR := prExists && prVal == "pull_request"

	// Get a list of all of the changed files in the commit.
	filepaths, err := buildconfig.GetChangedFiles(path, currentBranch, mainBranch)
	if err != nil {
		log.Fatal(err)
	}

	instructPath := ""
	for _, path := range filepaths {
		if strings.Contains(path, instructionsPath) && strings.HasSuffix(path, ".yaml") {
			instructPath = path
		}
	}
	fmt.Println(instructPath)

	// If the container is a spack environment, find the main spec.
	spackEnv, err := builder.ParseSpackEnv(defaultEnvPath, instructPath)
	if err != nil {
		log.Fatal(err)
	}
	spackEnv.Spack.Specs = []string{currentBuild}

	// Containerize SpackEnv to Dockerfile
	currentDockerfile, err = spack.Containerize(spackEnv, isPR, pubKeyURL)
	if err != nil {
		log.Fatal(err)
	}
	// Write the Dockerfile out to Disk
	f, err := os.Create(filepath.Join(path, "Dockerfile"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(currentDockerfile)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}
