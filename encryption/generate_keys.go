package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/luckyAkbar/stdlib/helper"
)

// KeyGenerationOpts option to generate key
type KeyGenerationOpts struct {
	Random          io.Reader
	Bits            int
	PEMFormat       bool
	GOBFormat       bool
	PublicFilename  string
	PrivateFilename string
}

var defaultKeyGenerationOpts = &KeyGenerationOpts{
	Random:          rand.Reader,
	Bits:            2048,
	PEMFormat:       true,
	GOBFormat:       true,
	PublicFilename:  "public",
	PrivateFilename: "private",
}

// GenerateKey generate private and public key. Based on value in config,
// you can choose what file to be generated or no file generated at all
// either way, will still return *rsa.PrivateKey.
// if failure happen when generating the file(s), will call os.Exit(1)
func GenerateKey(config *KeyGenerationOpts) (*rsa.PrivateKey, error) {
	if config == nil {
		config = defaultKeyGenerationOpts
	}

	privKey, err := rsa.GenerateKey(config.Random, config.Bits)
	if err != nil {
		return nil, err
	}

	publicKey := privKey.PublicKey

	if config.GOBFormat {
		saveGobKey(config.PrivateFilename+".key", privKey)
		saveGobKey(config.PublicFilename+".key", publicKey)
	}

	if config.PEMFormat {
		savePEMKey(config.PrivateFilename+".pem", privKey)
		savePublicPEMKey(config.PublicFilename+".pem", publicKey)
	}

	return privKey, nil
}

func saveGobKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer helper.WrapCloser(outFile.Close)

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer helper.WrapCloser(outFile.Close)

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer helper.WrapCloser(pemfile.Close)

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
