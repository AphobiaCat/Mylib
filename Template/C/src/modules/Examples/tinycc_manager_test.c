#include "thread_manager.h"
#include "public.h"

#include "tinycc_manager.h"

int main()
{
	TinyCC_Manager *tmp_tinycc = new TinyCC_Manager("int add(int a, int b){return a + b;}");

	int (*add)(int, int) = (int (*)(int, int))tmp_tinycc->Get_Func("add");

	DBG_LOG("add: 1 + 2 = %d", add(1, 2));

	delete tmp_tinycc;
	
	return 0;
}

