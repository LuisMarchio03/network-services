#include <stdio.h>
#include <bson/bson.h>
#include <string.h>
#include <time.h>
#include "session.h"

void convertTimeToString(time_t rawtime, char* buffer, size_t size) {
    struct tm * timeinfo;
    timeinfo = localtime(&rawtime);
    strftime(buffer, size, "%Y-%m-%d %H:%M:%S", timeinfo);
}

void sessionToString(const Session* session, char* output, size_t output_size) {
    char created_at_str[20], updated_at_str[20];
    struct tm *tm_info;

    tm_info = gmtime(&session->created_at);
    strftime(created_at_str, sizeof(created_at_str), "%Y-%m-%dT%H:%M:%SZ", tm_info);

    tm_info = gmtime(&session->updated_at);
    strftime(updated_at_str, sizeof(updated_at_str), "%Y-%m-%dT%H:%M:%SZ", tm_info);

    snprintf(output, output_size,
             "{"
             "\"session_id\": \"%s\", "
             "\"user_id\": \"%s\", "
             "\"client_ip\": \"%s\", "
             "\"client_mac\": \"%s\", "
             "\"dhcp_ip\": \"%s\", "
             "\"dhcp_subnet_mask\": \"%s\", "
             "\"token\": \"%s\", "
             "\"expires_in\": %d, "
             "\"created_at\": \"%s\", "
             "\"updated_at\": \"%s\""
             "}",
             session->session_id,
             session->user_id,
             session->client_ip,
             session->client_mac,
             session->dhcp_ip,
             session->dhcp_subnet_mask,
             session->token,
             session->expires_in,
             created_at_str,
             updated_at_str);
}


void extractString(const char* source, const char* key, char* destination, size_t max_len) {
    char* start = strstr(source, key);
    if (start) {
        start = strchr(start, ':');
        start = strchr(start, '\"') + 1;
        char* end = strchr(start, '\"');
        size_t length = end - start;

        if (length > max_len - 1) {
            length = max_len - 1;
        }

        strncpy(destination, start, length);
        destination[length] = '\0';
    }
}

int extractInt(const char* source, const char* key) {
    char* start = strstr(source, key);
    if (start) {
        start = strchr(start, ':') + 1;
        return atoi(start);
    }
    return 0;
}

time_t parseISO8601Time(const char* datetime) {
    struct tm tm;
    memset(&tm, 0, sizeof(struct tm));

    sscanf(datetime, "%4d-%2d-%2dT%2d:%2d:%2d",
           &tm.tm_year, &tm.tm_mon, &tm.tm_mday,
           &tm.tm_hour, &tm.tm_min, &tm.tm_sec);

    tm.tm_year -= 1900;
    tm.tm_mon -= 1;

    time_t t = mktime(&tm);

    return t - timezone;
}

time_t extractTime(const char* source, const char* key) {
    char buffer[20];
    extractString(source, key, buffer, sizeof(buffer));

    return parseISO8601Time(buffer);
}

void deserializeSession(const char* message, Session* session) {
    extractString(message, "\"session_id\"", session->session_id, sizeof(session->session_id));
    extractString(message, "\"user_id\"", session->user_id, sizeof(session->user_id));
    extractString(message, "\"client_ip\"", session->client_ip, sizeof(session->client_ip));
    extractString(message, "\"client_mac\"", session->client_mac, sizeof(session->client_mac));
    extractString(message, "\"dhcp_ip\"", session->dhcp_ip, sizeof(session->dhcp_ip));
    extractString(message, "\"dhcp_subnet_mask\"", session->dhcp_subnet_mask, sizeof(session->dhcp_subnet_mask));
    extractString(message, "\"token\"", session->token, sizeof(session->token));

    session->expires_in = extractInt(message, "\"expires_in\"");

    session->created_at = extractTime(message, "\"created_at\"");
    session->updated_at = extractTime(message, "\"updated_at\"");
}