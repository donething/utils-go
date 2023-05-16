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
//
// 当没有找到需要删除的元素时，返回原数组和 false
func DelItem[T any](data []T, equal func(item T) bool) ([]T, bool) {
	index := FindIndex(data, equal)
	if index == -1 {
		return data, false
	}

	return append(data[:index], data[index+1:]...), true
}
