// Copyright 2022 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

// DeduplicateStrings deduplicate string list, empty items are dropped
func DeduplicateStrings(src []string) []string {
	m := make(map[string]struct{}, len(src))
	dst := make([]string, 0, len(src))

	for _, v := range src {
		// Skip empty items
		if len(v) == 0 {
			continue
		}
		// Skip duplicates
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		dst = append(dst, v)
	}

	return dst
}
