#ifndef Web_Server_H
#define Web_Server_H

#include "esp_err.h"
#include "esp_http_server.h"

typedef esp_err_t (*router_proceser)(httpd_req_t *);

typedef struct
{
	router_proceser	*procesers;
	const char		*url;
	BOOL			is_get;		//or post
}Web_Server_Info;

BOOL Init_WebServer(Web_Server_Info *web_server_infos, int router_count);

#endif
