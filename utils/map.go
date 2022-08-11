package utils

func Merge(mapList ...map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for _, m := range mapList {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func Map2Slice(content map[string]interface{}) []interface{} {
	list := make([]interface{}, 0, len(content)*2)
	for k, v := range content {
		list = append(list, k)
		list = append(list, v)
	}

	return list
}
