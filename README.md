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

