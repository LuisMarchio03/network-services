#include "mongo.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <mongoc/mongoc.h>
#include <bson/bson.h>
#include <utils.h>

static mongoc_client_t *client = NULL; // Cliente MongoDB
static mongoc_collection_t *collection = NULL; // Coleção MongoDB

mongoc_client_t* connectMongoDB() {
    // Inicializar a biblioteca libmongoc
    mongoc_init();

    // Criar um cliente MongoDB
    const char *uri_string = "mongodb://localhost:27017"; // Substitua pela sua URI
    mongoc_uri_t *uri = mongoc_uri_new_with_error(uri_string, NULL);

    if (!uri) {
        fprintf(stderr, "Erro ao criar a URI do MongoDB: %s\n", uri_string);
        exit(EXIT_FAILURE);
    }

    client = mongoc_client_new_from_uri(uri);
    if (!client) {
        fprintf(stderr, "Erro ao criar o cliente MongoDB.\n");
        mongoc_uri_destroy(uri);
        exit(EXIT_FAILURE);
    }

    // Conectar-se à coleção
    collection = mongoc_client_get_collection(client, "mydb", "sessions"); // Nome do banco e coleção

    if (!collection) {
        fprintf(stderr, "Erro ao acessar a coleção.\n");
        mongoc_client_destroy(client);
        mongoc_uri_destroy(uri);
        exit(EXIT_FAILURE);
    }

    printf("Conexão com o MongoDB estabelecida com sucesso.\n");
    mongoc_uri_destroy(uri);

    return client;
}

void insertSession(const Session* session) {
    bson_t *document = NULL;
    bson_error_t error;
    bool success = NULL;

    if (!client) {
        fprintf(stderr, "Falha ao criar o cliente MongoDB.\n");
        return;
    }

    // Criar o documento BSON para a sessão
    document = bson_new();

    char created_at_str[20]; // Alocando buffer para a string
    char updated_at_str[20];
    convertTimeToString(session->created_at, created_at_str, sizeof(created_at_str));
    convertTimeToString(session->updated_at, updated_at_str, sizeof(updated_at_str));


    BSON_APPEND_UTF8(document, "session_id", session->session_id);
    BSON_APPEND_UTF8(document, "user_id", session->user_id);
    BSON_APPEND_UTF8(document, "client_ip", session->client_ip);
    BSON_APPEND_UTF8(document, "client_mac", session->client_mac);
    BSON_APPEND_UTF8(document, "dhcp_ip", session->dhcp_ip);
    BSON_APPEND_UTF8(document, "dhcp_subnet_mask", session->dhcp_subnet_mask);
    BSON_APPEND_UTF8(document, "token", session->token);
    BSON_APPEND_INT32(document, "expires_in", session->expires_in);
    BSON_APPEND_UTF8(document, "created_at", created_at_str);
    BSON_APPEND_UTF8(document, "updated_at", updated_at_str);

    // Inserir o documento na coleção
    success = mongoc_collection_insert_one(collection, document, NULL, NULL, &error);

    if (!success) {
        fprintf(stderr, "Erro ao inserir documento: %s\n", error.message);
    } else {
        printf("Sessão inserida com sucesso!\n");
    }

    // Limpar
    bson_destroy(document);
    mongoc_collection_destroy(collection);
    mongoc_client_destroy(client);
    mongoc_cleanup();
}

Session retrieveSession(const char* sessionId) {
    //TODO: Implementar a função retrieveSession
}

void updateSessionInMongo(const Session* session) {
    //TODO: Implementar a função updateSessionInMongo
}

void deleteSessionFromMongo(const char* sessionId) {
    bson_t *filter = NULL;
    bson_error_t error;
    bool success = NULL;

    if (!client) {
        fprintf(stderr, "Falha ao deletar o cliente do MongoDB.\n");
        return;
    }

    // Delete by session_id
    filter = bson_new();
    BSON_APPEND_UTF8(filter, "session_id", sessionId);

    // Inserir o documento na coleção
    success = mongoc_collection_delete_one(collection, filter, NULL, NULL, &error);
    if (!success) {
        fprintf(stderr, "Erro ao inserir documento: %s\n", error.message);
    } else {
        printf("Sessão inserida com sucesso!\n");
    }

    // Limpar
    bson_destroy(filter);
}

void disconnectMongoDB() {
    // Fechar a conexão com o MongoDB e liberar os recursos
    if (collection) {
        mongoc_collection_destroy(collection);
    }
    if (client) {
        mongoc_client_destroy(client);
    }
    mongoc_cleanup();
}
