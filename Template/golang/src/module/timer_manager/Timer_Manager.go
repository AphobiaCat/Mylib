package timer_manager

import (
	"time"
	"mylib/src/public"
)

type Timer_Callback func()

const max_timer_cons int = 10

var timer_manager Timer_Manager

type Timer_Manager struct{
	max_timer_index		int
	now_timer_index		int
	callbacks			[max_timer_cons]Timer_Callback
	callback_max_time	[max_timer_cons]uint32
	callback_now_time	[max_timer_cons]uint32

	init_once			bool
	ticker				*time.Ticker
	now_time			uint32
	timer_ongoing		bool
}

func (tm *Timer_Manager) Init(){

	if tm.init_once == false{
		tm.max_timer_index = max_timer_cons
		tm.now_timer_index = 0
	}	
	tm.init_once = true

	tm.ticker = time.NewTicker(1 * time.Second)
	tm.now_time = 0

	tm.timer_ongoing = true

    go func(){
		defer tm.ticker.Stop()
		for tm.timer_ongoing {
			select{
				case _ = <-tm.ticker.C:
					tm.now_time++
			}

			for i := 0; i < tm.now_timer_index; i++{
				tm.callback_now_time[i]++
				if tm.callback_now_time[i] >= tm.callback_max_time[i]{
					tm.callback_now_time[i] = 0
					go tm.callbacks[i]()
				}	
			}
		}
    }()
}

func Stop(){
    timer_manager.timer_ongoing = false
    timer_manager.now_timer_index = 0
}


func Reg_Timer(call_time_s uint32, callback Timer_Callback) int{
	if timer_manager.now_timer_index >= timer_manager.max_timer_index{
		public.DBG_LOG("should extented timer")
		return -1
	}

	timer_manager.callback_now_time[timer_manager.now_timer_index]	= 0
	timer_manager.callback_max_time[timer_manager.now_timer_index]	= call_time_s
	timer_manager.callbacks[timer_manager.now_timer_index]		 	= callback

	timer_manager.now_timer_index++

	return timer_manager.now_timer_index - 1
}

func Change_Timer_Set(timer_id int, call_time_s uint32) bool{
	if timer_id >= timer_manager.max_timer_index{
		public.DBG_LOG("error timer id")
		return false
	}

	timer_manager.callback_now_time[timer_id] = 0
	timer_manager.callback_max_time[timer_id] = call_time_s
	
	return true
}

func init(){
	timer_manager.Init()
}


