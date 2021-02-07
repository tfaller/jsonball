package mysql

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"time"
)

func encryptDoc(aead cipher.AEAD, id uint64, docType uint64, docName string, docData []byte) ([]byte, error) {
	docNameLen := len(docName)

	// Create dst buffer that contains additional auth values, nonce, encrypted data, auth tag.
	dst := make([]byte, 8+docNameLen, len(docData)+32+aead.Overhead()+docNameLen)

	// Prepare additional auth info.
	// Use docType, docName and entry id as additional authentication,
	// this prevents that an attacker can swap data to a different document.
	authInfo := dst[:docNameLen+16]
	binary.LittleEndian.PutUint64(authInfo, docType)
	copy(authInfo[8:], docName)
	binary.LittleEndian.PutUint64(authInfo[8+docNameLen:], id)

	// cipherText will contain a leading nonce,
	// the encrypted data and the auth tag in this order
	cipherText := dst[8+docNameLen:]

	// The nonce will contain the entry id, current time and random data.
	// This makes sure that a nonce will not be reused.
	// entry id is already set.
	nonce := cipherText[:24]

	// add current time to nonce
	binary.LittleEndian.PutUint64(nonce[8:], uint64(time.Now().UnixNano()))

	// add random data to nonce
	_, err := rand.Read(nonce[16:])
	if err != nil {
		return nil, err
	}

	// encrypt and tag data
	return aead.Seal(nonce, nonce, docData, authInfo), nil
}

func decryptDoc(aead cipher.AEAD, id uint64, docType uint64, docName string, data []byte) ([]byte, error) {
	docNameLen := len(docName)
	dataLen := len(data) - aead.NonceSize() - aead.Overhead()

	// create a buffer to hold auth info and decrypted data
	dst := make([]byte, docNameLen+16+dataLen)
	authInfo := dst[dataLen:]

	binary.LittleEndian.PutUint64(authInfo, docType)
	copy(authInfo[8:], docName)
	binary.LittleEndian.PutUint64(authInfo[8+docNameLen:], id)

	return aead.Open(dst[:0], data[:24], data[24:], authInfo)
}
