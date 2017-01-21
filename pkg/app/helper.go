package app

// UniqueIDs filters out duplicated ID from IDs
func UniqueIDs(in []string) []string {
	if in == nil {
		return []string{}
	}
	idMap := map[string]bool{}
	for _, id := range in {
		idMap[id] = true
	}
	res := make([]string, 0, len(idMap))
	for id := range idMap {
		res = append(res, id)
	}
	return res
}
