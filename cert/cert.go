package cert

import (
	"crypto/rand"
	"crypto/x509"
	"net"

	"github.com/square/certstrap/pkix"
)

type CertData struct {
	Country string
	State string
	Locality string
	Org string
	Unit string
	Common string
}

type Cert struct {
	Cert []byte
	Key []byte
}

func GenerateCaCert(cd *CertData) (*Cert, error) {
	key, err := pkix.CreateRSAKey(2048)
	if err != nil {
		return nil, err
	}

	cert, err := pkix.CreateCertificateAuthority(key, cd.Unit, 1, cd.Org, cd.Country, cd.State, cd.Locality, cd.Common)
	if err != nil {
		return nil, err
	}

	return newCert(cert, key)
}

func GenerateCert(pcertS string, pkeyS string, ip string, domain string) (*Cert, error) {
	key, err := pkix.CreateRSAKey(2048)
	if err != nil {
		return nil, err
	}

	pcert, err := pkix.NewCertificateFromPEM([]byte(pcertS))
	if err != nil {
		return nil, err
	}
	rawpcert, err := pcert.GetRawCertificate()
	if err != nil {
		return nil, err
	}

	pkey, err := pkix.NewKeyFromPrivateKeyPEM([]byte(pkeyS))
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
