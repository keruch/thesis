package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"golang.org/x/crypto/hkdf"
)

type SecurityLayer struct {
	privateKey crypto.PrivKey
}

func NewSecurityLayer(privateKey crypto.PrivKey) *SecurityLayer {
	return &SecurityLayer{
		privateKey: privateKey,
	}
}

func (sl *SecurityLayer) EncryptMessage(msg []byte, recipientPubKey crypto.PubKey) ([]byte, error) {
	sharedSecret, err := sl.generateSharedSecret(recipientPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shared secret: %w", err)
	}

	key := sl.deriveKey(sharedSecret, []byte("encryption"))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, msg, nil)
	return ciphertext, nil
}

func (sl *SecurityLayer) DecryptMessage(ciphertext []byte, senderPubKey crypto.PubKey) ([]byte, error) {
	sharedSecret, err := sl.generateSharedSecret(senderPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shared secret: %w", err)
	}

	key := sl.deriveKey(sharedSecret, []byte("encryption"))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (sl *SecurityLayer) SignMessage(msg []byte) ([]byte, error) {
	return sl.privateKey.Sign(msg)
}

func (sl *SecurityLayer) VerifySignature(msg, signature []byte, pubKey crypto.PubKey) (bool, error) {
	return pubKey.Verify(msg, signature)
}

func (sl *SecurityLayer) generateSharedSecret(peerPubKey crypto.PubKey) ([]byte, error) {
	//ecdheSharedSecret, err := sl.privateKey.(*crypto.ECDSAPrivateKey).GenerateShared(peerPubKey.(*crypto.ECDSAPublicKey))
	//if err != nil {
	//	return nil, fmt.Errorf("failed to generate ECDHE shared secret: %w", err)
	//}

	return nil, nil
}

func (sl *SecurityLayer) deriveKey(secret, info []byte) []byte {
	hash := sha256.New
	hkdf := hkdf.New(hash, secret, nil, info)
	key := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		panic(err)
	}
	return key
}

func (sl *SecurityLayer) GetPublicKey() crypto.PubKey {
	return sl.privateKey.GetPublic()
}

func (sl *SecurityLayer) GetPeerID() peer.ID {
	id, err := peer.IDFromPrivateKey(sl.privateKey)
	if err != nil {
		panic(err)
	}
	return id
}
