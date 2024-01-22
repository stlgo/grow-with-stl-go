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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"stl-go/grow-with-stl-go/pkg/log"
)

const (
	// keypair details
	keySize        = 4096 // 4k key
	privateKeyType = "RSA PRIVATE KEY"
	publicKeyType  = "CERTIFICATE"

	// certificate request details
	cn = "localhost" // common name
	o  = "stl-go"    // organization

	// ObfuscatedPrefix encrypted key prefix
	ObfuscatedPrefix string = "obf::"
)

// GenerateSecret returns a 32 byte AES key
func GenerateSecret() (*string, error) {
	key := make([]byte, 16)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	secret := fmt.Sprintf("%x", key)
	return &secret, nil
}

// HashPassword will hash passwords for storage
func HashPassword(password *string) (*string, error) {
	if password != nil {
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(*password), 4)
		hash := string(hashBytes)
		return &hash, err
	}
	return nil, errors.New("password is nil")
}

// HashCompare will hash & compare a password to the stored hash
func HashCompare(original, request *string) error {
	if original != nil && request != nil {
		return bcrypt.CompareHashAndPassword([]byte(*original), []byte(*request))
	}
	return errors.New("nil input of original or requested strings")
}

// Encrypt takes plaintext and turns it into ciphertext
func Encrypt(plaintext, secret *string) (*string, error) {
	if plaintext != nil && secret != nil {
		c, err := aes.NewCipher([]byte(*secret))
		if err != nil {
			return nil, err
		}

		gcm, err := cipher.NewGCM(c)
		if err != nil {
			return nil, err
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err = rand.Read(nonce); err != nil {
			return nil, err
		}

		cipherText := fmt.Sprintf("%s%s", ObfuscatedPrefix, base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(*plaintext), nil)))
		return &cipherText, nil
	}
	return nil, errors.New("plaintext or secret is nil")
}

// Decrypt takes ciphertext and returns plaintext
func Decrypt(encodedData, secret *string) (*string, error) {
	if encodedData != nil && secret != nil {
		if !strings.Contains(*encodedData, ObfuscatedPrefix) {
			return encodedData, nil
		}

		cipherTextEncoded, err := base64.URLEncoding.DecodeString(strings.ReplaceAll(*encodedData, ObfuscatedPrefix, ""))
		if err != nil {
			return nil, err
		}

		cipherBlock, err := aes.NewCipher([]byte(*secret))
		if err != nil {
			return nil, err
		}

		aead, err := cipher.NewGCM(cipherBlock)
		if err != nil {
			return nil, err
		}

		nonceSize := aead.NonceSize()
		if len(cipherTextEncoded) < nonceSize {
			return nil, err
		}

		nonce, cipherText := cipherTextEncoded[:nonceSize], cipherTextEncoded[nonceSize:]
		plainTextEncoded, err := aead.Open(nil, nonce, cipherText, nil)
		if err != nil {
			return nil, err
		}

		plainText := string(plainTextEncoded)
		return &plainText, nil
	}
	return nil, errors.New("encodedData or secret is nil")
}

// GeneratePrivateKey will a pem encoded private key and an rsa private key object
func GeneratePrivateKey() ([]byte, *rsa.PrivateKey, error) {
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
func GeneratePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
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
