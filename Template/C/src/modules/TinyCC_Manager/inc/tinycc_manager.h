#ifndef TinyCC_Manager_H
#define TinyCC_Manager_H

#include <libtcc.h>

class TinyCC_Manager
{

public:
	TinyCC_Manager(const char *init_source_content);
	~TinyCC_Manager();
	void* Get_Func(const char *func_name);
	
private:
	char*		m_source_content;
	bool		m_have_init;
	TCCState*	m_tcc;
};

#endif
