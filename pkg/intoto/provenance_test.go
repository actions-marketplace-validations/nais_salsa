package intoto

import (
	"fmt"
	"github.com/nais/salsa/pkg/digest"
	"os"
	"testing"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
	"github.com/stretchr/testify/assert"
)

func TestCreateProvenanceArtifact(t *testing.T) {
	deps := toDeps()
	artDeps := toArtDeps(deps)

	provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps, nil)
	assert.Equal(t, "artifact", provenanceArtifact.Name)
	assert.Equal(t, vcs.AdHocBuildType, provenanceArtifact.BuildType)
	assert.Equal(t, deps, provenanceArtifact.Dependencies.RuntimeDeps)
	assert.True(t, time.Now().UTC().After(provenanceArtifact.BuildStartedOn))
}

func TestCreateProvenanceArtifact_withContext(t *testing.T) {
	deps := toDeps()
	artDeps := toArtDeps(deps)
	env := toVcsEnvironment()
	err := os.Setenv("GITHUB_ACTIONS", "true")
	assert.NoError(t, err)

	provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps, env)
	slsaPredicate, err := GenerateSlsaPredicate(provenanceArtifact)
	assert.NoError(t, err)

	// VCS Context
	assert.Equal(t, "https://github.com/nais/salsa/actions/runs/1234", slsaPredicate.Metadata.BuildInvocationID)
	assert.Equal(t, vcs.BuildType, slsaPredicate.BuildType)
	assert.NotEmpty(t, slsaPredicate.Invocation)
	assert.NotEmpty(t, slsaPredicate.Invocation.Parameters)
	assert.NotEmpty(t, slsaPredicate.Invocation.Environment)
	assert.Equal(t, toExpectedConfigSource(), slsaPredicate.Invocation.ConfigSource)
	assert.Equal(t, "https://github.com/nais/salsa/Attestations/GitHubHostedActions@v1", slsaPredicate.Builder.ID)
	assert.Equal(t, toExpectedMaterials(), slsaPredicate.Materials)

	// completeness
	assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Environment)
	assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Materials)
	assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Parameters)
}

func TestProvenanceArtifact_GenerateSlsaPredicate(t *testing.T) {
	deps := toDeps()
	artDeps := toArtDeps(deps)

	provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps, nil)
	slsaPredicate, err := GenerateSlsaPredicate(provenanceArtifact)
	assert.NoError(t, err)

	// Predicate
	expectedConfigSource := slsa.ConfigSource{
		URI:        "",
		Digest:     slsa.DigestSet(nil),
		EntryPoint: "",
	}

	expectedMaterials := []slsa.ProvenanceMaterial{
		{
			URI: "pkg:groupId:artifactId:v01",
			Digest: slsa.DigestSet{
				"todo": "todo",
			},
		},
	}

	assert.Equal(t, vcs.AdHocBuildType, slsaPredicate.BuildType)
	assert.Equal(t, "https://github.com/nais/salsa", slsaPredicate.Builder.ID)
	assert.Equal(t, "Some commands that made this build", slsaPredicate.BuildConfig)
	assert.Equal(t, expectedConfigSource, slsaPredicate.Invocation.ConfigSource)
	assert.Equal(t, nil, slsaPredicate.Invocation.Parameters)
	assert.Equal(t, nil, slsaPredicate.Invocation.Environment)

	// metadata
	assert.Equal(t, "", slsaPredicate.Metadata.BuildInvocationID)
	assert.True(t, time.Now().UTC().After(*slsaPredicate.Metadata.BuildStartedOn))
	assert.True(t, time.Now().UTC().After(*slsaPredicate.Metadata.BuildFinishedOn))
	assert.Equal(t, false, slsaPredicate.Metadata.Reproducible)

	// completeness
	assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Environment)
	assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Materials)
	assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Parameters)

	// materials
	assert.Equal(t, expectedMaterials, slsaPredicate.Materials)
}

func toDeps() []build.Dependency {
	return []build.Dependency{
		{
			Coordinates: fmt.Sprintf("%s:%s", "groupId", "artifactId"),
			Version:     "v01",
			CheckSum: build.CheckSum{
				Algorithm: "todo",
				Digest:    "todo",
			},
		},
	}
}

func toArtDeps(deps []build.Dependency) *build.ArtifactDependencies {
	return &build.ArtifactDependencies{
		Cmd:         "lang list:deps",
		RuntimeDeps: deps,
	}
}

func toVcsEnvironment() *vcs.Environment {
	return &vcs.Environment{
		GitHubContext: vcs.GitHubContext{
			Repository: "nais/salsa",
			RunId:      "1234",
			SHA:        "4321",
			Workflow:   "Create a provenance",
			ServerUrl:  "https://github.com",
		},
		Event: vcs.Event{
			Inputs: []byte("some vents"),
		},
		RunnerContext: vcs.RunnerContext{
			OS:        "Linux",
			Temp:      "/home/runner/work/_temp",
			ToolCache: "/opt/hostedtoolcache",
		},
	}
}

func toExpectedMaterials() []slsa.ProvenanceMaterial {
	return []slsa.ProvenanceMaterial{
		{
			URI: "pkg:groupId:artifactId:v01",
			Digest: slsa.DigestSet{
				"todo": "todo",
			},
		},
		{
			URI: "git+https://github.com/nais/salsa",
			Digest: slsa.DigestSet{
				digest.AlgorithmSHA1: "4321",
			},
		},
	}
}

func toExpectedConfigSource() slsa.ConfigSource {
	return slsa.ConfigSource{
		URI: "git+https://github.com/nais/salsa",
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: "4321",
		},
		EntryPoint: "Create a provenance",
	}
}
