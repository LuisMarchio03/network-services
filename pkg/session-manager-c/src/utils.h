//
// Created by luis_ on 09/10/2024.
//

#ifndef UTILS_H
#define UTILS_H

#include <bson/bson.h>

void convertTimeToString(time_t rawtime, char* buffer, size_t size);
void sessionToString(const Session* session, char* output, size_t output_size);
void deserializeSession(const char* message, Session* session);

#endif //UTILS_H


