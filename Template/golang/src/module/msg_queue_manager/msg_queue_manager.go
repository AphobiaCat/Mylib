package msg_queue_manager

import (
	"mylib/src/public"
	cache "mylib/src/module/cachesql_manager"
	"github.com/nsqio/go-nsq"
)

var nsqd_ip			string
var nsqlookupd_ip	string
var max_retry		int
var retry_time		int
var redis_key		string
var error_save_time	int

func New_Msg_Queue_Producer(topic string)(new_msg_chan chan string, close_queue chan bool, succ bool){
	succ = true

	config := nsq.NewConfig()

	// new producer, and connect to nsqd
	producer, err := nsq.NewProducer(nsqd_ip, config)
	if err != nil {
		public.DBG_ERR("new producer err:", err)
		succ = false
		return
	}

	new_msg_chan = make(chan string, 10)
	close_queue = make(chan bool, 2)

	go func(){
		defer producer.Stop()

		for{
			select {
				case <-close_queue:
					public.DBG_LOG("nsq topic[", topic, "] quit")
					return 

				case msg := <-new_msg_chan:

					public.DBG_LOG("send msg:", msg)
				
					retry_i := 0
					for ; retry_i < max_retry; retry_i++{

						err = producer.Publish(topic, []byte(msg))
						if err != nil {
							public.DBG_ERR("Publish error:", err)
							public.Sleep(retry_time)
							continue
						}
						break
					}

					if retry_i == max_retry && len(redis_key) != 0 {
						cache.Set_Cache(redis_key + ":" + topic + ":" + public.ConvertNumToStr(int64(public.Rand_U64())), msg, int64(error_save_time))
					}
			}
		}
	}()

	return
}

func New_Msg_Queue_Consumer(topic string, channel string, call_back func(string))(stop_chan chan bool, succ bool){
	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		public.DBG_ERR("new consumer err:", err)
		return
	}
	// connect to nsqlookupd or nsqd (suggest nsqlookupd first)

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		call_back(string(message.Body))
		return nil
	}))
	
	//err = consumer.ConnectToNSQD(nsqd_ip)
	err = consumer.ConnectToNSQLookupd(nsqlookupd_ip)
	if err != nil {
		public.DBG_ERR("connect err:", err)
		return
	}else {
		public.DBG_LOG("consumer connected to lookupd OK, topic:", topic, " channel:", channel)
	}

	stop_chan = make(chan bool, 2)

	go func(){
		<-stop_chan
		consumer.Stop()
		public.DBG_LOG("exit consumer topic[", topic, "] channel[", channel, "]")
	}()

	succ = true

	return 
}


func init(){
	nsqd_ip			= public.Config.Nsq.NsqdIp
	nsqlookupd_ip	= public.Config.Nsq.NsqlookupdIp
	max_retry		= public.Config.Nsq.MaxRetry
	retry_time		= public.Config.Nsq.RetryTime
	redis_key		= public.Config.Nsq.ErrorRedisKey
	error_save_time	= public.Config.Nsq.ErrorSaveTime

	public.DBG_LOG("nsq config[", public.Config.Nsq, "]")
}

