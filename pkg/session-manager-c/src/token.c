#include "session.h"
#include <stdio.h>
#include <string.h>
#include <mongo.h>

char* generateToken(const char* sessionId) {
    // Gerar um token aleatório simples baseado no sessionId e um valor aleatório
    char* token = malloc(64); // alocar espaço para o token
    if (token == NULL) {
        return NULL;
    }
    snprintf(token, 64, "%s-%d", sessionId, rand());
    return token;
}

bool validateToken(const char* token) {
    // Verificar o token contra um valor esperado (apenas exemplo simplificado)
    if (strcmp(token, "expected_token") == 0) {
        return true;
    }
    return false;
}
