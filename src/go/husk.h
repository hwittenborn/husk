#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

typedef uintptr_t ptr;

#ifndef HUSK_RUST_READ
#define HUSK_RUST_READ
struct HuskRustRead {
	bool isOk;
	uint8_t byte;
	uint8_t bytesRead;
	void* errObj;
};
#endif

extern int HUSK_ERROR_IO;
extern int HUSK_ERROR_UNEXPECTED_COMMAND;
extern int HUSK_ERROR_UNSET_PARAMETER;
extern int HUSK_ERROR_EXIT_STATUS;
extern int HUSK_ERROR_LANG;
extern int HUSK_ERROR_PARSE;
extern int HUSK_ERROR_QUOTE;
extern int HUSK_ERROR_UNKNOWN;

extern int HUSK_SHELL_EXPAND_STRING;
extern int HUSK_SHELL_FIELDS_STRINGS;
extern int HUSK_SYNTAX_QUOTE_STRING;
extern int HUSK_SYNTAX_PARSER_FILE;

extern int HUSK_LANG_BASH;
extern int HUSK_LANG_POSIX;
extern int HUSK_LANG_MKSH;
extern int HUSK_LANG_BATS;
extern int HUSK_LANG_AUTO;

char* HuskRustGetHashMapKey(void* hmap, unsigned int index);
char* HuskRustGetHashMapValue(void* hmap, char* key);
struct HuskRustRead HuskRustGetByteFromReadTrait(void* container);
