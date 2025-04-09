#ifndef PUBLIC_H
#define PUBLIC_H

#include <stdio.h>
#include <assert.h>
#include <unistd.h>
#include <time.h>


#define DBG_LOG(str, ...)	\
	do{ 					\
		printf("INFO FUNC[%20s] | LINE[%4d] | LOG[" str "] \n\r", __FUNCTION__, __LINE__,##__VA_ARGS__); \
	}while(0);
	
#define DBG_ERR(str, ...)	\
	do{ 					\
		printf("\033[1;31mERR  FUNC[%20s] | LINE[%4d] | LOG[" str "] \033[0m\n\r", __FUNCTION__, __LINE__,##__VA_ARGS__); \
	}while(0);

#define SLEEP(x) sleep((x));
#define USLEEP(x) usleep((x));

typedef unsigned long int	uint64;
typedef unsigned int		uint32;
typedef unsigned short		uint16;
typedef unsigned char		uint8;
typedef enum{FALSE=0, TRUE=1} BOOL;

void Memory_Copy(const void* source, void* target, int len);
int Str_Len(const char *str);


#endif

