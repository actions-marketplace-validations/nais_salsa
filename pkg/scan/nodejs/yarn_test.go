package nodejs

import (
	"testing"

	"github.com/nais/salsa/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestYarnLockParsing(t *testing.T) {

	got := YarnDeps(yarnlLockContents)
	want := []scan.Dependency{
		yarnDep("js-tokens", "4.0.0", "sha512", "RdJUflcE3cUzKiMqQgsCu06FPu9UdIJO0beYbPhHN4k6apgJtifcoCtT9bcxOpYBtpD2kCM6Sbzg4CausW/PKQ=="),
		yarnDep("loose-envify", "1.4.0", "sha512", "lyuxPGr/Wfhrlem2CL/UcnUc1zcqKAImBDzukY7Y5F/yQiNdko6+fRLevlw1HgMySw7f611UIY408EtxRSoK3Q=="),
		yarnDep("object-assign", "4.1.1", "sha1", "IQmtx5ZYh8/AXLvUQsrIv7s2CGM="),
		yarnDep("react", "17.0.2", "sha512", "gnhPt75i/dq/z3/6q/0asP78D0u592D5L1pd7M8P+dck6Fu/jJeL6iVVK23fptSUZj8Vjf++7wXA8UNclGQcbA=="),
	}

	if !assert.ElementsMatch(t, got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func yarnDep(coordinates, version, alg, digest string) scan.Dependency {
	return scan.Dependency{
		Coordinates: coordinates,
		Version:     version,
		CheckSum: scan.CheckSum{
			Algorithm: alg,
			Digest:    digest,
		},
	}
}

const yarnlLockContents = `
# THIS IS AN AUTOGENERATED FILE. DO NOT EDIT THIS FILE DIRECTLY.
# yarn lockfile v1

"js-tokens@^3.0.0 || ^4.0.0":
  version "4.0.0"
  resolved "https://registry.yarnpkg.com/js-tokens/-/js-tokens-4.0.0.tgz#19203fb59991df98e3a287050d4647cdeaf32499"
  integrity sha512-RdJUflcE3cUzKiMqQgsCu06FPu9UdIJO0beYbPhHN4k6apgJtifcoCtT9bcxOpYBtpD2kCM6Sbzg4CausW/PKQ==

"loose-envify@^1.1.0":
  version "1.4.0"
  resolved "https://registry.yarnpkg.com/loose-envify/-/loose-envify-1.4.0.tgz#71ee51fa7be4caec1a63839f7e682d8132d30caf"
  integrity sha512-lyuxPGr/Wfhrlem2CL/UcnUc1zcqKAImBDzukY7Y5F/yQiNdko6+fRLevlw1HgMySw7f611UIY408EtxRSoK3Q==
  dependencies:
    js-tokens "^3.0.0 || ^4.0.0"

"object-assign@^4.1.1":
  version "4.1.1"
  resolved "https://registry.yarnpkg.com/object-assign/-/object-assign-4.1.1.tgz#2109adc7965887cfc05cbbd442cac8bfbb360863"
  integrity sha1-IQmtx5ZYh8/AXLvUQsrIv7s2CGM=

"react@>=0.5.1, react@^17.0.2":
  version "17.0.2"
  resolved "https://registry.yarnpkg.com/react/-/react-17.0.2.tgz#d0b5cc516d29eb3eee383f75b62864cfb6800037"
  integrity sha512-gnhPt75i/dq/z3/6q/0asP78D0u592D5L1pd7M8P+dck6Fu/jJeL6iVVK23fptSUZj8Vjf++7wXA8UNclGQcbA==
  dependencies:
    loose-envify "^1.1.0"
    object-assign "^4.1.1"
`
