# Servidor DHCP (Dynamic Host Configuration Protocol)

Um servidor DHCP é um componente de rede responsável por atribuir automaticamente endereços IP e outras configurações de rede para dispositivos que se conectam a uma rede. Ele simplifica o processo de configuração de rede, eliminando a necessidade de configurar manualmente cada dispositivo com um endereço IP exclusivo, máscara de sub-rede, gateway padrão, servidor DNS, entre outros.

O servidor DHCP opera seguindo uma abordagem de cliente/servidor. Quando um dispositivo se conecta à rede, ele envia uma solicitação DHCP (DHCP Discover) para localizar um servidor DHCP disponível. O servidor DHCP responde com uma oferta de configuração (DHCP Offer), que inclui um endereço IP disponível para o dispositivo. O dispositivo aceita a oferta e solicita a configuração completa (DHCP Request), e o servidor DHCP concede a configuração (DHCP Acknowledge), fornecendo todas as informações necessárias para que o dispositivo se conecte à rede.

A utilização de um servidor DHCP simplifica significativamente a administração de redes, especialmente em redes grandes, onde a configuração manual de cada dispositivo seria impraticável.

