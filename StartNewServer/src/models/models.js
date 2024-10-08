const mongoose = require('mongoose');

/**
 * Modelo DHCP
 *
 * O modelo DHCP é responsável por armazenar os endereços IP que o servidor DHCP aloca ou atribui a dispositivos na rede.
 *
 * Campos:
 * - ip_address: Endereço IP que será atribuído. Deve ser um endereço IP válido no formato IPv4 (Ex: '192.168.0.1').
 * - assigned: Booleano indicando se o endereço IP foi atribuído (true) ou está disponível (false).
 */
const DhcpSchema = new mongoose.Schema({
    ip_address: {
        type: String,
        required: true,
        match: /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, // Regex para validar IP
    },
    assigned: {
        type: Boolean,
        required: true, // Obrigatório indicar se o IP está ou não atribuído
    },
});

const Dhcp = mongoose.model('Dhcp', DhcpSchema);

/**
 * Modelo DNS
 *
 * O modelo DNS armazena os registros DNS da rede, mapeando nomes de domínio para endereços IP.
 *
 * Campos:
 * - domain: Nome de domínio associado ao endereço IP (Ex: 'example.com').
 * - ip_address: Endereço IP correspondente ao domínio. Deve ser um IP válido no formato IPv4.
 */
const DnsSchema = new mongoose.Schema({
    domain: {
        type: String,
        required: true, // Domínio obrigatório
    },
    ip_address: {
        type: String,
        required: true,
        match: /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, // Regex para validar IP
    },
});

const Dns = mongoose.model('Dns', DnsSchema);

/**
 * Modelo HTTP Server
 *
 * O modelo do servidor HTTP define portas e URLs dos servidores web na rede.
 *
 * Campos:
 * - port: Número da porta usada pelo servidor HTTP (Ex: 80 para HTTP, 443 para HTTPS).
 * - url: URL associada ao servidor.
 * - status: Indica se o servidor está "up" (em operação) ou "down" (fora de operação). O valor padrão é 'up'.
 */
const HttpServerSchema = new mongoose.Schema({
    port: {
        type: Number,
        required: true, // Porta obrigatória
    },
    url: {
        type: String,
        required: true, // URL associada ao servidor é obrigatória
    },
    status: {
        type: String,
        enum: ['up', 'down'], // Servidor pode estar "up" ou "down"
        default: 'up', // Padrão é 'up'
    },
});

const HttpServer = mongoose.model('HttpServer', HttpServerSchema);

/**
 * Modelo FTP Server
 *
 * Define a configuração de servidores FTP na rede.
 *
 * Campos:
 * - port: Porta utilizada pelo servidor FTP.
 * - status: Define se o servidor está 'active' (ativo) ou 'inactive' (inativo). O padrão é 'active'.
 */
const FtpServerSchema = new mongoose.Schema({
    port: {
        type: Number,
        required: true, // Porta FTP obrigatória
    },
    status: {
        type: String,
        enum: ['active', 'inactive'], // Pode estar 'active' ou 'inactive'
        default: 'active', // Padrão é 'active'
    },
});

const FtpServer = mongoose.model('FtpServer', FtpServerSchema);

/**
 * Modelo Email Server
 *
 * Representa servidores de e-mail (SMTP) na rede.
 *
 * Campos:
 * - smtp_server: Nome do servidor SMTP.
 * - port: Porta utilizada pelo servidor SMTP.
 * - is_active: Define se o servidor de e-mail está ativo. O valor padrão é true.
 */
const EmailServerSchema = new mongoose.Schema({
    smtp_server: {
        type: String,
        required: true, // Nome do servidor SMTP é obrigatório
    },
    port: {
        type: Number,
        required: true, // Porta SMTP obrigatória
    },
    is_active: {
        type: Boolean,
        default: true, // O servidor começa como ativo por padrão
    },
});

const EmailServer = mongoose.model('EmailServer', EmailServerSchema);

/**
 * Modelo File Server
 *
 * Representa um servidor de arquivos na rede.
 *
 * Campos:
 * - path: Caminho do servidor onde os arquivos estão armazenados.
 * - permissions: Permissões atribuídas aos arquivos, sendo 'read-write' o valor padrão.
 */
const FileServerSchema = new mongoose.Schema({
    path: {
        type: String,
        required: true, // Caminho dos arquivos é obrigatório
    },
    permissions: {
        type: String,
        default: 'read-write', // Padrão: leitura e escrita
    },
});

const FileServer = mongoose.model('FileServer', FileServerSchema);

/**
 * Modelo Firewall
 *
 * Define regras do firewall na rede.
 *
 * Campos:
 * - rule: Regra de firewall (Ex: bloquear ou permitir uma porta).
 * - action: Ação a ser realizada pela regra ('allow' para permitir ou 'deny' para negar).
 */
const FirewallSchema = new mongoose.Schema({
    rule: {
        type: String,
        required: true, // Regra é obrigatória
    },
    action: {
        type: String,
        enum: ['allow', 'deny'], // Pode ser 'allow' ou 'deny'
        required: true, // Ação obrigatória
    },
});

const Firewall = mongoose.model('Firewall', FirewallSchema);

/**
 * Modelo LDAP Server
 *
 * Define a configuração de servidores LDAP na rede, usados para serviços de diretório.
 *
 * Campos:
 * - host: Endereço do servidor LDAP.
 * - port: Porta usada pelo LDAP.
 * - base_dn: Base DN (Distinguished Name) da árvore LDAP.
 */
const LdapServerSchema = new mongoose.Schema({
    host: {
        type: String,
        required: true, // Host do LDAP é obrigatório
    },
    port: {
        type: Number,
        required: true, // Porta do LDAP obrigatória
    },
    base_dn: {
        type: String,
        required: true, // Base DN é obrigatória
    },
});

const LdapServer = mongoose.model('LdapServer', LdapServerSchema);

/**
 * Modelo Log Server
 *
 * Registra logs-producer-teste de atividades na rede.
 *
 * Campos:
 * - log_type: Tipo de log (Ex: 'error', 'access', 'event').
 * - created_at: Data de criação do log, padrão é a data atual.
 * - message: Mensagem ou descrição do log.
 */
const LogServerSchema = new mongoose.Schema({
    log_type: {
        type: String,
        required: true, // Tipo de log obrigatório
    },
    created_at: {
        type: Date,
        default: Date.now, // Padrão é a data atual
    },
    message: {
        type: String,
        required: true, // Mensagem de log obrigatória
    },
});

const LogServer = mongoose.model('LogServer', LogServerSchema);

/**
 * Modelo User Computer
 *
 * Representa um computador de usuário na rede.
 *
 * Campos:
 * - hostname: Nome do computador.
 * - ip_address: Endereço IP do computador.
 * - os: Sistema operacional do computador (Ex: Windows, Linux).
 * - browser: Navegador utilizado (Ex: Chrome, Firefox).
 */
const UserComputerSchema = new mongoose.Schema({
    hostname: {
        type: String,
        required: true, // Nome do computador obrigatório
    },
    ip_address: {
        type: String,
        required: true,
        match: /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, // Regex para validar IP
    },
    os: {
        type: String,
        required: true, // Sistema operacional é obrigatório
    },
    browser: {
        type: String,
        required: true, // Navegador é obrigatório
    },
});

const UserComputer = mongoose.model('UserComputer', UserComputerSchema);

module.exports = {
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
};
