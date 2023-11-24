package lnpdf

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"latona-pdf/pkg/lnpdf/utils"
	"os"
	"time"

	pdfsign "github.com/digitorus/pdfsign/sign"
)

func Sign(pdfPath *string, privateKeyPath *string, certificatePath *string, chainPath *string, signingInfo *SigningInfo) {
	certificate_data, err := os.ReadFile(*certificatePath)
	if err != nil {
		panic(err)
	}
	certificate_data_block, _ := pem.Decode(certificate_data)
	if certificate_data_block == nil {
		panic(err)
	}

	cert, err := x509.ParseCertificate(certificate_data_block.Bytes)
	if err != nil {
		panic(err)
	}

	key_data, err := os.ReadFile(*privateKeyPath)
	if err != nil {
		panic(err)
	}
	key_data_block, _ := pem.Decode(key_data)
	if key_data_block == nil {
		panic(err)
	}

	pkey, err := x509.ParsePKCS1PrivateKey(key_data_block.Bytes)
	if err != nil {
		panic(err)
	}

	certificate_chains := make([][]*x509.Certificate, 0)

	if *chainPath != "" {
		certificate_pool := x509.NewCertPool()
		if err != nil {
			panic(err)
		}

		chain_data, err := os.ReadFile(*chainPath)
		if err != nil {
			panic(err)
		}

		certificate_pool.AppendCertsFromPEM(chain_data)
		certificate_chains, err = cert.Verify(x509.VerifyOptions{
			Intermediates: certificate_pool,
			CurrentTime:   cert.NotBefore,
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		})
		if err != nil {
			panic(err)
		}
	}

	tempPath := *pdfPath + "_"
	utils.Copy(pdfPath, &tempPath)
	defer os.Remove(tempPath)

	err = pdfsign.SignFile(*pdfPath, tempPath, pdfsign.SignData{
		Signature: pdfsign.SignDataSignature{
			Info: pdfsign.SignDataSignatureInfo{
				Name:        (*signingInfo).Name,
				Location:    (*signingInfo).Location,
				Reason:      (*signingInfo).Reason,
				ContactInfo: (*signingInfo).ContactInfo,
				Date:        time.Now().Local(),
			},
			CertType:   pdfsign.CertificationSignature,
			DocMDPPerm: pdfsign.AllowFillingExistingFormFieldsAndSignaturesPerms,
		},
		Signer:            pkey,
		DigestAlgorithm:   crypto.SHA256,
		Certificate:       cert,
		CertificateChains: certificate_chains,
		TSA: pdfsign.TSA{
			URL: (*signingInfo).TsrUrl,
		},
	})
	if err != nil {
		panic(err)
	}

	utils.Copy(&tempPath, pdfPath)
}

type SigningInfo struct {
	Name        string
	Location    string
	Reason      string
	ContactInfo string
	TsrUrl      string
}
