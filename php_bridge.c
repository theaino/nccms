#include <sapi/embed/php_embed.h>
#include <stdlib.h>
#include "php_bridge.h"

static char *php_out_buf = NULL;
static size_t php_out_len = 0;

size_t buf_write(const char *str, size_t len) {
	php_out_buf = realloc(php_out_buf, php_out_len + len + 1);
	memcpy(php_out_buf + php_out_len, str, len);
	php_out_len += len;
	php_out_buf[php_out_len] = '\0';
	return len;
}

char *php_get_output() {
	return php_out_buf;
}

void php_reset_output() {
	free(php_out_buf);
	php_out_buf = NULL;
	php_out_len = 0;
}

void render_php(char *path) {
	PHP_EMBED_START_BLOCK(0, NULL)

	sapi_module.ub_write = buf_write;

	zend_file_handle file_handle;
	zend_stream_init_filename(&file_handle, path);

	if (php_execute_script(&file_handle) == FAILURE) {
		php_printf("Failed to execute PHP script.\n");
	}
	PHP_EMBED_END_BLOCK()
}
