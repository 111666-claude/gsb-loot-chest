package chest

import "testing"

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
	world := New(2, 1000)
	world.AddKeys("p1", 1)
	if !world.Open("p1", "c1", 0, "r1") {
		t.Fatal("有钥匙应该能开")
	}
	if world.Keys("p1") != 0 {
		t.Fatal("开一次应该扣一把钥匙")
	}
}

func TestOpenWithoutKeysFails(t *testing.T) {
	world := New(2, 1000)
	if world.Open("p1", "c1", 0, "r1") {
		t.Fatal("没有钥匙不该开成功")
	}
	if world.Keys("p1") != 0 || world.Opens("c1", 0) != 0 {
		t.Fatal("失败不该改变任何状态")
	}
}

func TestMaxOpensCap(t *testing.T) {
	world := New(2, 1000)
	world.AddKeys("p1", 5)
	world.Open("p1", "c1", 0, "r1")
	world.Open("p1", "c1", 0, "r2")
	if world.Open("p1", "c1", 0, "r3") {
		t.Fatal("超出次数上限不该开成功")
	}
	if world.Opens("c1", 0) != 2 {
		t.Fatal("已开次数不应超过上限")
	}
}

func TestRefreshWindow(t *testing.T) {
	world := New(1, 1000)
	world.AddKeys("p1", 5)
	world.Open("p1", "c1", 0, "r1")
	if world.Open("p1", "c1", 900, "r2") {
		t.Fatal("周期内第二次应该被拦")
	}
	if !world.Open("p1", "c1", 1500, "r3") {
		t.Fatal("刷新后应该能重新开")
	}
	if world.Opens("c1", 1500) != 1 {
		t.Fatal("刷新后已开次数应该清零重计")
	}
}

func TestIdempotentReqID(t *testing.T) {
	world := New(5, 1000)
	world.AddKeys("p1", 5)
	world.Open("p1", "c1", 0, "same")
	if world.Open("p1", "c1", 0, "same") {
		t.Fatal("同一个 reqID 不该开第二次")
	}
	if world.Keys("p1") != 4 || world.Opens("c1", 0) != 1 {
		t.Fatal("重复请求不该扣钥匙也不该加次数")
	}
}

func TestReplayDeterministic(t *testing.T) {
	run := func() []bool {
		world := New(2, 1000)
		world.AddKeys("p1", 3)
		return []bool{
			world.Open("p1", "c1", 0, "r1"),
			world.Open("p1", "c1", 100, "r2"),
			world.Open("p1", "c1", 200, "r3"),
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
