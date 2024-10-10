#ifndef RABBITMQ_H
#define RABBITMQ_H

#include <stdbool.h>

// Funções para manipulação de mensagens no RabbitMQ
void connectRabbitMQ();
void publishMessage(const char* message);
void consumeMessages();
void disconnectRabbitMQ();

#endif // RABBITMQ_H
