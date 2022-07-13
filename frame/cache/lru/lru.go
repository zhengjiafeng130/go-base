package lru

import "container/list"

// FIFO | 先进先出 -> 队列 | 没有使用频率,命中率低
// LFU  | 维护频率排序 | 维护消耗大,受历史数据影响大
// LRU  | 双向链表+map

// LRU cache, not concurrent safe
type Cache struct {
	maxBytes int64
	nbytes int64
	ll  *list.List
	cache map[string]*list.Element
}
