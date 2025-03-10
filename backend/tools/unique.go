package tools

func Unique(stringslice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringslice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}