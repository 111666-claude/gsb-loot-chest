package chest

import (
	"strconv"
	"testing"
)

func TestKeysStartAtZero(t *testing.T) {
	if New(2, 1000).Keys("p1") != 0 {
		t.Fatal("钥匙初始应该是 0")
	}
}

func TestAddKeysThenRead(t *testing.T) {
	world := New(2, 1000)
	world.AddKeys("p1", 3)
	if world.Keys("p1") != 3 {
		t.Fatal("发钥匙之后应该能读回来")
	}
}

func TestOpensStartAtZero(t *testing.T) {
	if New(2, 1000).Opens("c1", 0) != 0 {
		t.Fatal("开启次数初始应该是 0")
	}
}

func TestMaxOpensStored(t *testing.T) {
	if New(3, 1000).maxOpens != 3 {
		t.Fatal("次数上限应该能读回来")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New(2, 1000).Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestOpenConsumesKey(t *testing.T) {
	world := New(2, 100000)
	world.AddKeys("p1", 1)
	if !world.Open("p1", "c1", 0, "r1") {
		t.Fatal("有钥匙时应该开成功")
	}
	if world.Keys("p1") != 0 {
		t.Fatal("每成功开一次应该消耗一把钥匙")
	}
}

func TestOpenWithoutKeysFailsAndKeepsState(t *testing.T) {
	world := New(2, 100000)
	if world.Open("p1", "c1", 0, "r1") {
		t.Fatal("没有钥匙时应该返回 false")
	}
	if world.Keys("p1") != 0 || world.Opens("c1", 0) != 0 {
		t.Fatal("没有钥匙时不应该改变任何状态")
	}
}

func TestMaxOpensCap(t *testing.T) {
	world := New(2, 100000)
	world.AddKeys("p1", 5)
	if !world.Open("p1", "c1", 0, "r1") || !world.Open("p1", "c1", 0, "r2") {
		t.Fatal("前两次应该开成功")
	}
	if world.Open("p1", "c1", 0, "r3") {
		t.Fatal("超过 maxOpens 应该返回 false")
	}
	if world.Opens("c1", 0) != 2 {
		t.Fatal("已开次数不应该超过 maxOpens")
	}
	if world.Keys("p1") != 3 {
		t.Fatal("被上限拦下的开箱不应该扣钥匙")
	}
}

func TestPeriodRefresh(t *testing.T) {
	world := New(1, 1000)
	world.AddKeys("p1", 5)
	if !world.Open("p1", "c1", 0, "r1") {
		t.Fatal("第一次应该开成功")
	}
	if world.Open("p1", "c1", 900, "r2") {
		t.Fatal("refreshMs 未满时达到上限应该失败")
	}
	if !world.Open("p1", "c1", 1500, "r3") {
		t.Fatal("距周期起点满 refreshMs 后应该能重新开")
	}
	if world.Opens("c1", 1500) != 1 {
		t.Fatal("周期刷新后已开次数应该清零重计")
	}
}

func TestIdempotentReqID(t *testing.T) {
	world := New(5, 100000)
	world.AddKeys("p1", 5)
	if !world.Open("p1", "c1", 0, "same") {
		t.Fatal("第一次应该开成功")
	}
	if world.Open("p1", "c1", 0, "same") {
		t.Fatal("同一个 reqID 重复调用应该返回 false")
	}
	if world.Keys("p1") != 4 {
		t.Fatal("重复请求不应该再扣钥匙")
	}
	if world.Opens("c1", 0) != 1 {
		t.Fatal("重复请求不应该增加已开次数")
	}
}

func TestSameSequenceRepeatsSameResult(t *testing.T) {
	run := func() []bool {
		world := New(2, 1000)
		world.AddKeys("p1", 3)
		return []bool{
			world.Open("p1", "c1", 0, "r1"),
			world.Open("p1", "c1", 500, "r2"),
			world.Open("p1", "c1", 900, "r3"),
			world.Open("p1", "c1", 1500, "r4"),
		}
	}
	first, second := run(), run()
	for i := range first {
		if first[i] != second[i] {
			t.Fatal("同一串开启序列重复执行结果应该相同")
		}
	}
}

func TestScannedNotLinearInChestCount(t *testing.T) {
	world := New(5, 100000)
	world.AddKeys("p1", 100000)
	for i := 0; i < 100000; i++ {
		world.Open("p1", "c-"+strconv.Itoa(i), 0, "r-"+strconv.Itoa(i))
	}
	// 10 万宝箱台账下查询一次只扫描一次，不随宝箱数增长。
	before := world.Scanned()
	world.Opens("c-99999", 0)
	if world.Scanned()-before != 1 {
		t.Fatal("单次 Opens 的扫描量应该是 1")
	}
}
