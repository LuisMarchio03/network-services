// Importar dependências
const amqp = require('amqplib');
const { Client } = require('@elastic/elasticsearch');

// Configurações do RabbitMQ e Elasticsearch
const RABBITMQ_URL = 'amqp://guest:guest@localhost:5672/';
const ELASTICSEARCH_URL = 'http://localhost:9200';
const QUEUE_NAME = 'logs_queue';
const INDEX_NAME = 'logs';

// Dados de log de exemplo para enviar ao RabbitMQ
const logEntry = {
    timestamp: new Date().toISOString(),
    service: "my_service",
    level: "INFO",
    message: "This is a test log",
    context: { user: "user123", action: "login" },
    computer_client_ip: "192.168.1.1",
};

// Função para criar a fila no RabbitMQ
async function createQueue() {
    try {
        const conn = await amqp.connect(RABBITMQ_URL); // Conecta ao RabbitMQ
        const channel = await conn.createChannel(); // Cria um canal
        await channel.assertQueue(QUEUE_NAME, { durable: true }); // Assegura que a fila existe
        console.log(`Queue "${QUEUE_NAME}" created or exists.`);

        // Enviar mensagem de log para a fila
        channel.sendToQueue(QUEUE_NAME, Buffer.from(JSON.stringify(logEntry)));
        console.log(`Log sent to RabbitMQ: ${JSON.stringify(logEntry)}`);

        setTimeout(() => conn.close(), 500); // Fecha a conexão após meio segundo
    } catch (err) {
        console.error('Error creating queue or sending message:', err);
    }
}

// Função para criar o índice no Elasticsearch
async function createIndex() {
    const client = new Client({ node: ELASTICSEARCH_URL });

    try {
        const indexExists = await client.indices.exists({ index: INDEX_NAME });

        if (!indexExists.body) {
            // Se o índice não existe, criá-lo com mapeamento
            await client.indices.create({
                index: INDEX_NAME,
                body: {
                    mappings: {
                        properties: {
                            timestamp: { type: 'date' },
                            service: { type: 'text' },
                            level: { type: 'keyword' },
                            message: { type: 'text' },
                            context: { type: 'object' },
                            computer_client_ip: { type: 'ip' }
                        }
                    }
                }
            });
            console.log(`Index "${INDEX_NAME}" created in Elasticsearch.`);
        } else {
            console.log(`Index "${INDEX_NAME}" already exists.`);
        }

        // Enviar log para o Elasticsearch
        const response = await client.index({
            index: INDEX_NAME,
            body: logEntry
        });
        console.log('Log sent to Elasticsearch:', response.body);

    } catch (err) {
        console.error('Error creating index or sending log:', err);
    }
}

// Executar funções
(async () => {
    await createQueue();     // Automatizar a criação da fila no RabbitMQ e enviar o log
    await createIndex();     // Automatizar a criação do índice no Elasticsearch e enviar o log
})();
