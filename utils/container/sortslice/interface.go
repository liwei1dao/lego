package sortslice

//排序工具
func quickSort(arr []interface{}, start, end int, compete func(a interface{}, b interface{}) int8) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for compete(arr[i], key) == -1 {
				i++
			}
			for compete(arr[j], key) == 1 {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			quickSort(arr, start, j, compete)
		}
		if end > i {
			quickSort(arr, i, end, compete)
		}
	}
}

func Sort(a []interface{}, compete func(a interface{}, b interface{}) int8) {
	if len(a) < 2 {
		return
	}
	quickSort(a, 0, len(a)-1, compete)
}
