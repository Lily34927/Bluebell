package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

/* 投票的几种情况：
direction = 1时，有两种情况：
	1.之前没有投过票，现在投赞成票  --》 更新分数和投票记录  差值的绝对值：1  +432
	2.之前投反对票，现在改投赞成票  --》 更新分数和投票记录  差值的绝对值：2  +432 * 2
direction = 0时，有两种情况：
	1.之前透过反对票，现在要取消投票  --》 更新分数和投票记录  差值的绝对值：1  +432
	2.之前投过赞成票，现在要取消投票  --》 更新分数和投票记录  差值的绝对值：1  -432
direction = -1，有两种情况：
	1.之前没有投过票，现在投反对票  --》 更新分数和投票记录  差值的绝对值：1  -432
	2.之前投赞成票，现在改投反对票  --》 更新分数和投票记录  差值的绝对值：2  -432 * 2

投票的限制：
	每个帖子自发表之日起，一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成拍哦书及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(postID)),
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(postID)),
	})

	// 更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, communityID)
	_, err := pipeline.Exec()
	return err
}
func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票的限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if (float64(time.Now().Unix()) - postTime) > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2. 更新帖子的分数
	// 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()

	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	}
	op = -1

	diff := math.Abs(ov - value) // 计算两次投票的差值

	// 2和3需要放到一个pipeline事务中操作
	pipeline := client.TxPipeline()
	//_, err := client.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID).Result()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		// ZREM函数的作用是从有序集合中删除一个或多个成员。
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	}
	pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
		Score:  value, // 赞成票还是反对票
		Member: userID,
	})
	_, err := pipeline.Exec()
	return err
}
