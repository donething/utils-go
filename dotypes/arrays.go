package dotypes

// FindIndex 从数组中查找指定项的索引
func FindIndex[T any](data []T, equal func(item T) bool) int {
	for k, v := range data {
		if equal(v) {
			return k
		}
	}

	return -1
}

// DelItem 从数组从删除指定项
func DelItem[T any](data []T, equal func(item T) bool) []T {
	index := FindIndex(data, equal)

	return append(data[:index], data[index+1:]...)
}
