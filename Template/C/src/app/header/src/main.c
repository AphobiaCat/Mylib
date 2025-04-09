#include "thread_manager.h"
#include "public.h"


Create_Thread_Info thread_info[3];

void* test_thread(void *init_data)
{
	DBG_LOG("thread start");

	Thread_Args *init_id = (Thread_Args*)init_data;
	
	int *recv_data = (int*)Recv_Msg_From_Thread(init_id->thread_id);

	DBG_LOG("thread start id: %d, init data: %d, recv data: %d", init_id->thread_id, *(int*)init_id->data, *recv_data);

	return NULL;
}

int main()
{
	int t1, t2, t3, d1, d2, d3;

	t1 = 1;
	t2 = 2;
	t3 = 3;
	d1 = 4;
	d2 = 5;
	d3 = 6;
	

	thread_info[0].thread_func = test_thread;
	thread_info[0].msg_queue_len = 5;
	thread_info[0].init_data = &t1;

	thread_info[1].thread_func = test_thread;
	thread_info[1].msg_queue_len = 5;
	thread_info[1].init_data = &t2;

	thread_info[2].thread_func = test_thread;
	thread_info[2].msg_queue_len = 5;
	thread_info[2].init_data = &t3;
	
	Init_Thread_Manager(thread_info, sizeof(thread_info) / sizeof(Create_Thread_Info));

	Send_Msg_To_Thread(0, &d1);
	Send_Msg_To_Thread(1, &d2);
	Send_Msg_To_Thread(2, &d3);

	SLEEP(5);

	Close_All_Thread();
	
	return 0;
}


