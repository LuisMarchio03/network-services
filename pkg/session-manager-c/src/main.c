#include "../include/session_manager.h"
#include <stdio.h>
#include <signal.h>
#include <mongo.h>
#include <rabbitmq.h>

// Variável para controlar o loop principal
volatile int keepRunning = 1;

void intHandler(int dummy) {
    keepRunning = 0;
}

int main() {
    // Captura sinal de interrupção para fechar o serviço corretamente
    signal(SIGINT, intHandler);

    // Inicializar conexões com RabbitMQ e MongoDB
    connectRabbitMQ();
    connectMongoDB();

    printf("Session Manager Service Running...\n");

    // Loop principal do serviço
    while (keepRunning) {
        consumeMessages();  // Consome mensagens do RabbitMQ e realiza as operações necessárias
    }

    // Fechar conexões antes de sair
    disconnectMongoDB();
    disconnectRabbitMQ();

    printf("Session Manager Service Stopped.\n");
    return 0;
}
