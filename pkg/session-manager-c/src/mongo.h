#ifndef MONGO_H
#define MONGO_H

#include "../include/session_manager.h"
#include <mongoc/mongoc.h>


mongoc_client_t* connectMongoDB();
void insertSession(const Session* session);
Session retrieveSession(const char* sessionId);
void updateSessionInMongo(const Session* session);
void deleteSessionFromMongo(const char* sessionId);
void disconnectMongoDB();

#endif // MONGO_H
