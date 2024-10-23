#include "Web_Server.h"

esp_err_t example_router(httpd_req_t *req) {
    // 打开index.html文件
    FILE* f = fopen("/spiffs/index.html", "r");
    if (f == NULL) {
        httpd_resp_send_404(req);
        return ESP_FAIL;
    }

    char line[256];
    while (fgets(line, sizeof(line), f)) {
        httpd_resp_send_chunk(req, line, strlen(line));
    }

    httpd_resp_send_chunk(req, NULL, 0);  // 结束发送
    fclose(f);
    return ESP_OK;
}


BOOL Init_WebServer(Web_Server_Info *web_server_infos, int router_count)
{
	httpd_handle_t server = NULL;
    httpd_config_t config = HTTPD_DEFAULT_CONFIG();


	if (httpd_start(&server, &config) != ESP_OK) 
	{
		DBG_ERR("start http error");
		return FALSE;
    }


	for(int i = 0; i < router_count; i++)
	{	
		httpd_uri_t httpd_config = {
	        .uri       = web_server_infos[i].url,
	        .method    = HTTP_GET,
	        .handler   = web_server_infos[i].procesers,
	        .user_ctx  = NULL
		};

		if(!web_server_infos[i].is_get)
		{
			httpd_config.method = HTTP_POST;
		}
			
		httpd_register_uri_handler(server, &httpd_config);
	}
}


