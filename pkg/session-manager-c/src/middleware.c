#include "middleware.h"
#include <stdio.h>

bool authenticateRequest(const HttpRequest* request) {
    // TODO: Implementar autenticação
    return false; // Retornar true se autenticado, false caso contrário
}

void logRequest(const HttpRequest* request) {
    // TODO: Implementar logging
    printf("Log: %s %s\n", request->method, request->path);
}
