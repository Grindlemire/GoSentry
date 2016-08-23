package start

import (
	"io"
	"os"
	SYS "syscall"

	log "github.com/cihub/seelog"
	CONF "github.com/grindlemire/GoSentry/c"
	"github.com/grindlemire/GoSentry/watch"
	"github.com/vrecan/death"
)

// Run Entry Point for GoSentry
func Run() {

	c, err := CONF.GetConf("c.yml")
	if nil != err {
		log.Critical(err)
		os.Exit(2)
	}

	log.Info("Starting GoSentry Service")

	var goRoutines []io.Closer
	death := death.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	watch, err := watch.NewWatch(c)
	if err != nil {
		log.Error("Unable to create Watcher: ", err)
		return
	}
	watch.Start()
	goRoutines = append(goRoutines, watch)

	death.WaitForDeath(goRoutines...)

}
