package standard

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

func remove(a []string, s []string) (result []string) {
	exists := make(map[string]bool)
	for _, v := range s {
		exists[v] = true
	}

	for _, v := range a {
		if !exists[v] {
			result = append(result, v)
		}
	}
	return
}
