package test

import (
	"github.com/nais/salsa/pkg/build"
	"reflect"
	"testing"
)

func Dependency(coordinates, version, algo, checksum string) build.Dependency {
	return build.Dependence(coordinates, version,
		build.Verification(algo, checksum),
	)
}

func AssertEqual(t *testing.T, got, want map[string]build.Dependency) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}