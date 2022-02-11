package php

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/digest"
	"os"

	"github.com/nais/salsa/pkg/build"
)

const composerLockFileName = "composer.lock"

type Composer struct {
	BuildFilePatterns []string
}

func (c Composer) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, composerLockFileName))
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps, err := ComposerDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("scan: %v\n", err)
	}
	return &build.ArtifactDependencies{
		Cmd:         composerLockFileName,
		RuntimeDeps: deps,
	}, nil
}

func NewComposer() build.BuildTool {
	return &Composer{
		BuildFilePatterns: []string{composerLockFileName},
	}
}

func (c Composer) BuildFiles() []string {
	return c.BuildFilePatterns
}

type dist struct {
	Shasum string `json:"shasum"`
}

type dep struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Dist    dist   `json:"dist"`
}

type composerLock struct {
	Dependencies []dep `json:"packages"`
}

func ComposerDeps(composerLockJsonContents string) ([]build.Dependency, error) {
	var lock composerLock
	err := json.Unmarshal([]byte(composerLockJsonContents), &lock)
	if err != nil {
		return nil, fmt.Errorf("unable to parse composer.lock: %v", err)
	}

	return transform(lock.Dependencies), nil
}

func transform(dependencies []dep) []build.Dependency {
	deps := make([]build.Dependency, 0)
	for _, d := range dependencies {
		deps = append(deps, build.Dependency{
			Coordinates: d.Name,
			Version:     d.Version,
			CheckSum: build.CheckSum{
				Algorithm: digest.AlgorithmSHA1,
				Digest:    d.Dist.Shasum,
			},
		})
	}
	return deps
}