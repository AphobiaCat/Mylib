package exec_manager

import(
	"fmt"
	"mylib/src/public"
	"os/exec"
	"io"
	"bufio"
)

const max_retry = 5

func Exec(retry bool, cmd ...string) (chan bool, bool){

	stop_chan := make(chan bool)

	if len(cmd) < 1{
		public.DBG_ERR("null cmd.")
		return stop_chan, false
	}

	app			:= exec.Command(cmd[0], cmd[1:]...)
	succ		:= true
	retry_times	:= 0

	var stdout io.ReadCloser
	var err error

	for ; retry; {
		
		stdout, err = app.StdoutPipe()
		if err != nil {
			public.DBG_ERR("get standard output failed. err:", err)

			if retry && retry_times < max_retry{
				public.Sleep(200)
				retry_times++
				continue
			}else{
				succ = false
				break
			}
		}

		if err = app.Start(); err != nil {
			public.DBG_ERR("start app failed. err:", err)
			if retry && retry_times < max_retry{
				public.Sleep(200)
				retry_times++
				continue
			}else{
				succ = false
				break
			}
		}

		break
	}

	public.DBG_LOG("here")

	if !succ{
		return stop_chan, succ
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	public.DBG_LOG("here")

	go func(){
		stop := <- stop_chan
		public.DBG_LOG("stop[", cmd[0], "], ", stop)
		app.Process.Kill()
	}()

	return stop_chan, succ
}



