package mysql

import (
	"crypto/rand"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)

	chacha, err := chacha20poly1305.NewX(key)
	assert.NoError(t, err)

	testCases := []struct {
		entryID, docTypeID uint64
		docName, plainText string
	}{
		{entryID: 0, docTypeID: 1, docName: "test", plainText: "{}"},
		{entryID: 0xdeadbeef, docTypeID: ^uint64(0), docName: "a", plainText: "null"},
		{entryID: 255 << 56, docTypeID: 255, docName: "long-document-name", plainText: "{\"age\":100}"},
	}

	for _, test := range testCases {
		entryID := test.entryID

		// encrypt
		cipherText, err := encryptDoc(chacha, entryID, test.docTypeID, test.docName, []byte(test.plainText))
		assert.NoError(t, err)
		assert.Equal(t, chacha.NonceSize()+len(test.plainText)+chacha.Overhead(), len(cipherText))
		assert.Equal(t, entryID, binary.LittleEndian.Uint64(cipherText))

		// changed entry id
		_, err = decryptDoc(chacha, entryID+1, test.docTypeID, test.docName, cipherText)
		assert.Error(t, err)

		// successfull decrypt
		plainText, err := decryptDoc(chacha, entryID, test.docTypeID, test.docName, cipherText)
		assert.NoError(t, err)
		assert.Equal(t, []byte(test.plainText), plainText)

		// changed docType
		_, err = decryptDoc(chacha, entryID, test.docTypeID+1, test.docName, cipherText)
		assert.Error(t, err)

		// changed docName
		_, err = decryptDoc(chacha, entryID, test.docTypeID, strings.ToTitle(test.docName), cipherText)
		assert.Error(t, err)
	}
}
