#include "session.h"
#include <stdio.h>
#include <string.h>
#include <mongo.h>
#include <bson/bson.h>
#include <utils.h>

#define MAX_STRING_SIZE 512

int createSession(const Session* data) {
    Session new_session;
    char new_session_string[MAX_STRING_SIZE];

    // Gerar um ID único para a sessão (exemplo simples de string ID)
    snprintf(new_session.session_id, sizeof(new_session.session_id), "session-%ld", time(NULL));

    // Copiar os dados da sessão fornecidos
    memcpy(&new_session, data, sizeof(SessionData));

    time_t now = time(NULL);
    struct tm *t = gmtime(&now);
    strftime(new_session.created_at, sizeof(new_session.created_at), "%Y-%m-%dT%H:%M:%SZ", t);
    strcpy(new_session.updated_at, new_session.created_at);  // Inicialmente, `updated_at` é o mesmo que `created_at`

    // Gerar Token
    strcpy(new_session.token, generateToken(new_session.session_id));

    // Salvar a sessão no MongoDB
    insertSession(&new_session);

    // Salvar a sessão no RabbitMQ
    sessionToString(&new_session, new_session_string, sizeof(new_session_string));
    publishMessage(new_session_string);

    return 0;
}


Session getSession(const char* sessionId) {
    Session session;
    // TODO: Recuperar sessão do MongoDB
    return session; // Retornar sessão recuperada
}

int updateSession(const char* sessionId, const SessionData* newData) {
    // TODO: Atualizar sessão no MongoDB
    return 0; // Retornar 0 se bem-sucedido, -1 se falhar
}

int deleteSession(const char* sessionId) {
    // TODO: Excluir sessão do MongoDB
    return 0; // Retornar 0 se bem-sucedido, -1 se falhar
}

void manageSessionExpiration() {
    // TODO: Gerenciar expiração de sessões
}
