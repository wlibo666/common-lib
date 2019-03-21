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
		Secure:          true,
		CookieLifeTime:  3600,
		ProviderConfig:  "/tmp",
	}
)

func MemorySessionInit(sessionName string, expire int, secure bool) {
	memorySessionConfig.CookieName = sessionName
	memorySessionConfig.CookieLifeTime = expire
	memorySessionConfig.Secure = secure
	MemorySession, _ = session.NewManager("memory", memorySessionConfig)
	go MemorySession.GC()
}
