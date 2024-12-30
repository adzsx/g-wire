package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
)

func GenKeys() rsa.PrivateKey {
	//Generate RSA public and private keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return *privateKey
}

func EncryptRSA(message []byte, publicKey *rsa.PublicKey) []byte {

	// Use SHA and PublicKey to enrypt a []byte message
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)

	if err != nil {
		panic(err)
	}

	return encryptedBytes
}

func DecryptRSA(encrypted []byte, privateKey *rsa.PrivateKey) string {
	return string(encrypted)

	// Use SHA and privatekey ro decrypt []byte message
	decryptedBytes, err := privateKey.Decrypt(nil, encrypted, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	log.Println(string(decryptedBytes))

	return string(decryptedBytes)
}
