// Copyright 2016 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package ini

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBOM(t *testing.T) {
	t.Run("test handling BOM", func(t *testing.T) {
		t.Run("UTF-8-BOM", func(t *testing.T) {
			f, err := Load("testdata/UTF-8-BOM.ini")
			require.NoError(t, err)
			require.NotNil(t, f)

			assert.Equal(t, "example@email.com", f.Section("author").Key("E-MAIL").String())
		})

		t.Run("UTF-16-LE-BOM", func(t *testing.T) {
			f, err := Load("testdata/UTF-16-LE-BOM.ini")
			require.NoError(t, err)
			require.NotNil(t, f)
		})

		t.Run("UTF-16-BE-BOM", func(t *testing.T) {
		})
	})
}

func TestBadLoad(t *testing.T) {
	t.Run("load with bad data", func(t *testing.T) {
		t.Run("bad section name", func(t *testing.T) {
			_, err := Load([]byte("[]"))
			require.Error(t, err)

			_, err = Load([]byte("["))
			require.Error(t, err)
		})

		t.Run("bad keys", func(t *testing.T) {
			_, err := Load([]byte(`"""name`))
			require.Error(t, err)

			_, err = Load([]byte(`"""name"""`))
			require.Error(t, err)

			_, err = Load([]byte(`""=1`))
			require.Error(t, err)

			_, err = Load([]byte(`=`))
			require.Error(t, err)

			_, err = Load([]byte(`name`))
			require.Error(t, err)
		})

		t.Run("bad values", func(t *testing.T) {
			_, err := Load([]byte(`name="""Unknwon`))
			require.Error(t, err)
		})
	})
}

func TestCommentInsideQuotes(t *testing.T) {
	t.Run("preserve comment characters inside single quoted values", func(t *testing.T) {
		f, err := Load([]byte(`
[background]
primary-color   = '#e18477'
secondary-color = '#000000'
`))
		require.NoError(t, err)
		require.NotNil(t, f)

		assert.Equal(t, "#e18477", f.Section("background").Key("primary-color").String())
		assert.Equal(t, "#000000", f.Section("background").Key("secondary-color").String())
	})

	t.Run("preserve comment characters inside double quoted values", func(t *testing.T) {
		f, err := Load([]byte(`
[background]
primary-color   = "#e18477"
secondary-color = "#000000"
`))
		require.NoError(t, err)
		require.NotNil(t, f)

		assert.Equal(t, "#e18477", f.Section("background").Key("primary-color").String())
		assert.Equal(t, "#000000", f.Section("background").Key("secondary-color").String())
	})

	t.Run("quoted values with trailing inline comments", func(t *testing.T) {
		f, err := Load([]byte(`
[background]
primary-color   = '#e18477' # inline comment
secondary-color = "#000000" ; another comment
`))
		require.NoError(t, err)
		require.NotNil(t, f)

		assert.Equal(t, "#e18477", f.Section("background").Key("primary-color").String())
		assert.Equal(t, "# inline comment", f.Section("background").Key("primary-color").Comment)
		assert.Equal(t, "#000000", f.Section("background").Key("secondary-color").String())
		assert.Equal(t, "; another comment", f.Section("background").Key("secondary-color").Comment)
	})

	t.Run("preserve surrounded quotes with comment characters inside", func(t *testing.T) {
		f, err := LoadSources(LoadOptions{
			PreserveSurroundedQuote: true,
		}, []byte(`
[background]
primary-color   = '#e18477'
secondary-color = "#000000"
`))
		require.NoError(t, err)
		require.NotNil(t, f)

		assert.Equal(t, "'#e18477'", f.Section("background").Key("primary-color").String())
		assert.Equal(t, "\"#000000\"", f.Section("background").Key("secondary-color").String())
	})

	t.Run("space before inline comment with comment characters inside quotes", func(t *testing.T) {
		f, err := LoadSources(LoadOptions{
			SpaceBeforeInlineComment: true,
		}, []byte(`
[background]
primary-color = '#e18477' # comment
with-space    = 'color # 1' # comment 2
`))
		require.NoError(t, err)
		require.NotNil(t, f)

		assert.Equal(t, "#e18477", f.Section("background").Key("primary-color").String())
		assert.Equal(t, "# comment", f.Section("background").Key("primary-color").Comment)
		assert.Equal(t, "color # 1", f.Section("background").Key("with-space").String())
		assert.Equal(t, "# comment 2", f.Section("background").Key("with-space").Comment)
	})

	t.Run("url containing hash symbol inside quotes", func(t *testing.T) {
		f, err := Load([]byte(`
[links]
docs  = 'https://example.com/page#overview'
issue = "https://github.com/go-ini/ini/issues#393" ; see issue
`))
		require.NoError(t, err)
		require.NotNil(t, f)

		assert.Equal(t, "https://example.com/page#overview", f.Section("links").Key("docs").String())
		assert.Equal(t, "https://github.com/go-ini/ini/issues#393", f.Section("links").Key("issue").String())
		assert.Equal(t, "; see issue", f.Section("links").Key("issue").Comment)
	})

	t.Run("unquoted contraction with inline comment", func(t *testing.T) {
		f, err := Load([]byte(`
[general]
action = don't # comment
`))
		require.NoError(t, err)
		require.NotNil(t, f)

		assert.Equal(t, "don't", f.Section("general").Key("action").String())
		assert.Equal(t, "# comment", f.Section("general").Key("action").Comment)
	})
}
