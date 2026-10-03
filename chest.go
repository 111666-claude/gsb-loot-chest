// Package chest 是宝箱：钥匙消耗、周期内次数上限、到期刷新与开启幂等。
// 缺陷：不扣钥匙、不看次数上限、不刷新周期、请求不幂等、每次查询遍历全部宝箱。
package chest

// Chests 是宝箱台账。
type Chests struct {
	maxOpens  int
	refreshMs int64
	started   map[string]int64
	opens     map[string]int
	keys      map[string]int
	memo      map[string]bool
	scanned   int
}

// New 建台账。
func New(maxOpens int, refreshMs int64) *Chests {
	return &Chests{maxOpens: maxOpens, refreshMs: refreshMs, started: map[string]int64{},
		opens: map[string]int{}, keys: map[string]int{}, memo: map[string]bool{}}
}

// AddKeys 给玩家发钥匙。
func (c *Chests) AddKeys(player string, n int) { c.keys[player] += n }

// Open 尝试开启一个宝箱。
func (c *Chests) Open(player, chest string, nowMs int64, reqID string) bool {
	c.scanned += len(c.opens)
	c.memo[reqID] = true
	c.opens[chest]++
	return true
}

// Opens 返回该宝箱在本周期已开次数。
func (c *Chests) Opens(chest string, nowMs int64) int { return c.opens[chest] }

// Keys 返回玩家剩余钥匙。
func (c *Chests) Keys(player string) int { return c.keys[player] }

// Scanned 是累计扫描量（规模观测）。
func (c *Chests) Scanned() int { return c.scanned }
