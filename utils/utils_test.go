package utils

import (
	"testing"
)

func TestJsonIndent(t *testing.T) {
	data := `{"errno":10000,"err":"","data":{"index":4,"hash":"0x232b8fb20767ef09bd36ac1cf16673506c379996250d2e631e92056e42d73674","parent":"0xc39c279d1dea2889f30c3ee2a002d0848052777969ffeb6e738c431e110665b3","received":1535100622,"size":775,"tx":1,"transactions":[{"hash":"0x2a758cfa94d14672fab74ecb30a38ae3b9545a3c7dcf678696b0f79e0d7a0eca"}]}}`
	b, err := JsonIndent([]byte(data))
	if err != nil {
		t.Fatalf("JsonIndent failed,err:%s", err.Error())
	}
	t.Logf("byte:\n%s\n", b)
}

func TestRemovePath(t *testing.T) {
	f := RemovePath("/root/work/WorkCode/src/git.smartisan.com/infrastructure/eththirdproxy/controller/business/mongo.go:167", 7)
	t.Logf("path:%s", f)
}

func TestGetRoutineId(t *testing.T) {
	id := GetRoutineId()

	t.Logf("routine:%s end", id)
}
