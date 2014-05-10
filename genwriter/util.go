package genwriter

func remove(a []string, s ...string) (result []string) {
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
