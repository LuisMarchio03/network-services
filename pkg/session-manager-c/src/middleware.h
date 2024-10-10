#ifndef MIDDLEWARE_H
#define MIDDLEWARE_H

#include <stdbool.h>

// Estrutura para uma solicitação HTTP (exemplo simples)
typedef struct {
    char method[8]; // Método da requisição (GET, POST, etc.)
    char path[256]; // Caminho da requisição
    char body[1024]; // Corpo da requisição
} HttpRequest;

// Funções para middleware
bool authenticateRequest(const HttpRequest* request);
void logRequest(const HttpRequest* request);

#endif // MIDDLEWARE_H
