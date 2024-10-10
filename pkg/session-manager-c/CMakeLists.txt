cmake_minimum_required(VERSION 3.10)
project(session_manager_c C)

set(CMAKE_C_STANDARD 11)

# Adicione a linha abaixo para especificar o arquivo de toolchain do vcpkg
set(CMAKE_TOOLCHAIN_FILE "C:/Users/luis_/vcpkg/scripts/buildsystems/vcpkg.cmake")

include_directories(include)
include_directories(src)

add_executable(session_manager_c
        src/main.c
        src/session.c
        src/session.h
        src/token.c
        src/token.h
        src/mongo.c
        src/middleware.c
        src/middleware.h
        src/rabbitmq.c
        src/rabbitmq.h
        include/session_manager.h
        config/config.h
        src/utils.c
        src/utils.h)

# Use o caminho do vcpkg para as configurações do mongo-c-driver
set(CMAKE_PREFIX_PATH "${CMAKE_PREFIX_PATH};C:/Users/luis_/vcpkg/installed/x64-windows")

# Encontre o pacote do mongo-c-driver
find_package(mongoc-1.0 REQUIRED CONFIG)
target_link_libraries(session_manager_c PRIVATE mongo::mongoc_shared)

# Encontrar o pacote RabbitMQ (librabbitmq)
find_package(rabbitmq REQUIRED)
target_link_libraries(session_manager_c PRIVATE rabbitmq::rabbitmq)

# Encontrar a biblioteca OpenSSL (necessária pelo RabbitMQ)
find_package(OpenSSL REQUIRED)
target_link_libraries(session_manager_c PRIVATE OpenSSL::SSL OpenSSL::Crypto)