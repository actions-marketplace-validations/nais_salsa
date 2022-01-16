package intoto

import (
	"fmt"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"strings"
	"time"
)

func GenerateSlsaPredicate(app App) slsa.ProvenancePredicate {
	return withPredicate(app)
}

func withPredicate(app App) slsa.ProvenancePredicate {
	return slsa.ProvenancePredicate{
		Builder:   slsa.ProvenanceBuilder{},
		BuildType: "yolo",
		Invocation: slsa.ProvenanceInvocation{
			ConfigSource: slsa.ConfigSource{},
			Parameters:   nil,
			Environment:  nil,
		},
		BuildConfig: nil,
		Metadata:    withMetadata(app, false, time.Now().UTC()),
		Materials:   withMaterials(app),
	}
}

func FindMaterials(materials []slsa.ProvenanceMaterial, value string) []slsa.ProvenanceMaterial {
	found := make([]slsa.ProvenanceMaterial, 0)
	for _, m := range materials {
		if find(m, value) {
			found = append(found, m)
		}
	}
	return found
}

func find(material slsa.ProvenanceMaterial, value string) bool {
	return strings.Contains(material.URI, value)
}

func withMetadata(app App, rp bool, buildFinished time.Time) *slsa.ProvenanceMetadata {
	return &slsa.ProvenanceMetadata{
		BuildInvocationID: app.BuildInvocationId,
		BuildStartedOn:    &app.BuildStartedOn,
		BuildFinishedOn:   &buildFinished,
		Completeness:      withCompleteness(false, false),
		Reproducible:      rp,
	}
}

func withCompleteness(environment, materials bool) slsa.ProvenanceComplete {
	return slsa.ProvenanceComplete{
		Environment: environment,
		Materials:   materials,
	}
}

// TODO: use other type of materials aswell, e.g. github actions run in the build
func withMaterials(app App) []slsa.ProvenanceMaterial {
	materials := make([]slsa.ProvenanceMaterial, 0)
	for k, v := range app.Dependencies {
		var uri = fmt.Sprintf("pkg:%s:%s", k, v)
		m := slsa.ProvenanceMaterial{
			URI:    uri,
			Digest: slsa.DigestSet{app.Checksums[k].Algorithm: app.Checksums[k].Digest},
		}
		materials = append(materials, m)
	}
	return materials
}
