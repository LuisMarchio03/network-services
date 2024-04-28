# Cliente DHCP

O cliente DHCP é um programa que utiliza o Protocolo de Configuração Dinâmica de Host (DHCP) para obter configurações de rede automaticamente de um servidor DHCP. Ele é responsável por enviar mensagens DHCP para o servidor, como DHCP Discover e DHCP Request, e receber as respostas correspondentes, como DHCP Offer e DHCP ACK.

## DHCP (Protocolo de Configuração Dinâmica de Host)

O DHCP é um protocolo de rede utilizado para atribuir automaticamente endereços IP e outras configurações de rede (como máscara de sub-rede, gateway padrão, servidores DNS, etc.) para dispositivos em uma rede IP. O DHCP simplifica a administração de redes, permitindo que os dispositivos obtenham configurações de rede dinamicamente sem intervenção manual.

## Checklist de Implementação do Cliente DHCP

- [x] Configuração de uma conexão UDP para receber respostas do servidor DHCP.
- [x] Envio de uma mensagem DHCP Discover para solicitar configurações de rede.
- [x] Processamento da oferta (DHCP Offer) recebida do servidor DHCP.
- [x] Envio de uma mensagem DHCP Request para confirmar as configurações oferecidas.
- [x] Recebimento e processamento do ACK (DHCP ACK) do servidor DHCP.
- [ ] Tratamento de NAK (Not Acknowledge).
- [ ] Lógica de Reenvio para mensagens não recebidas.
- [ ] Gerenciamento de Conexão robusto.
- [ ] Implementação de Logging e Tratamento de Erros.
- [ ] Adição de suporte para mais opções DHCP (DNS, gateway, etc.).
