package util

import (
	"fmt"
	"testing"
)

func TestStr(t *testing.T) {
	da := ProcessString(`【类型】车找人
	【时间】明早8点
	【出发】文旅城
	【到达】回龙
	【车位】2
	【备注】18349156327
	 
	
	********************
	转自:文旅拼车群（10元每人）`)
	fmt.Println(da)
}
