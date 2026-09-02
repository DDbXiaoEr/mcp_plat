// Copyright (C) 2026 Zhaoquan Wang
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package i18n provides lightweight, dependency-free message translation.
//
// Chinese is the canonical source language written across the codebase
// (handlers, middleware and service error templates). English translations are
// produced on demand from an Accept-Language-derived language tag by replacing
// known Chinese phrases with their English counterparts. Any phrase without an
// entry is kept verbatim, so localization is best-effort and safe.
package i18n

// Author: deepseek-v4-pro / opencode

import (
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Lang normalizes an Accept-Language header value to a supported language tag.
// Only English ("en*") is translated; everything else falls back to Chinese.
func Lang(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return "zh"
	}
	first := strings.ToLower(strings.TrimSpace(strings.Split(header, ",")[0]))
	if strings.HasPrefix(first, "en") {
		return "en"
	}
	return "zh"
}

// Translate converts a Chinese-origin message into the given language.
// lang should be one of the tags produced by Lang; any other value returns text
// unchanged. Translation is a longest-match-first substitution of the phrases
// registered in zhEnMessages, which tolerates interpolated dynamic values.
func Translate(lang, text string) string {
	if lang != "en" || !containsHan(text) {
		return text
	}
	for _, p := range compiledPatterns {
		if p.re.MatchString(text) {
			return p.re.ReplaceAllString(text, p.tmpl)
		}
	}
	for _, zh := range sortedPhrases {
		if strings.Contains(text, zh) {
			text = strings.ReplaceAll(text, zh, zhEnMessages[zh])
		}
	}
	return text
}

func containsHan(s string) bool {
	for _, r := range s {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

var sortedPhrases []string

var compiledPatterns []struct {
	re   *regexp.Regexp
	tmpl string
}

func init() {
	sortedPhrases = make([]string, 0, len(zhEnMessages))
	for zh := range zhEnMessages {
		sortedPhrases = append(sortedPhrases, zh)
	}
	sort.Slice(sortedPhrases, func(i, j int) bool {
		return utf8.RuneCountInString(sortedPhrases[i]) > utf8.RuneCountInString(sortedPhrases[j])
	})

	compiledPatterns = make([]struct {
		re   *regexp.Regexp
		tmpl string
	}, 0, len(zhEnPatterns))
	for _, p := range zhEnPatterns {
		compiledPatterns = append(compiledPatterns, struct {
			re   *regexp.Regexp
			tmpl string
		}{re: regexp.MustCompile(p.re), tmpl: p.tmpl})
	}
}
