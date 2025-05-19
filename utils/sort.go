package utils

import "sort"

// Returns a slice containing the map's values, sorted based on
// the alphabetical order of the map's keys.
func SortMapValuesByKey[V any](m map[string]V) []V {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	sortedValues := make([]V, 0, len(m))

	for _, k := range keys {
		sortedValues = append(sortedValues, m[k])
	}

	return sortedValues
}
