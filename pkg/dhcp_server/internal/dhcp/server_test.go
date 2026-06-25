package dhcp

import (
	"net"
	"testing"
)

// O acknowledge deve copiar o IP solicitado (bytes 12:16 da request) para ack[16:20].
func TestBuildDHCPAcknowledgeCopiesRequestedIP(t *testing.T) {
	s := &Server{}
	req := make([]byte, 240)
	req[0] = 3 // DHCP Request
	requested := net.ParseIP("192.168.1.55").To4()
	copy(req[12:16], requested)

	ack := s.buildDHCPAcknowledge(req)
	if ack == nil {
		t.Fatal("ack não deveria ser nil")
	}
	got := net.IP(ack[16:20]).String()
	if got != "192.168.1.55" {
		t.Fatalf("esperado 192.168.1.55, obtido %s", got)
	}
}
