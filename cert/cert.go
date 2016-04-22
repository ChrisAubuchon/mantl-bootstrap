package cert

import (
	"crypto/rand"
	"crypto/x509"
	"net"

	"github.com/spf13/viper"
	"github.com/square/certstrap/pkix"
)

type Cert struct {
	Cert []byte
	Key []byte
}

func GenerateCaCert(config *viper.Viper) (*Cert, error) {
	key, err := pkix.CreateRSAKey(2048)
	if err != nil {
		return nil, err
	}

	cert, err := pkix.CreateCertificateAuthority(
		key, 
		viper.GetString("cert-unit"),
		1,
		viper.GetString("cert-organization"),
		viper.GetString("cert-country"),
		viper.GetString("cert-state"),
		viper.GetString("cert-locality"),
		viper.GetString("cert-common"),
		)
	if err != nil {
		return nil, err
	}

	return newCert(cert, key)
}

func GenerateCert(caCert []byte, caKey []byte, ip string, domain string) (*Cert, error) {
	key, err := pkix.CreateRSAKey(2048)
	if err != nil {
		return nil, err
	}

	pcert, err := pkix.NewCertificateFromPEM(caCert)
	if err != nil {
		return nil, err
	}
	rawpcert, err := pcert.GetRawCertificate()
	if err != nil {
		return nil, err
	}

	pkey, err := pkix.NewKeyFromPrivateKeyPEM(caKey)
	if err != nil {
		return nil, err
	}

	ips := make([]net.IP, 2)
	ips[0] = net.ParseIP("127.0.0.1")
	ips[1] = net.ParseIP(ip)
	dnsnames := make([]string, 3)
	dnsnames[0] = "*.service.consul"
	dnsnames[0] = "*.node.consul"
	dnsnames[0] = "localhost"

	csrTemplate := &x509.CertificateRequest{
		Subject:     rawpcert.Subject,
		IPAddresses: ips,
		DNSNames:    dnsnames,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, key.Private)
	if err != nil {
		return nil, err
	}

	csr := pkix.NewCertificateSigningRequestFromDER(csrBytes)

	cert, err := pkix.CreateCertificateHost(pcert, pkey, csr, 1)
	if err != nil {
		return nil, err
	}

	return newCert(cert, key)
}

func newCert(cert *pkix.Certificate, key *pkix.Key) (*Cert, error) {
	certBytes, err := cert.Export()
	if err != nil {
		return nil, err
	}

	keyBytes, err := key.ExportPrivate()
	if err != nil {	
		return nil, err
	}

	return &Cert{
		Cert: certBytes,
		Key: keyBytes,
	}, nil
}
