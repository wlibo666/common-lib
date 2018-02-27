package atomic

import "testing"

func equal(t *testing.T, number1, number2 uint64) {
	if number1 != number2 {
		t.Fatalf("number1:%d,is not equal with number2:%d", number1, number2)
	}
}

func TestIncrDefUInt64(t *testing.T) {
	t.Logf("def cnt:%d", GetDefUInt64())
	equal(t, GetDefUInt64(), 0)

	IncrDefUInt64()
	t.Logf("def cnt:%d,should be 1", GetDefUInt64())
	equal(t, GetDefUInt64(), 1)

	IncrDefUInt64()
	t.Logf("def cnt:%d,should be 2", GetDefUInt64())
	equal(t, GetDefUInt64(), 2)

	DecrDefUInt64()
	t.Logf("def cnt:%d,should be 1", GetDefUInt64())
	equal(t, GetDefUInt64(), 1)

	ResetDefUInt64()
	t.Logf("def cnt:%d,should be 0", GetDefUInt64())
	equal(t, GetDefUInt64(), 0)

	var customNum uint64 = 100
	t.Logf("def cnt:%d,should be 100", GetUInt64(&customNum))
	equal(t, GetUInt64(&customNum), 100)

	IncrUInt64(&customNum)
	t.Logf("def cnt:%d,should be 101", GetUInt64(&customNum))
	equal(t, GetUInt64(&customNum), 101)

	IncrUInt64(&customNum)
	t.Logf("def cnt:%d,should be 102", GetUInt64(&customNum))
	equal(t, GetUInt64(&customNum), 102)

	DecrUInt64(&customNum)
	t.Logf("def cnt:%d,should be 101", GetUInt64(&customNum))
	equal(t, GetUInt64(&customNum), 101)

	ResetUInt64(&customNum)
	t.Logf("def cnt:%d,should be 0", GetUInt64(&customNum))
	equal(t, GetUInt64(&customNum), 0)
}
