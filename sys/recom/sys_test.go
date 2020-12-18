package recom

import (
	"fmt"
	"testing"
)

/*
测试推荐系统
		1	2	3
1001		5	5
1002	5
1003			5
*/

func Test_Sys(t *testing.T) {
	// sys, err := NewSys(SetUserIds([]uint32{1001, 1002, 1003, 1001, 1002, 1003}), SetItemIds([]uint32{2, 1, 3, 3, 5, 5}), SetRatings([]float64{1, 5, 1, 5, 5, 5}))
	sys, err := NewSys(SetItemIdsScore(map[uint32]map[uint32]float64{
		1:  map[uint32]float64{1002: 5},
		2:  map[uint32]float64{1001: 1},
		10: map[uint32]float64{},
		3:  map[uint32]float64{1001: 4, 1003: 5},
		4:  map[uint32]float64{1002: 5, 1003: 5},
		11: map[uint32]float64{},
		5:  map[uint32]float64{1002: 5},
		6:  map[uint32]float64{1001: 2},
		7:  map[uint32]float64{},
	}))
	if err != nil {
		fmt.Printf("sys err:%v", err)
	}
	sys.Wait()
	items := sys.RecommendItems(1003, 10)
	fmt.Printf("items:%v", items)
}
