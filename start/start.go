package start

import (
	"io"
	SYS "syscall"

	log "github.com/cihub/seelog"
	"github.com/grindlemire/GoSentry/watch"
	"github.com/vrecan/death"
)

// Run Entry Point for GoSentry
func Run() {

	log.Info("Starting GoSentry Service")

	var goRoutines []io.Closer
	death := death.NewDeath(SYS.SIGINT, SYS.SIGTERM) //pass the signals you want to end your application

	watch := watch.NewWatch()
	watch.Start()
	goRoutines = append(goRoutines, watch) // this will work as long as the type implements a Close method

	//when you want to block for shutdown signals
	death.WaitForDeath(goRoutines...) // this will finish when a signal of your type is sent to your application

}
