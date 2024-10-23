#include "Thread_Manager.h"
#include "Pub_Func.h"



Thread_Manager *thread_manager;
BOOL			finish_init	= FALSE;

void Send(uint32 thread_id, void* msg)
{
	if(thread_id >= thread_manager->m_thread_num)
	{	
		DBG_ERR("Send msg to thread[%d] overflow thread num[%d]", thread_id, thread_manager->m_thread_num);
		return ;
	}
	
	Thread_Info *thread_info = &thread_manager->m_threads_info[thread_id];

	if(thread_info->msg_queue_len == 0)
	{
		DBG_ERR("this thread[%d] have not queue", thread_id);
		return ;
	}

	SPIN_GET(&(thread_info->msg_set_or_get_queue_lock));

	//wait queue process
	while(TRUE)
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

void* Recv(uint32 thread_id)
{
	void * ret_msg = NULL;
	
	if(thread_id >= thread_manager->m_thread_num)
	{	
		DBG_ERR("Recv msg from thread[%d] overflow thread num[%d]", thread_id, thread_manager->m_thread_num);
		return ret_msg;
	}
	
	Thread_Info *thread_info = &thread_manager->m_threads_info[thread_id];

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
	finish_init = TRUE;

	thread_manager = (Thread_Manager*)malloc(sizeof(Thread_Manager));

	thread_manager->m_thread_num		= thread_num;
	thread_manager->m_threads_info	= (Thread_Info*)malloc(sizeof(Thread_Info) * thread_num);

	for(int i = 0; i < thread_num; i++)
	{
	    int ret;
		pthread_t thread_id;

		thread_manager->m_threads_info[i].thread_id		= thread_id;
		thread_manager->m_threads_info[i].thread_func	= thread_info[i].thread_func;

		thread_manager->m_threads_info[i].thread_info.thread_func	= thread_info[i].thread_func;
		thread_manager->m_threads_info[i].thread_info.init_data 	= thread_info[i].init_data;
		thread_manager->m_threads_info[i].thread_info.msg_queue_len	= thread_info[i].msg_queue_len;
		
		thread_manager->m_threads_info[i].thread_args.data	= thread_info[i].init_data;

		if(thread_info[i].msg_queue_len != 0)
		{

			thread_manager->m_threads_info[i].msg_queue				= (void**)malloc(sizeof(void*) * thread_info[i].msg_queue_len);
			thread_manager->m_threads_info[i].msg_queue_len			= thread_info[i].msg_queue_len;
			thread_manager->m_threads_info[i].msg_queue_now_index	= 0;
			thread_manager->m_threads_info[i].msg_queue_enter_index	= 0;
			
			SEM_Init(&(thread_manager->m_threads_info[i].msg_queue_get_data_lock));

			SPIN_Init(&(thread_manager->m_threads_info[i].msg_set_or_get_queue_lock));
			SPIN_RELEASE(&(thread_manager->m_threads_info[i].msg_set_or_get_queue_lock));
			
			thread_manager->m_threads_info[i].thread_args.msg_sem = &(thread_manager->m_threads_info[i].msg_queue_get_data_lock);
		}

		ret = pthread_create(&thread_id, NULL, thread_info[i].thread_func, &(thread_manager->m_threads_info[i].thread_args));
		if (ret != 0) {
	        DBG_ERR("create thread error");
	        assert(0);
	    }
	}	
}

void Close_All_Thread()
{
	if(finish_init){
		free(thread_manager);
		free(thread_manager->m_threads_info);
	}else{
		DBG_ERR("thread no init");
	}
}

void Send_Msg_To_Thread(int thread_id, void* msg)
{
	if(finish_init){
		Send((uint32)thread_id, msg);
	}else{
		DBG_ERR("thread no init");
	}
}

void* Recv_Msg_From_Thread(int thread_id)
{
	if(finish_init){
		return Recv((uint32)thread_id);
	}else{
		DBG_ERR("thread no init");
		return NULL;
	}
}


