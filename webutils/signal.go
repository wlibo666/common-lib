package webutils

import (
	"os"
	"os/signal"

	"git.smartisan.com/infrastructure/golang-common/log"
	"github.com/wlibo666/common-lib/utils"
)

func SignalHandle() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		log.InfoWithFields(log.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"signal":           s.String(),
		}, "recv signal,will exit with code 0.")
		utils.ExitWaitDef(0)
	}()
	return nil
}
