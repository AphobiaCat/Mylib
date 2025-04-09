#include "tinycc_manager.h"
#include "public.h"

#include <libtcc.h>


TinyCC_Manager::TinyCC_Manager(const char *init_source_content)
{
	this->m_have_init = false;

	this->m_tcc	= tcc_new();
    if (!this->m_tcc) {
		DBG_ERR("Could not create TCC state")
        return ;
    }
	
	
	int copy_len = Str_Len(init_source_content) + 1;
	this->m_source_content = new char[copy_len];
	
	Memory_Copy(init_source_content, this->m_source_content, copy_len);

	tcc_set_output_type(this->m_tcc, TCC_OUTPUT_MEMORY);

	if (tcc_compile_string(this->m_tcc, init_source_content) == -1)
	{
		DBG_ERR("Compilation failed");		
		delete []this->m_source_content;
		tcc_delete(this->m_tcc);
		return ;
    }

	if (tcc_relocate(this->m_tcc) < 0) {
		DBG_ERR("Relocation failed");
		delete []this->m_source_content;
		tcc_delete(this->m_tcc);
		return ;
	}
	
	this->m_have_init = true;
}

TinyCC_Manager::~TinyCC_Manager()
{
	if(this->m_have_init == true)
	{
		delete []this->m_source_content;

		tcc_delete(this->m_tcc);

		this->m_have_init = false;	
	}
}

void* TinyCC_Manager::Get_Func(const char *func_name)
{
	if(this->m_have_init == true)
	{
		void* ret_func = tcc_get_symbol(this->m_tcc, func_name);

	    if (!ret_func) {
	        DBG_ERR("Function[%s] not found", func_name);
	        return NULL;
	    }
		return ret_func;
	}

	return NULL;
}

