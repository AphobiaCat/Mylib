#include "Thread_Manager.h"
#include "Pub_Func.h"



Thread_Manager *thread_manager;
bool			finish_init	= false;


Thread_Manager::Thread_Manager(Create_Thread_Info *thread_info, int thread_num)
{
	m_thread_num		= thread_num;
	m_threads_info	= new Thread_Info[thread_num];

	for(int i = 0; i < thread_num; i++)
	{
	    int ret;
		pthread_t thread_id;

		m_threads_info[i].thread_id		= thread_id;
		m_threads_info[i].thread_func	= thread_info[i].thread_func;
		m_threads_info[i].thread_info	= thread_info[i];
		
		m_threads_info[i].thread_args.data	= thread_info[i].init_data;

		if(thread_info[i].msg_queue_len != 0)
		{
			m_threads_info[i].msg_queue				= new void*[thread_info[i].msg_queue_len];
			m_threads_info[i].msg_queue_len			= thread_info[i].msg_queue_len;
			m_threads_info[i].msg_queue_now_index	= 0;
			m_threads_info[i].msg_queue_enter_index	= 0;
			
			SEM_Init(&(m_threads_info[i].msg_queue_get_data_lock));

			SPIN_Init(&(m_threads_info[i].msg_set_or_get_queue_lock));
			SPIN_RELEASE(&(m_threads_info[i].msg_set_or_get_queue_lock));
			
			m_threads_info[i].thread_args.msg_sem = &(m_threads_info[i].msg_queue_get_data_lock);
		}

		ret = pthread_create(&thread_id, NULL, thread_info[i].thread_func, &(m_threads_info[i].thread_args));
		if (ret != 0) {
	        DBG_ERR("create thread error");
	        assert(0);
	    }
	}	
}

Thread_Manager::~Thread_Manager()
{
	delete []m_threads_info;
}

void Thread_Manager::Send(uint32 thread_id, void* msg)
{
	if(thread_id >= m_thread_num)
	{	
		DBG_ERR("Send msg to thread[%d] overflow thread num[%d]", thread_id, m_thread_num);
		return ;
	}
	
	Thread_Info *thread_info = &m_threads_info[thread_id];

	if(thread_info->msg_queue_len == 0)
	{
		DBG_ERR("this thread[%d] have not queue", thread_id);
		return ;
	}

	SPIN_GET(&(thread_info->msg_set_or_get_queue_lock));

	//wait queue process
	while(true)
	{		
		if((thread_info->msg_queue_enter_index + 1) % thread_info->msg_queue_len == thread_info->msg_queue_now_index)
		{
			DBG_LOG("queue full");
			SPIN_RELEASE(&(thread_info->msg_set_or_get_queue_lock));
			USLEEP(15000);	//wait 15ms
			SPIN_GET(&(thread_info->msg_set_or_get_queue_lock));
		}
		else
		{
			break;
		}
	}

	thread_info->msg_queue[thread_info->msg_queue_enter_index] = msg;
	thread_info->msg_queue_enter_index = (thread_info->msg_queue_enter_index + 1) % thread_info->msg_queue_len;

	SEM_RELEASE(&(thread_info->msg_queue_get_data_lock)); 
	
	SPIN_RELEASE(&(thread_info->msg_set_or_get_queue_lock));
}

void* Thread_Manager::Recv(uint32 thread_id)
{
	void * ret_msg = NULL;
	
	if(thread_id >= m_thread_num)
	{	
		DBG_ERR("Recv msg from thread[%d] overflow thread num[%d]", thread_id, m_thread_num);
		return ret_msg;
	}
	
	Thread_Info *thread_info = &m_threads_info[thread_id];

	if(thread_info->msg_queue_len == 0)
	{
		DBG_ERR("this thread[%d] have not queue", thread_id);
		return ret_msg;
	}

	SEM_GET(&(thread_info->msg_queue_get_data_lock));

	SPIN_GET(&(thread_info->msg_set_or_get_queue_lock));

	ret_msg = thread_info->msg_queue[thread_info->msg_queue_now_index];
	thread_info->msg_queue_now_index = (thread_info->msg_queue_now_index + 1) % thread_info->msg_queue_len;

	SPIN_RELEASE(&(thread_info->msg_set_or_get_queue_lock));

	return ret_msg;
}

void Init_Thread_Manager(Create_Thread_Info *thread_info, uint32 thread_num)
{
	if(finish_init){
		DBG_ERR("already init thread");
		return;
	}
	finish_init = true;

	thread_manager = new Thread_Manager(thread_info, thread_num);
}

void Close_All_Thread()
{
	if(finish_init){
		delete thread_manager;
	}else{
		DBG_ERR("thread no init");
	}
}

void Send_Msg_To_Thread(int thread_id, void* msg)
{
	if(finish_init){
		thread_manager->Send((uint32)thread_id, msg);
	}else{
		DBG_ERR("thread no init");
	}
}

void* Recv_Msg_From_Thread(int thread_id)
{
	if(finish_init){
		return thread_manager->Recv((uint32)thread_id);
	}else{
		DBG_ERR("thread no init");
		return NULL;
	}
}


