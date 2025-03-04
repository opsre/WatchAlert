package tools

// GetSliceDifference 获取 slice1 中存在, slice2 不存在的差异key
func GetSliceDifference(slice1 []string, slice2 []string) []string {
	difference := []string{}

	// 遍历缓存
	for _, item1 := range slice1 {
		found := false
		// 遍历当前key
		for _, item2 := range slice2 {
			if item1 == item2 {
				found = true
				break
			}
		}
		// 添加到差异切片中
		if !found {
			difference = append(difference, item1)
		}
	}

	return difference
}

// GetSliceSame 获取相同key, 当slice1中存在, slice2也存在则标记为正在告警中撤销告警恢复
func GetSliceSame(slice1 []string, slice2 []string) []string {
	same := []string{}
	for _, item1 := range slice1 {
		for _, item2 := range slice2 {
			if item1 == item2 {
				same = append(same, item1)
			}
		}
	}
	return same
}
