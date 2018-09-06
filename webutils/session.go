package webutils

import (
	"github.com/astaxie/beego/session"
)

var (
	MemorySession       *session.Manager
	memorySessionConfig = &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "/tmp",
	}
)

func MemorySessionInit(sessionName string, expire int) {
	memorySessionConfig.CookieName = sessionName
	memorySessionConfig.CookieLifeTime = expire
	MemorySession, _ = session.NewManager("memory", memorySessionConfig)
	go MemorySession.GC()
}
