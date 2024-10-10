#ifndef TOKEN_H
#define TOKEN_H

#include <stdbool.h>

// Funções para geração e validação de tokens
char* generateToken(const char* sessionId);
bool validateToken(const char* token);

#endif // TOKEN_H
