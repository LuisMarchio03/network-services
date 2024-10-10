#ifndef SESSION_MANAGER_H
#define SESSION_MANAGER_H

#include <stdbool.h>
#include <time.h>
#include <mongoc/mongoc.h>

// Estrutura que representa uma requisição HTTP
typedef struct {
    char method[8];       // Método HTTP (GET, POST, etc.)
    char url[256];        // URL da requisição
    char headers[512];    // Cabeçalhos da requisição
    char body[1024];      // Corpo da requisição
} HttpRequest;

// Estrutura que representa os dados de uma sessão
typedef struct {
    char session_id[64];          // Identificador único da sessão
    char user_id[64];             // Identificador do usuário
    char client_ip[46];           // Endereço IP do cliente (compatível com IPv4 e IPv6)
    char client_mac[18];          // Endereço MAC do cliente
    char dhcp_ip[46];             // Endereço IP atribuído pelo servidor DHCP
    char dhcp_subnet_mask[46];    // Máscara de sub-rede atribuída pelo servidor DHCP
    char token[256];              // Token da sessão (JWT ou outro)
    int expires_in;               // Tempo de expiração da sessão (em segundos)
    time_t created_at;            // Hora de criação da sessão
    time_t updated_at;            // Hora da última atualização da sessão
} Session;

// Estrutura para os dados da sessão, se precisar separar
typedef struct {
    char userId[50];
    char data[256];
} SessionData;

// Declarações das funções de gerenciamento de sessão
int createSession(const SessionData* data);
Session getSession(const char* sessionId);
int updateSession(const char* sessionId, const SessionData* newData);
int deleteSession(const char* sessionId);
void manageSessionExpiration();

char* generateToken(const char* sessionId);
bool validateToken(const char* token);

// Funções para conectar ao MongoDB
mongoc_client_t* connectMongoDB();
void insertSession(const Session* session);
Session retrieveSession(const char* sessionId);
void updateSessionInMongo(const Session* session);
void deleteSessionFromMongo(const char* sessionId);

// Funções para RabbitMQ (placeholder para implementação futura)
void connectRabbitMQ();
void publishMessage(const char* message);
void consumeMessages();

bool authenticateRequest(const HttpRequest* request);
void logRequest(const HttpRequest* request);

#endif // SESSION_MANAGER_H
