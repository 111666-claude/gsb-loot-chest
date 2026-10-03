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
