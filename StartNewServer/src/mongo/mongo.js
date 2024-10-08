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
} = require('../models/models'); // Ajuste o caminho se necessário

// Conexão com o MongoDB
mongoose.connect('mongodb://localhost:27017/mydb', {
    useNewUrlParser: true,
    useUnifiedTopology: true,
}).then(() => {
    console.log('Conectado ao MongoDB');
}).catch(err => {
    console.error('Erro ao conectar ao MongoDB:', err);
});

// Função para criar um novo documento genérico
const createDocument = async (Model, data) => {
    try {
        const document = new Model(data);
        return await document.save();
    } catch (error) {
        console.error('Erro ao criar documento:', error);
        throw new Error('Erro ao salvar no banco de dados');
    }
};

module.exports = {
    createDocument,
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
};
