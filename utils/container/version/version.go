package version

import (
	"strings"
)

func CompareStrVer(verA, verB string) int8 {
	verStrArrA := spliteStrByNet(verA)
	verStrArrB := spliteStrByNet(verB)
	lenStrA := len(verStrArrA)
	lenStrB := len(verStrArrB)
	if lenStrA > lenStrB {
		return 1
	}
	if lenStrA < lenStrB {
		return -1
	}
	return compareArrStrVers(verStrArrA, verStrArrB)
}

// 比较版本号字符串数组
func compareArrStrVers(verA, verB []string) int8 {
	for index, _ := range verA {
		littleResult := compareLittleVer(verA[index], verB[index])
		if littleResult != 0 {
			return littleResult
		}
	}
	return 0
}

//
// 比较小版本号字符串
//
func compareLittleVer(verA, verB string) int8 {
	bytesA := []byte(verA)
	bytesB := []byte(verB)
	lenA := len(bytesA)
	lenB := len(bytesB)
	if lenA > lenB {
		return 1
	}
	if lenA < lenB {
		return -1
	}
	//如果长度相等则按byte位进行比较
	return compareByBytes(bytesA, bytesB)
}

// 按byte位进行比较小版本号
func compareByBytes(verA, verB []byte) int8 {
	for index, _ := range verA {
		if verA[index] > verB[index] {
			return 1
		}
		if verA[index] < verB[index] {
			return -1
		}
	}
	return 0
}

// 按“.”分割版本号为小版本号的字符串数组
func spliteStrByNet(strV string) []string {
	return strings.Split(strV, ". ")
}
