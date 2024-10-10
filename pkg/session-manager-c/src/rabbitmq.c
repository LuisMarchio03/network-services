#include "rabbitmq.h"
#include <amqp.h>
#include <amqp_tcp_socket.h>
#include <stdio.h>
#include <stdlib.h>
#include "session.h"
#include "utils.h"

#define RABBITMQ_HOST "localhost"
#define RABBITMQ_PORT 5672
#define RABBITMQ_QUEUE_PUBLISH "session_created_queue"
#define RABBITMQ_QUEUE_CONSUME "dhcp_ack_queue"

amqp_connection_state_t conn;

void connectRabbitMQ() {
    conn = amqp_new_connection();
    amqp_socket_t *socket = amqp_tcp_socket_new(conn);

    if (!socket) {
        fprintf(stderr, "Erro ao criar socket\n");
        exit(1);
    }

    int status = amqp_socket_open(socket, RABBITMQ_HOST, RABBITMQ_PORT);
    if (status) {
        fprintf(stderr, "Erro ao abrir conexão TCP para o RabbitMQ\n");
        exit(1);
    }

    amqp_rpc_reply_t login = amqp_login(conn, "/", 0, 131072, 0, AMQP_SASL_METHOD_PLAIN, "guest", "guest");
    if (login.reply_type != AMQP_RESPONSE_NORMAL) {
        fprintf(stderr, "Erro ao autenticar no RabbitMQ\n");
        exit(1);
    }

    amqp_channel_open(conn, 1);
    amqp_rpc_reply_t channel_open_reply = amqp_get_rpc_reply(conn);
    if (channel_open_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        fprintf(stderr, "Erro ao abrir canal no RabbitMQ\n");
        exit(1);
    }

    printf("Conectado ao RabbitMQ\n");
}

void publishMessage(const char* message) {
    amqp_bytes_t message_bytes = amqp_cstring_bytes(message);
    amqp_basic_properties_t props;
    props._flags = AMQP_BASIC_CONTENT_TYPE_FLAG | AMQP_BASIC_DELIVERY_MODE_FLAG;
    props.content_type = amqp_cstring_bytes("text/plain");
    props.delivery_mode = 2; // Persistente

    int status = amqp_basic_publish(conn, 1, amqp_empty_bytes, amqp_cstring_bytes(RABBITMQ_QUEUE_PUBLISH), 0, 0, &props, message_bytes);
    if (status < 0) {
        fprintf(stderr, "Erro ao publicar mensagem no RabbitMQ\n");
        exit(1);
    }

    printf("Mensagem publicada: %s\n", message);
}

void consumeMessages() {
    amqp_basic_consume(conn, 1, amqp_cstring_bytes(RABBITMQ_QUEUE_CONSUME), amqp_empty_bytes, 0, 1, 0, amqp_empty_table);
    amqp_rpc_reply_t consume_reply = amqp_get_rpc_reply(conn);
    if (consume_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        fprintf(stderr, "Erro ao iniciar consumo de mensagens\n");
        exit(1);
    }

    while (1) {
        amqp_envelope_t envelope;
        amqp_maybe_release_buffers(conn);

        amqp_rpc_reply_t recv_reply = amqp_consume_message(conn, &envelope, NULL, 0);
        if (recv_reply.reply_type != AMQP_RESPONSE_NORMAL) {
            fprintf(stderr, "Erro ao consumir mensagem\n");
            break;
        }

        printf("Mensagem recebida: %.*s\n", (int)envelope.message.body.len, (char *)envelope.message.body.bytes);

        // Construir Session a partir da mensagem
        Session session;
        deserializeSession(envelope.message.body.bytes, &session);

        // Criar uma nova sessão com base na mensagem
        if (createSession(&session) != 0) {
            fprintf(stderr, "Erro ao criar sessão a partir da mensagem recebida\n");
        }

        amqp_destroy_envelope(&envelope);
    }
}
