package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node

func Init(startTime string, matchineID int64) (err error) {
	// 传递过来的时间因子，从什么时候开始，可以用69年 -- startTime
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000 // 初始化开始的时间，毫秒级别
	node, err = sf.NewNode(matchineID) // 指定机器的ID
	return
}

func GenID() int64 {
	// 生成ID
	return node.Generate().Int64()
}

//func main() {
//	if err := Init("2020-07-01", 1); err != nil {
//		fmt.Printf("init failed, err %v\n", err)
//		return
//	}
//	id := GenID()
//	fmt.Println(id)
//}
