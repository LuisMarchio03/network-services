# Documentação do Projeto

Este documento fornece informações sobre como acessar os diferentes serviços configurados no projeto.

## Serviços e Credenciais

### MongoDB

- **URL de Acesso**: `mongodb://localhost:27017/mydb`
- **Usuário**: Não é necessário usuário e senha para a configuração padrão do MongoDB.
- **Senha**: Não aplicável.

### RabbitMQ

- **URL de Acesso**: [RabbitMQ Management](http://localhost:15672)
- **Usuário**: `guest`
- **Senha**: `guest`
- **Porta do Broker**: `5672`

### Elasticsearch

- **URL de Acesso**: `http://localhost:9200`
- **Usuário**: Não é necessário usuário e senha para a configuração padrão do Elasticsearch.
- **Senha**: Não aplicável.

### Kibana

- **URL de Acesso**: [Kibana](http://localhost:5601)
- **Usuário**: Não é necessário usuário e senha para a configuração padrão do Kibana.
- **Senha**: Não aplicável.

## Informações Adicionais

- Certifique-se de que os serviços estão em execução antes de tentar acessá-los. Você pode iniciar todos os serviços utilizando o comando:

  ```bash
  docker-compose up --build

## Dockerfile

docker build -t my-network-services .
docker run --network my-network my-network-services

## RabbitMQ

- É Necessario criar a fila manualmente no RabbitMQ.

## Elasticsearch + Kibana

- É Necessario criar o INDEX manualmente no Kibana.
- Alem de ser necessario passar todos os dados para o INDEX, também é necessário passar os dados para o INDEX do Elasticsearch - Seguindo a estrutura dentro do log-receiver GOLANG.