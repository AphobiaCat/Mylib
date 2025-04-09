#ifndef test_H
#define test_H

class Test
{
public:
	static void Init();

	~Test();
	
private:
	static Test *g_test;	
};

#endif

