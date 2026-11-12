// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package text_test

import (
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/language"
	"code.superseriousbusiness.org/gotosocial/internal/text"
)

func TestContentToContentLanguage(t *testing.T) {
	type testcase struct {
		content           gtsmodel.Content
		instanceLanguages language.Languages
		expectedContent   string
		expectedLang      string
	}

	ctx := t.Context()

	for i, testcase := range []testcase{
		{
			content: gtsmodel.Content{
				Content:    "hello world",
				ContentMap: nil,
			},
			expectedContent: "hello world",
			expectedLang:    "",
		},
		{
			content: gtsmodel.Content{
				MediaType:  "text/markdown",
				Content:    "**hello world**",
				ContentMap: nil,
			},
			expectedContent: "<p><strong>hello world</strong></p>",
			expectedLang:    "",
		},
		{
			content: gtsmodel.Content{
				Content: "",
				ContentMap: map[string]string{
					"en": "hello world",
				},
			},
			expectedContent: "hello world",
			expectedLang:    "en",
		},
		{
			content: gtsmodel.Content{
				Content: "bonjour le monde",
				ContentMap: map[string]string{
					"en": "hello world",
					"fr": "bonjour le monde",
				},
			},
			expectedContent: "bonjour le monde",
			expectedLang:    "fr",
		},
		{
			content: gtsmodel.Content{
				Content: "bonjour le monde",
				ContentMap: map[string]string{
					"en": "hello world",
				},
			},
			expectedContent: "bonjour le monde",
			expectedLang:    "",
		},
		{
			content: gtsmodel.Content{
				Content: "",
				ContentMap: map[string]string{
					"en": "hello world",
					"ru": "Привет, мир!",
					"nl": "hallo wereld!",
					"ca": "Hola món!",
				},
			},
			instanceLanguages: language.Languages{
				{TagStr: "en"},
				{TagStr: "ca"},
			},
			expectedContent: "hello world",
			expectedLang:    "en",
		},
		{
			content: gtsmodel.Content{
				Content: "",
				ContentMap: map[string]string{
					"en": "hello world",
					"ru": "Привет, мир!",
					"nl": "hallo wereld!",
					"ca": "Hola món!",
				},
			},
			instanceLanguages: language.Languages{
				{TagStr: "ca"},
				{TagStr: "en"},
			},
			expectedContent: "Hola món!",
			expectedLang:    "ca",
		},
	} {
		langs, err := language.InitLangs(testcase.instanceLanguages.TagStrs())
		if err != nil {
			t.Fatal(err)
		}
		config.SetInstanceLanguages(langs)

		content, language := text.ContentToContentLanguage(ctx, testcase.content)
		if content != testcase.expectedContent {
			t.Errorf(
				"test %d expected content '%s' got '%s'",
				i, testcase.expectedContent, content,
			)
		}

		if language != testcase.expectedLang {
			t.Errorf(
				"test %d expected language '%s' got '%s'",
				i, testcase.expectedLang, language,
			)
		}
	}
}
