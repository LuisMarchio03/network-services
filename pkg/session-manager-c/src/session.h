#ifndef SESSION_H
#define SESSION_H

#include <time.h>
#include "../include/session_manager.h"

#define SESSION_ID_LENGTH 64 // Comprimento do ID da sessão


// Funções para gerenciamento de sessões
int createSession(const SessionData* data);
Session getSession(const char* sessionId);
int updateSession(const char* sessionId, const SessionData* newData);
int deleteSession(const char* sessionId);
void manageSessionExpiration();

#endif // SESSION_H
