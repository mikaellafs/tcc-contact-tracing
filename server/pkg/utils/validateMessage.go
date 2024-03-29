package utils

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
)

type ecdsaSignature struct {
	R, S *big.Int
}

func ValidateMessage(message any, pk string, signature []byte) (bool, error) {
	// Parse public key base64 encoded
	// pkBytes, err := base64.StdEncoding.DecodeString(pk)
	pkBytes, err := hex.DecodeString(pk)
	pub, err := x509.ParsePKIXPublicKey(pkBytes)
	if err != nil {
		return false, err
	}
	pubKey, _ := pub.(*ecdsa.PublicKey)

	// Parse signature
	var esig ecdsaSignature
	_, err = asn1.Unmarshal(signature, &esig)
	if err != nil {
		return false, err
	}

	log.Printf("ValidateMessage - Signature R: %d , S: %d", esig.R, esig.S)

	// sha1withECDSA
	msgBytes, err := json.Marshal(message)
	fmt.Println(string(msgBytes))
	h := sha1.New()
	h.Write([]byte(string(msgBytes)))
	msgSha1 := h.Sum(nil)

	return ecdsa.Verify(pubKey, msgSha1, esig.R, esig.S), err
}
