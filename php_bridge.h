#include <stdlib.h>

size_t buf_write(const char *str, size_t len);

char *php_get_output();

void php_reset_output();

void render_php(char *path);
