const mongoose = require('mongoose');
const {
    Dhcp,
    Dns,
    HttpServer,
    FtpServer,
    EmailServer,
    FileServer,
    Firewall,
    LdapServer,
    LogServer,
    UserComputer,
    Session,
} = require('./models/models');
const fs = require("fs"); // Ajuste o caminho se necessário

require('./rabbitmq-elasticsearch')

// String de conexão com o MongoDB
const uri = 'mongodb://mongodb:27017/mydb';

const dataToInsert = {
    dhcpEntries: [
        { ip_address: '192.168.1.10', assigned: true },
        { ip_address: '192.168.1.11', assigned: false },
    ],
    dnsEntries: [
        { domain: 'example.com', ip_address: '192.168.1.20' },
    ],
    httpEntries: [
        { port: 80, url: 'http://example.com', status: 'up' },
    ],
    ftpEntries: [
        { port: 21, status: 'active' },
    ],
    emailEntries: [
        { smtp_server: 'smtp.example.com', port: 587, is_active: true },
    ],
    fileEntries: [
        { path: '/var/files', permissions: 'read-write' },
    ],
    firewallEntries: [
        { rule: 'ALLOW 192.168.1.0/24', action: 'allow' },
    ],
    ldapEntries: [
        { host: 'ldap.example.com', port: 389, base_dn: 'dc=example,dc=com' },
    ],
    logEntries: [
        { log_type: 'INFO', message: 'Serviço DHCP iniciado' },
    ],
    computerEntries: [
        { hostname: 'PC1', ip_address: '192.168.1.10', os: 'Windows', browser: 'Chrome' },
        { hostname: 'PC2', ip_address: '192.168.1.11', os: 'Linux', browser: 'Firefox' },
    ],
    sessionEntries: [
        { clientId: '5f81a3a2c3b4c1a1f2e3d4e5', sessionToken: '1234567890', createdAt: new Date(), expiresAt: new Date(), machineId: '5f81a3a2c3b4c1a1f2e3d4e6' },
    ],
};

// Função para criar documentos no Mongoose
async function createDocument(Model, document) {
    try {
        const result = await Model.create(document);
        console.log(`Documento inserido com sucesso: ${result._id}`);
    } catch (err) {
        console.error(`Erro ao inserir documento: ${err.message}`);
    }
}

async function run() {
    try {
        function generateIPs() {
            return `${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}`;
        }

        function createDhcp(numeroDeRegistros) {
            let dhcpArray = [];

            for (let i = 0; i < numeroDeRegistros; i++) {
                // Criando o DHCP com um IP aleatório
                let ipAddress = generateIPs();
                let dhcpRecord = {
                    ip_address: ipAddress,
                    assigned: false,
                };

                dhcpArray.push(dhcpRecord);
            }

            console.log(JSON.stringify(dhcpArray));

            console.log(`Gerados ${numeroDeRegistros} registros de DHCP e DNS!`);

            return JSON.stringify(dhcpArray, null, 2)
        }

        const numeroDeRegistros = 100;

        const DhcpJSON = createDhcp(numeroDeRegistros);

        // Conectando ao cliente MongoDB via Mongoose
        await mongoose.connect(uri, { useNewUrlParser: true, useUnifiedTopology: true });

        // Inserindo os documentos nas coleções usando Mongoose
        await Promise.all([
            ...dataToInsert.dhcpEntries.map(entry => createDocument(Dhcp, DhcpJSON)),
            ...dataToInsert.dnsEntries.map(entry => createDocument(Dns, entry)),
            ...dataToInsert.httpEntries.map(entry => createDocument(HttpServer, entry)),
            ...dataToInsert.ftpEntries.map(entry => createDocument(FtpServer, entry)),
            ...dataToInsert.emailEntries.map(entry => createDocument(EmailServer, entry)),
            ...dataToInsert.fileEntries.map(entry => createDocument(FileServer, entry)),
            ...dataToInsert.firewallEntries.map(entry => createDocument(Firewall, entry)),
            ...dataToInsert.ldapEntries.map(entry => createDocument(LdapServer, entry)),
            ...dataToInsert.logEntries.map(entry => createDocument(LogServer, entry)),
            ...dataToInsert.computerEntries.map(entry => createDocument(UserComputer, entry)),
            ...dataToInsert.sessionEntries.map(entry => createDocument(Session, entry)),
        ]);

        console.log('Todos os dados foram inseridos com sucesso');
    } catch (error) {
        console.error('Erro ao executar o script:', error.message);
    } finally {
        // Fecha a conexão com o MongoDB
        await mongoose.connection.close();
    }
}

run().catch(console.dir);
