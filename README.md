# network-services [DEV]

## Simulação de uma rede local, desenvolvida do zero

### Componentes Essenciais
- Servidor DHCP: Responsável por atribuir endereços IP dinâmicos aos dispositivos da rede.
  
- Servidor DNS: Resolve nomes de domínio para endereços IP.

- Servidor TCP/HTTP: Para hospedar serviços web e outras aplicações.

- Servidor FTP: Para transferências de arquivos dentro da rede.

- Servidor de E-mail: Para gerenciar e-mails na rede local.

- Servidor de Arquivos (NFS/Samba): Para compartilhamento de arquivos.

- Firewall: Para gerenciar e proteger o tráfego de rede.

- Servidor de Diretório (LDAP): Para gerenciar usuários e permissões.

- Servidor de Logs: Para monitoramento e registro de eventos na rede.

- Switch e Roteador: Para conectar dispositivos e gerenciar o tráfego entre sub-redes.

### Infraestrutura da Rede

**Setup Físico e Virtual:**
  
  - Hardware: Switches, roteadores, servidores físicos ou máquinas virtuais.
  
  - Software: Sistemas operacionais para servidores (Linux/Windows), software de firewall (iptables/pf), etc.
    
**Configuração da Rede:**
  - Endereçamento IP: Definir uma faixa de IPs para a rede local.

  - Sub-redes: Dividir a rede em sub-redes se necessário (usando VLSM).
  
  - Serviços de Rede:
  
  - DHCP: Configurar para fornecer IPs automaticamente.
  
  - DNS: Configurar para resolver nomes de domínio localmente.
  
  - HTTP/FTP: Configurar servidores para hospedagem e transferência de arquivos.
  
  - E-mail: Configurar servidor de e-mail local.
  
  - Compartilhamento de Arquivos: Configurar NFS/Samba para compartilhamento de arquivos.
  
  - Firewall: Configurar regras para proteger a rede.
  
  - LDAP: Configurar servidor para gerenciamento de usuários.
  
  - Logs: Configurar servidor de logs para monitoramento.

# Mini Curso de Redes Integrado com Desenvolvimento Prático

## Mini Curso de Redes - Parte 1: Fundamentos de Redes de Computadores (Semana de SI)
### Tema Principal: Introdução aos Fundamentos de Redes
#### Aplicação Prática: Implementação de Servidor TCP
- **Título do Curso:** "Construindo e Compreendendo Redes - Parte 1"
- **Conteúdo:**
  - Visão geral dos princípios básicos de redes de computadores
  - Implementação de um servidor TCP simples
  - Demonstração de comunicação básica cliente-servidor
  - Uso de sockets para comunicação de rede
- **Atividade Prática:**
  - Desenvolvimento de um servidor de eco TCP
  - Exercícios para entender o fluxo de dados entre cliente e servidor

## Mini Curso de Redes - Parte 2: Tecnologias Avançadas e Segurança em Redes (Semana Universitária)
### Tema Principal: Explorando Tecnologias Avançadas em Redes
#### Aplicação Prática em Go: Desenvolvimento de Servidor HTTP com Roteamento
- **Título do Curso:** "Explorando Redes Avançadas - Parte 2"
- **Conteúdo:**
  - Implementação de um servidor HTTP com gerenciamento de rotas
  - Exploração de conceitos avançados de redes, incluindo roteamento e segurança
  - Introdução à segurança básica em servidores web
- **Atividade Prática:**
  - Desenvolvimento de um servidor HTTP com roteamento de URLs
  - Discussão sobre práticas de segurança e proteção contra ataques comuns

# Estratégias de Ensino e Aprendizagem
- **Integração Conceitual:** Cada curso aborda um conjunto específico de conceitos de redes, integrados com implementações práticas em Go ou C.
- **Progressão Lógica:** Os participantes iniciam com os fundamentos básicos de redes e avançam para tecnologias mais avançadas, aplicando esses conhecimentos na criação de servidores em Go ou C.
- **Projeto Cumulativo:** Os participantes desenvolvem habilidades ao longo do curso, culminando na criação de servidores TCP e HTTP funcionais.

# Benefícios da Abordagem Integrada
- **Aprendizado Holístico:** Os participantes adquirem conhecimento teórico e prático, compreendendo os fundamentos das redes enquanto aplicam esses conceitos em projetos reais.
- **Preparação para o Mercado:** O desenvolvimento de servidores em Go ou C os prepara para enfrentar desafios reais no campo da computação e redes.

# Recursos Necessários
- **Ambiente de Desenvolvimento:** Voltado para configuração de redes.
- **Material Didático:** Prepare materiais detalhados que combinam teoria de redes com exemplos de código em Go ou C, incluindo exercícios práticos e projetos.

> **Obs:** O desenvolvimento prático não necessariamente será em tempo real, será mais uma apresentação da lógica para criação dos serviços e apresentações dos códigos para os participantes de forma básica. O foco é o desenvolvimento de habilidades em redes e não em programação.

# Etapas Gerais de Desenvolvimento do Projeto para Recriar a Estrutura de uma Rede
- Server DHCP
- Server DNS
- Server TCP
- Firewall
- Server HTTP
- Server FTP
- Mail Server
- Outros servidores opcionais para complementar nossa rede...


