package vuid

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

/**
1 符号位  |  41 时间戳                                      | 10 工作机器    | 12 （毫秒内）自增ID
---------|----------------------------------------------------------------|--------------------
0        | 0 0000000 00000000 00000000 00000000 00000000  | 00000000 00   | 0000 00000000
---------|----------------------------------------------------------------|--------------------
0或者1   |    2199023255551                               |    1023       | 4095

解释：
1 bit：不使用，可以是 1 或 0
41 bit：记录时间戳 (当前时间戳减去用户设置的初始时间，毫秒表示)，可记录最多 69 年的时间戳数据
10 bit：用来记录分布式节点 ID，一般每台机器一个唯一 ID，也可以多进程每个进程一个唯一 ID，最大可部署 1024 个节点
12 bit：序列号，用来记录不同 ID 同一毫秒时的序列号，最多可生成 4096 个序列号

*/

const (
	workIDBits uint8 = 10                      // 节点 ID 的位数
	stepBits   uint8 = 12                      // 序列号的位数
	nodeMax    int64 = -1 ^ (-1 << workIDBits) // 节点 ID 的最大值，用于检测溢出
	stepMax    int64 = -1 ^ (-1 << stepBits)   // 序列号的最大值，用于检测溢出
	timeShift        = workIDBits + stepBits   // 时间戳向左的偏移量
	nodeShift        = stepBits                // 节点 ID 向左的偏移量
)

var (
	InitEpoch        = "2019-01-01 00:00:00" //初始时间
	TimeFormatLayout = "2006-01-02 15:04:05" //时间统一个格式
)

type SetEpoch func()

type SnowFlake struct {
	mu        sync.Mutex // 添加互斥锁，保证并发安全
	epoch     int64      // 初始时间戳
	timestamp int64      // 时间戳部分
	work      int64      // 工作机器部分
	step      int64      // 序列号部分
}

func NewSnowFlake(workID int64) (*SnowFlake, error) {
	if workID > nodeMax {
		return nil, errors.New("work id exceeds maximum for " + strconv.Itoa(int(nodeMax)))
	}
	ep, err := time.ParseInLocation(TimeFormatLayout, InitEpoch, time.Local)
	return &SnowFlake{
		work:      workID,
		timestamp: time.Now().UnixNano() / 1e6,
		step:      0,
		epoch:     ep.UnixNano() / 1000000,
	}, err
}

// @Desc 	: 设置初始的时间戳，后面生成的Uuid初始时间值为这个
// @param	: timeStr 格式为 YYYY-MM-DD HH:mm:ss
func (sel *SnowFlake) SetEpoch(timeStr string) error {
	tm, err := time.ParseInLocation(TimeFormatLayout, timeStr, time.Local)
	// 设置初始时间的时间戳 (毫秒表示)，这个可以随意设置 ，比现在的时间靠前即可。
	sel.epoch = tm.UnixNano() / 1000000
	return err
}

// @Desc	: 生成全局唯一ID
func (sel *SnowFlake) Generate() (id int64) {

	// 保证并发安全, 加锁
	sel.mu.Lock()
	// 方法运行完毕后解锁
	defer sel.mu.Unlock()
	// 获取当前时间的时间戳 (毫秒数显示)
	nowStamp := time.Now().UnixNano() / 1e6

	// 当前时间戳是否相等
	if sel.timestamp == nowStamp {
		// 自曾序列号
		sel.step++
		if sel.step > stepMax {
			// 等待本毫秒结束
			for nowStamp <= sel.timestamp {
				//fmt.Println("timestamp same")
				nowStamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 不是在同一个毫秒内就
		sel.step = 0
	}
	sel.timestamp = nowStamp
	id = ((nowStamp - sel.epoch) << timeShift) | (sel.work << nodeShift) | (sel.step)
	return
}
