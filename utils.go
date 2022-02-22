package main

func getMapKeys(myMap map[string]string) []string {
	keys := make([]string, len(myMap))

	i := 0
	for k := range myMap {
		keys[i] = k
		i++
	}
	return keys
}
func getMappingKeys(myMap map[string][]Mapping) []string {
	keys := make([]string, len(myMap))

	i := 0
	for k := range myMap {
		keys[i] = k
		i++
	}
	return keys
}

type response func(int, string)
