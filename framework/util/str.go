package util

func ConvertStrSlice2Map(data []string) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, v := range data {
		ret[v] = struct{}{}
	}

	return ret
}
