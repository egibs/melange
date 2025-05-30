// Copyright 2023 Chainguard, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package build

import (
	"bytes"
	"testing"
	"time"

	"chainguard.dev/melange/pkg/config"

	"github.com/stretchr/testify/require"
)

func Test_removeSelfProvidedDeps_WithVersionedProvides(t *testing.T) {
	provides := []string{"so:libfoo.so.3=3", "so:libbar.so.2=2"}
	depends := []string{"so:libbaz.so.4", "so:libfoo.so.3"}

	final := removeSelfProvidedDeps(depends, provides)

	require.Equal(t, len(final), 1, "only one depend in the list")
	require.Equal(t, final[0], "so:libbaz.so.4", "remaining depend should be so:libbaz.so.4")
}

func Test_removeSelfProvidedDeps_WithoutVersionedProvides(t *testing.T) {
	provides := []string{"so:libfoo.so.3", "so:libbar.so.2"}
	depends := []string{"so:libbaz.so.4", "so:libfoo.so.3"}

	final := removeSelfProvidedDeps(depends, provides)

	require.Equal(t, len(final), 1, "only one depend in the list")
	require.Equal(t, final[0], "so:libbaz.so.4", "remaining depend should be so:libbaz.so.4")
}

func Test_removeSelfProvidedDeps_WithEmptyProvides(t *testing.T) {
	provides := []string{}
	depends := []string{"so:libbaz.so.4", "so:libfoo.so.3"}

	final := removeSelfProvidedDeps(depends, provides)

	require.Equal(t, len(final), 2, "only two depends in the list")
	require.Equal(t, final[0], "so:libbaz.so.4", "first remaining depend should be so:libbaz.so.4")
	require.Equal(t, final[1], "so:libfoo.so.3", "second remaining depend should be so:libfoo.so.3")
}

func Test_GenerateControlData(t *testing.T) {
	pkg := &config.Package{
		Version: "1.2.3",
		Epoch:   4,
	}

	tests := []struct {
		name string
		pb   *PackageBuild
		want string
	}{{
		name: "simple",
		pb: &PackageBuild{
			Build: &Build{
				SourceDateEpoch: time.Unix(0, 0),
				Namespace:       "wolfi",
			},
			Origin:        pkg,
			PackageName:   "glibc",
			Arch:          "aarch64",
			InstalledSize: 666,
			OriginName:    "bigbang",
			Description:   "I'm a unit test",
			URL:           "https://chainguard.dev",
			Commit:        "deadbeef",
			DataHash:      "baadf00d",
		},
		want: `# Generated by melange
pkgname = glibc
pkgver = 1.2.3-r4
arch = aarch64
size = 666
origin = bigbang
pkgdesc = I'm a unit test
url = https://chainguard.dev
commit = deadbeef
maintainer = wolfi
datahash = baadf00d
`,
	}, {
		name: "source date epoch",
		pb: &PackageBuild{
			Build: &Build{
				SourceDateEpoch: time.Unix(12345678, 0),
			},
			Origin:        pkg,
			PackageName:   "glibc",
			Arch:          "aarch64",
			InstalledSize: 666,
			OriginName:    "bigbang",
			Description:   "I'm a unit test",
			URL:           "https://chainguard.dev",
			Commit:        "deadbeef",
			DataHash:      "baadf00d",
		},
		want: `# Generated by melange
pkgname = glibc
pkgver = 1.2.3-r4
arch = aarch64
size = 666
origin = bigbang
pkgdesc = I'm a unit test
url = https://chainguard.dev
commit = deadbeef
builddate = 12345678
datahash = baadf00d
`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotBuf := bytes.NewBuffer(nil)
			err := test.pb.GenerateControlData(gotBuf)
			if err != nil {
				t.Fatalf("GenerateControlData() = %v", err)
			}
			got := gotBuf.String()
			if test.want != got {
				t.Errorf("GenerateControlData() = %s, wanted %s", got, test.want)
			}
		})
	}
}
