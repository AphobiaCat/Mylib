#ifndef Pub_Func_H
#define Pub_Func_H

#include <stdio.h>
#include <assert.h>
#include <unistd.h>
#include <time.h>


#define DEBUG_MODE

#ifdef DEBUG_MODE

#define DBG_LOG(str, ...)	\
	do{						\
		printf("INFO FUNC[%20s] | LINE[%5d] | LOG[" str "] \n\r", __FUNCTION__, __LINE__,##__VA_ARGS__);	\
	}while(0);

#define DBG_ERR(str, ...)	\
	do{						\
		printf("ERR FUNC[%20s] | LINE[%5d] | LOG[" str "] \n\r", __FUNCTION__, __LINE__,##__VA_ARGS__);	\
	}while(0);

#else

#define DBG_LOG(str, ...);
#define DBG_ERR(str, ...);

#endif


#define SLEEP(x) sleep((x));
#define USLEEP(x) usleep((x));

typedef unsigned long int	uint64;
typedef unsigned int		uint32;
typedef unsigned short		uint16;
typedef unsigned char		uint8;

#endif
