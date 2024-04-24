# Estrutura do DHCP Offer

O DHCP Offer (Oferta DHCP) é uma mensagem enviada por um servidor DHCP em resposta a uma solicitação DHCP Discover de um cliente. Esta mensagem contém informações sobre a configuração de rede oferecida pelo servidor ao cliente.

## Componentes do DHCP Offer:

### Tipo de Mensagem:
- **Campo**: Tipo (OpCode)
- **Descrição**: Identifica o tipo de mensagem DHCP. No DHCP Offer, o valor deste campo é 2.

### Endereço IP Oferecido:
- **Campo**: Your IP Address (YIAddr)
- **Descrição**: Indica o endereço IP oferecido pelo servidor DHCP ao cliente.

### Endereço do Servidor DHCP:
- **Campo**: Server Identifier (SIAddr)
- **Descrição**: Especifica o endereço IP do servidor DHCP que está oferecendo a configuração de rede ao cliente.

### Outras Configurações de Rede:
- Outros campos opcionais podem conter informações adicionais, como máscara de sub-rede, gateway padrão, servidor DNS, etc.

- Além do endereço IP oferecido e do endereço do servidor DHCP, as mensagens DHCP Offer e DHCP Acknowledge podem conter várias outras configurações de rede. Essas configurações adicionais ajudam o cliente a configurar sua conexão de rede de maneira adequada.

## Máscara de Sub-rede (Subnet Mask)
- **Campo**: Subnet Mask (Subnet Mask Option)
- **Descrição**: Especifica a máscara de sub-rede a ser usada pelo cliente para determinar quais bits do endereço IP estão relacionados à rede e aos hosts.

## Gateway Padrão (Default Gateway)
- **Campo**: Router (Router Option)
- **Descrição**: Indica o endereço IP do gateway padrão, que é o roteador utilizado pelo cliente para encaminhar pacotes de rede para destinos fora de sua própria rede local.

## Servidor DNS (Domain Name Server)
- **Campo**: Domain Name Server (DNS Option)
- **Descrição**: Fornece os endereços IP dos servidores DNS que o cliente pode usar para resolver nomes de domínio em endereços IP.

## Tempo de Aluguel (Lease Time)
- **Campo**: Lease Time (Lease Time Option)
- **Descrição**: Especifica a duração do tempo de aluguel do endereço IP atribuído ao cliente. Após esse período, o cliente deve renovar sua concessão DHCP.

## Servidor NTP (Network Time Protocol)
- **Campo**: Network Time Protocol Server (NTP Option)
- **Descrição**: Indica o endereço IP do servidor NTP, que o cliente pode usar para sincronizar seu relógio com o horário correto.

## Outras Configurações Específicas do Cliente
- Além das configurações mencionadas acima, outras opções DHCP podem ser incluídas para fornecer configurações específicas do cliente, como informações de VLAN, parâmetros de segurança, etc.
