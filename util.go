package main

func findDuplicateStrings(a []string) (result []string) {
	found := make(map[string]bool)

	for _, s := range a {
		if found[s] {
			result = append(result, s)
		}
		found[s] = true
	}
	return
}
