package etc

// Discord proxy bot
// Copyright (C) 2020  Starshine System

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import (
	"fmt"
	"strings"
)

// TrimPrefixesSpace trims all given prefixes as well as whitespace from the given string
func TrimPrefixesSpace(s string, prefixes ...string) string {
	for _, prefix := range prefixes {
		s = strings.TrimPrefix(s, prefix)
	}
	s = strings.TrimSpace(s)
	return s
}

// HasAnyPrefix checks if the string has *any* of the given prefixes
func HasAnyPrefix(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// HasAnySuffix checks if the string has *any* of the given suffixes
func HasAnySuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// PrintfAll ...
func PrintfAll(template string, in []string) []string {
	out := make([]string, 0)
	for _, i := range in {
		out = append(out, fmt.Sprintf(template, i))
	}

	return out
}
