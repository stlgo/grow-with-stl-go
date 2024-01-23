/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cryptography

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"stl-go/grow-with-stl-go/pkg/log"
	"time"
)

// GenerateDevSSL generates self signed SSL keypair and writes it to file
func GenerateDevSSL(etcDir *string) (privateKey, publicKey *string, err error) {
	if etcDir != nil {
		privateKeyFile := filepath.Join(*etcDir, "key.pem")
		publicKeyFile := filepath.Join(*etcDir, "cert.pem")

		// generate and write out the private key
		log.Warnf("Generating private key %s.  DO NOT USE THIS FOR PRODUCTION", privateKeyFile)
		privateKey, err := generateAndWritePrivateKey(privateKeyFile)
		if err != nil {
			return nil, nil, err
		}

		// generate and write out the private key
		log.Warnf("Generating public key %s.  DO NOT USE THIS FOR PRODUCTION", publicKeyFile)
		err = generateAndWritePublicKey(publicKeyFile, privateKey)
		if err != nil {
			return nil, nil, err
		}

		return &privateKeyFile, &publicKeyFile, nil
	}
	return nil, nil, errors.New("etcDir input is nill for GenerateDevSSL")
}

func generateAndWritePrivateKey(fileName string) (*rsa.PrivateKey, error) {
	privateKeyBytes, privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(fileName, privateKeyBytes, 0o600)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func generateAndWritePublicKey(fileName string, privateKey *rsa.PrivateKey) error {
	if privateKey != nil {
		publicKeyBytes, err := generatePublicKey(privateKey)
		if err != nil {
			return err
		}
		err = os.WriteFile(fileName, publicKeyBytes, 0o600)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("nil input for private key in generateAndWritePublicKey")
}

// GeneratePrivateKey will a pem encoded private key and an rsa private key object
func generatePrivateKey() ([]byte, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Error("Problem generating private key", err)
		return nil, nil, err
	}

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  privateKeyType,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Error("Problem generating private key pem", err)
		return nil, nil, err
	}

	return buf.Bytes(), privateKey, nil
}

// GeneratePublicKey will create a pem encoded cert
func generatePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey != nil {
		template := generateCSR()
		derCert, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return nil, err
		}

		buf := &bytes.Buffer{}
		err = pem.Encode(buf, &pem.Block{
			Type:  publicKeyType,
			Bytes: derCert,
		})
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, errors.New("private key is nil")
}

// generateCSR creates the base information needed to create the certificate
func generateCSR() x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{o},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
}

// CheckCertValidity will check if the cert defined in the conf is not past its not after date
func CheckCertValidity(pemFile *string) error {
	if pemFile != nil {
		r, err := os.ReadFile(*pemFile)
		if err != nil {
			log.Error(err)
			return err
		}

		block, _ := pem.Decode(r)
		_, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	// calculate the validity of the cert
	// TODO: Add a cert check for time based validity here
	// fmt.Println(cert.NotAfter)
	return nil
}
