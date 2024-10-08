# Entendendo o DHCP Request

O DHCP (Dynamic Host Configuration Protocol) é um protocolo de rede amplamente utilizado para configurar automaticamente os dispositivos em uma rede IP, atribuindo endereços IP e outras configurações de rede de forma dinâmica. O DHCP funciona com base em um conjunto de mensagens trocadas entre clientes DHCP e servidores DHCP.

## O que é o DHCP Request?

O DHCP Request é uma mensagem enviada por um cliente DHCP para solicitar a confirmação de um endereço IP oferecido pelo servidor DHCP. Essa mensagem é uma parte crucial do processo de atribuição dinâmica de endereços IP e é utilizada para que o cliente comunique ao servidor DHCP sua escolha final de configurações de rede.

### Detalhes da Mensagem DHCP Request

- **Finalidade**: Confirmar a oferta de configuração de rede feita pelo servidor DHCP.
- **Tipo de Mensagem**: Identificado pelo código 3 (DHCP Request) no campo `op` (Operation Code) da mensagem DHCP.
- **Conteúdo**:
    - **Endereço IP Oferecido**: Incluído no campo `ciaddr` (Client IP Address) da mensagem DHCP Request, indicando o endereço IP que o cliente deseja confirmar.
    - **Endereço MAC do Cliente**: Especificado no campo `chaddr` (Client Hardware Address), fornecendo o endereço MAC único do cliente DHCP.
    - **Identificador da Transação (XID)**: Garante a associação correta entre mensagens DHCP Request e DHCP Offer.
    - **Servidor DHCP Selecionado**: Opcionalmente, o cliente pode especificar o servidor DHCP escolhido no campo `server identifier` para identificar a fonte da oferta DHCP.

### Fluxo de Mensagens DHCP

1. **DHCP Discover**: O cliente DHCP envia uma mensagem DHCP Discover para encontrar servidores DHCP disponíveis na rede.
2. **DHCP Offer**: O servidor DHCP responde com uma mensagem DHCP Offer, oferecendo ao cliente um endereço IP e outras configurações de rede.
3. **DHCP Request**: O cliente DHCP seleciona uma oferta recebida no DHCP Offer e envia uma mensagem DHCP Request para confirmar a escolha do endereço IP oferecido.
4. **DHCP Acknowledge**: O servidor DHCP recebe o DHCP Request, confirma a escolha do cliente e envia uma mensagem DHCP Acknowledge, atribuindo definitivamente o endereço IP ao cliente.

### Exemplo de Uso do DHCP Request

```plaintext
-- Exemplo de Mensagem DHCP Request --
Operação: DHCP Request
Endereço IP Oferecido: 192.168.1.100
Endereço MAC do Cliente: 01:23:45:67:89:ab
Identificador da Transação (XID): 0x12345678
```

- Na mensagem acima, o cliente DHCP está solicitando ao servidor DHCP a confirmação do endereço IP oferecido (192.168.1.100) associado ao seu endereço MAC específico.

