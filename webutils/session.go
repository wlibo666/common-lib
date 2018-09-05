package webutils

import (
	"fmt"
	"github.com/astaxie/beego/session"
)

var (
	MemorySession     *session.Manager
	memorySessionConf = `{"cookieName":"%s", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "cookieLifeTime": %d, "providerConfig": ""}`
)

func MemorySessionInit(sessionName string, expire int) {
	config := fmt.Sprintf(memorySessionConf, sessionName, expire)
	MemorySession, _ = session.NewManager("memory", config)
	go MemorySession.GC()
}
