#include <stdio.h>
#include "public.h"

void Memory_Copy(const void* source, void* target, int byte_len)
{
	uint8* tmp_s = (uint8*)source;
	uint8* tmp_t = (uint8*)target;
	
	for (int i = 0; i < byte_len; i++)
	{
		tmp_t[i] = tmp_s[i];
	}
}

int Str_Len(const char *str)
{
	for (int i = 0; ; i++)
	{
		if (str[i] == '\0')
		{
			return i;
		}
	}
	return -1;
}

