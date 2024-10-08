# Usar uma imagem base do Node.js
FROM node:14

# Criar diretório de trabalho
WORKDIR /usr/src/app

# Copiar package.json e package-lock.json
COPY package*.json ./

# Instalar dependências
RUN npm install

# Copiar o restante dos arquivos
COPY . .

# Comando para iniciar a aplicação
CMD ["node", "src/setup.js"]
