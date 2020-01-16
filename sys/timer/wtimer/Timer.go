package wtimer

import (
	"github.com/liwei1dao/lego/utils/id"
	"sync"
	"time"
)

const wheel_cnt uint8 = 5 //时间轮数量5个

func NewTimer(inteval time.Duration) *Timer {
	TimerSys := new(Timer)
	TimerSys.element_cnt_per_wheel = [wheel_cnt]uint32{256, 64, 64, 64, 64}                          //每个时间轮的槽(元素)数量。在 256+64+64+64+64 = 512 个槽中，表示的范围为 2^32
	TimerSys.right_shift_per_wheel = [wheel_cnt]uint32{8, 6, 6, 6, 6}                                //当指针指向当前时间轮最后一位数，再走一位就需要向上进位。每个时间轮进位的时候，使用右移的方式，最快实现进位。这里是每个轮的进位二进制位数
	TimerSys.base_per_wheel = [wheel_cnt]uint32{1, 256, 256 * 64, 256 * 64 * 64, 256 * 64 * 64 * 64} //记录每个时间轮指针当前指向的位置
	TimerSys.TimerMap = make(map[string]*Node)                                                       //保存待执行的计时器，方便按链表节点指针地址直接删除定时器
	var bucket_no uint8 = 0
	for bucket_no = 0; bucket_no < wheel_cnt; bucket_no++ {
		var i uint32 = 0
		for ; i < TimerSys.element_cnt_per_wheel[bucket_no]; i++ {
			TimerSys.timewheels[bucket_no] = append(TimerSys.timewheels[bucket_no], new(Node))
		}
	}
	go run(TimerSys, inteval)
	return TimerSys
}

func run(t *Timer, inteval time.Duration) {
	for {
		time.Sleep(inteval)
		t.step()
	}
}

type TimerTask struct {
	Key         string                       //定时器名称
	Inteval     uint32                       //时间间隔，即以插入该定时器的时间为起点，Inteval秒之后执行回调函数DoSomething()。例如进程插入该定时器的时间是2015-04-05 10:23:00，Inteval=5，则执行DoSomething()的时间就是2015-04-05 10:23:05。
	DoSomething func(string, ...interface{}) //自定义事件处理函数，需要触发的事件
	Args        []interface{}                //上述函数的输入参数
}

type Timer struct {
	element_cnt_per_wheel [wheel_cnt]uint32
	right_shift_per_wheel [wheel_cnt]uint32
	base_per_wheel        [wheel_cnt]uint32
	mutex                 sync.Mutex
	rwmutex               sync.RWMutex
	newest                [wheel_cnt]uint32  //每个时间轮当前指针所指向的位置
	timewheels            [wheel_cnt][]*Node //定义5个时间轮
	TimerMap              map[string]*Node   //保存待执行的计时器，方便按链表节点指针地址直接删除定时器
}

func (this *Timer) Add(inteval uint32, handler func(string, ...interface{}), args ...interface{}) string {
	timerid := id.GenerateID().String()
	this.setTimer(timerid, inteval, handler, args...)
	return timerid
}

func (this *Timer) Remove(key string) {
	this.rwmutex.Lock()
	node := this.TimerMap[key]
	delete(this.TimerMap, key)
	Delete(node)
	this.rwmutex.Unlock()
}

func (this *Timer) setTimer(key string, inteval uint32, handler func(string, ...interface{}), args ...interface{}) {
	if inteval <= 0 {
		return
	}
	var bucket_no uint8 = 0
	var offset uint32 = inteval
	var left uint32 = inteval
	for offset >= this.element_cnt_per_wheel[bucket_no] { //偏移量大于当前时间轮容量，则需要向高位进位
		offset >>= this.right_shift_per_wheel[bucket_no] //计算高位的值。偏移量除以低位的进制。比如低位当前是256，则右移8个二进制位，就是除以256，得到的结果是高位的值。
		var tmp uint32 = 1
		if bucket_no == 0 {
			tmp = 0
		}
		left -= this.base_per_wheel[bucket_no] * (this.element_cnt_per_wheel[bucket_no] - this.newest[bucket_no] - tmp)
		bucket_no++
	}
	if offset < 1 {
		return
	}
	if inteval < this.base_per_wheel[bucket_no]*offset {
		return
	}
	left -= this.base_per_wheel[bucket_no] * (offset - 1)
	pos := (this.newest[bucket_no] + offset) % this.element_cnt_per_wheel[bucket_no] //通过类似hash的方式，找到在时间轮上的插入位置

	var node Node
	node.SetData(TimerTask{key, left, handler, []interface{}(args)})

	this.rwmutex.Lock()
	this.TimerMap[key] = this.timewheels[bucket_no][pos].InsertHead(node) //插入定时器
	this.rwmutex.Unlock()
}

func (this *Timer) step() {
	this.rwmutex.RLock()
	//遍历所有桶
	var bucket_no uint8 = 0
	for bucket_no = 0; bucket_no < wheel_cnt; bucket_no++ {
		this.newest[bucket_no] = (this.newest[bucket_no] + 1) % this.element_cnt_per_wheel[bucket_no] //当前指针递增1
		//fmt.Println(newest)
		var head *Node = this.timewheels[bucket_no][this.newest[bucket_no]] //返回当前指针指向的槽位置的表头
		var firstElement *Node = head.Next()
		for firstElement != nil { //链表不为空
			if value, ok := firstElement.Data().(TimerTask); ok { //如果element里面确实存储了Timer类型的数值，那么ok返回true，否则返回false。
				inteval := value.Inteval
				doSomething := value.DoSomething
				args := value.Args
				if nil != doSomething { //有遇到函数为nil的情况，所以这里判断下非nil
					if 0 == bucket_no || 0 == inteval {
						//dolist.PushBack(value) //执行自定义处理函数
						go doSomething(value.Key, args...)
					} else {
						this.setTimer(value.Key, inteval, doSomething, args) //重新插入计时器
					}
				}
				delete(this.TimerMap, value.Key)
				Delete(firstElement) //删除定时器
			}
			firstElement = head.Next() //重新定位到链表第一个元素头
		}
		if 0 != this.newest[bucket_no] { //指针不是0，还未转回到原点，跳出。如果回到原点，则说明转完了一圈，需要向高位进位1，则继续循环入高位步进一步。
			break
		}
	}
	this.rwmutex.RUnlock()
}
