package bootstrap

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type CertData struct {
	Country string
	State string
	Locality string
	Org string
	Unit string
	Common string
}

func GenerateCaCert(cd *CertData) (string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	template, err := getTemplate(cd)
	if err != nil {
		return "", err
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	bytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return "", err
	}

	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: bytes})), nil
}

func getTemplate(cd *CertData) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365*24*time.Hour)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country: []string{cd.Country},
			Organization: []string{cd.Org},
			OrganizationalUnit: []string{cd.Unit},
			Locality: []string{cd.Locality},
			Province: []string{cd.State},
			CommonName: cd.Common,
		},
		NotBefore: notBefore,
		NotAfter: notAfter,
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	return &template, nil
}
