package recom

import (
	"fmt"
	"testing"
)

/*
测试推荐系统
		1	2	3	4	5
1001		1	5
1002	5				1
1003			1		5
*/

func Test_Sys(t *testing.T) {
	sys, err := NewSys(SetUserIds([]uint32{1001, 1002, 1003, 1001, 1002, 1003}), SetItemIds([]uint32{2, 1, 3, 3, 5, 5}), SetRatings([]float64{1, 5, 1, 5, 1, 5}))
	if err != nil {
		fmt.Printf("sys err:%v", err)
	}
	sys.Wait()
	items := sys.RecommendItems(1003, 10)
	fmt.Printf("items:%v", items)
}
