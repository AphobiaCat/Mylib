#ifndef Thread_Manager_H
#define Thread_Manager_H

#include "Pub_Func.h"

#include <pthread.h>
#include <semaphore.h>

#ifdef __linux__

#define SEM					sem_t
#define SEM_Init(s_ptr)		sem_init(s_ptr, 0, 0);
#define SEM_GET(s_ptr)		sem_wait(s_ptr);
#define SEM_TRY_GET(s_ptr)	sem_trywait(s_ptr);
#define SEM_RELEASE(s_ptr)	sem_post(s_ptr);

#define SPIN				pthread_spinlock_t
#define SPIN_Init(s_ptr)	pthread_spin_init(s_ptr, 0);
#define SPIN_GET(s_ptr)		pthread_spin_lock(s_ptr);
#define SPIN_TRYGET(s_ptr)	pthread_spin_trylock(s_ptr);
#define SPIN_RELEASE(s_ptr)	pthread_spin_unlock(s_ptr);

#else
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"

#define SEM					sem_t
#define SEM_Init(s_ptr)		sem_init(s_ptr, 0, 0);
#define SEM_GET(s_ptr)		sem_wait(s_ptr);
#define SEM_TRY_GET(s_ptr)	sem_trywait(s_ptr);
#define SEM_RELEASE(s_ptr)	sem_post(s_ptr);

#define SPIN				portMUX_TYPE
#define SPIN_Init(s_ptr)	do{(s_ptr)->owner = portMUX_FREE_VAL; (s_ptr)->count = 0; portENTER_CRITICAL(s_ptr);}while(0);
#define SPIN_GET(s_ptr)		portENTER_CRITICAL(s_ptr);
#define SPIN_TRYGET(s_ptr)	portTRY_ENTER_CRITICAL(s_ptr);
#define SPIN_RELEASE(s_ptr)	portEXIT_CRITICAL(s_ptr);

#endif

typedef void* (*Thread_Function)(void*);

typedef struct{
	void*			data;
	SEM*			msg_sem;
}Thread_Args;

typedef struct
{
	Thread_Function	thread_func;
	void *			init_data;
	uint32			msg_queue_len;
}Create_Thread_Info;



typedef struct
{
	Thread_Function		thread_func;
	Create_Thread_Info	thread_info;
	Thread_Args			thread_args;
	pthread_t			thread_id;
	
	void **				msg_queue;
	uint32				msg_queue_len;
	SEM					msg_queue_get_data_lock;
	SPIN				msg_set_or_get_queue_lock;
	uint32				msg_queue_now_index;
	uint32				msg_queue_enter_index;
}Thread_Info;

typedef struct
{
	Thread_Info		*m_threads_info;
	uint32			m_thread_num;
}Thread_Manager;

void Send(uint32 thread_id, void* msg);
void* Recv(uint32 thread_id);

void Init_Thread_Manager(Create_Thread_Info *thread_info, uint32 thread_num);
void Close_All_Thread();

void Send_Msg_To_Thread(int thread_id, void* msg);
void* Recv_Msg_From_Thread(int thread_id);

#endif

