package example


import(
	"mylib/src/public"
	mq "mylib/src/module/msg_queue_manager"
)

func example_producer(){
	send_msg, close_msg, succ := mq.New_Msg_Queue_Producer("hello")

	if succ{
		public.DBG_LOG("hello_queue start")
	}

	for i:= 0; i < 10; i++{
		send_msg <- "hello " + public.ConvertNumToStr(int64(i))
	}

	public.Sleep(10000)

	close_msg <- true
}

func example_consumer(){
	stop1, succ1 := mq.New_Msg_Queue_Consumer("hello", "ch_1", func(recv string){
		public.DBG_LOG("recv msg:", recv)
	})

	stop2, succ2 := mq.New_Msg_Queue_Consumer("hello", "ch_2", func(recv string){
		public.DBG_LOG("recv msg:", recv)
	})

	if succ1 && succ2{
		public.DBG_LOG("succ connect")
	}

	public.Sleep(10000)
	stop1<- true
	stop2<- true

	public.DBG_LOG("close connect")
}

func Example_msg_queue(){
	go example_producer()

	example_consumer()
}

