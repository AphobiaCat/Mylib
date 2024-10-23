#include <stdio.h>
#include "Thread_Manager.h"
#include "Pub_Func.h"

void* test(void* args)
{
	Thread_Args *ta = (Thread_Args*)args;

	int *id = (int*)(ta->data);

	while(true)
	{
		void* msg = Recv_Msg_From_Thread(*id);
		DBG_LOG("thread[%d] recv msg[%d]", *id, *(int*)msg);
	}
	
	return NULL;
}


int main()
{

	int a = 0;
	int b = 1;
	Create_Thread_Info tmp[] = {\
		{test, &a, 10},\
		{test, &b, 10}\
	};

	Init_Thread_Manager(tmp, sizeof(tmp) / sizeof(Create_Thread_Info));

	while(true)
	{	
		Send_Msg_To_Thread(0, &a);
		Send_Msg_To_Thread(1, &b);
	}


	SLEEP(5);

	Close_All_Thread();
	return 0;
}
