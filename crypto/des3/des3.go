/*
	Package des3 supports password encrypt/decrypt
	It depends on 3DES
	Support by 张松
*/
package des3

import (
	"encoding/json"
	"sync"

	"github.com/bluele/gcache"
	"github.com/pkg/errors"
	curl "gitlab.ifchange.com/data/cordwood/rpc/rpc-curl"
)

var (
	constUrl       string = ""
	constWithCache bool   = false

	setupLock sync.RWMutex
)

func Setup(url string, withCache bool) {
	setupLock.Lock()
	defer setupLock.Unlock()
	constUrl = url
	constWithCache = withCache
}

var (
	encryptCache = gcache.New(2000).
		LRU().
		// Expiration(time.Minute).
		LoaderFunc(func(plaintextInterface interface{}) (interface{}, error) {
			plaintext, ok := plaintextInterface.(string)
			if !ok {
				return nil, errors.Errorf("Encrypt GCache Loader %s is Not String",
					plaintextInterface)
			}
			ciphertext, err := encrypt(plaintext)
			if err != nil {
				return nil, errors.Wrap(err, "encrypt")
			}
			return ciphertext, nil
		}).
		Build()
)

func Encrypt(plaintext string) (ciphertext string, err error) {
	setupLock.RLock()
	defer setupLock.RUnlock()

	if !constWithCache {
		return encrypt(plaintext)
	}
	ciphertextInterface, err := encryptCache.Get(plaintext)
	if err != nil {
		return "", errors.Wrap(err, "encryptCache.Get")
	}
	ciphertext, ok := ciphertextInterface.(string)
	if !ok {
		return "", errors.Errorf("Encrypt GCache Get %s is Not String",
			ciphertextInterface)
	}
	return ciphertext, nil
}

func encrypt(plaintext string) (ciphertext string, err error) {
	req := curl.NewRequest()
	req.SetC("/security")
	req.SetM("/encrypt")
	req.SetP(map[string]interface{}{"plaintext": []string{plaintext}})
	rsp := curl.NewResponse()
	if err := curl.Curl(constUrl, req, rsp); err != nil {
		return "", errors.Wrap(err, "CURL")
	}
	if errNo := rsp.GetErrNo(); errNo != 0 {
		return "", errors.Errorf("encrypt ErrNo:%d ErrMsg:%s",
			errNo, rsp.GetErrMsg())
	}
	ciphertextes := make(map[string]string)
	err = json.Unmarshal(rsp.GetResults(), &ciphertextes)
	if err != nil {
		return "", errors.Wrap(err, "encrypt Unmarshal")
	}
	if len(ciphertextes) == 0 {
		return "", errors.Errorf("encrypt Reply is Empty Array")
	}
	ciphertext = ciphertextes[plaintext]
	if len(ciphertext) == 0 {
		return "", errors.Errorf("encrypt Reply is Empty")
	}
	return ciphertext, nil
}

var (
	decryptCache = gcache.New(2000).
		LRU().
		// Expiration(time.Minute).
		LoaderFunc(func(ciphertextInterface interface{}) (interface{}, error) {
			ciphertext, ok := ciphertextInterface.(string)
			if !ok {
				return nil, errors.Errorf("Decrypt GCache Loader %s is Not String",
					ciphertextInterface)
			}
			plaintext, err := decrypt(ciphertext)
			if err != nil {
				return nil, errors.Wrap(err, "decrypt")
			}
			return plaintext, nil
		}).
		Build()
)

func Decrypt(ciphertext string) (plaintext string, err error) {
	setupLock.RLock()
	defer setupLock.RUnlock()

	if !constWithCache {
		return decrypt(ciphertext)
	}
	plaintextInterface, err := decryptCache.Get(ciphertext)
	if err != nil {
		return "", errors.Wrap(err, "decryptCache.Get")
	}
	plaintext, ok := plaintextInterface.(string)
	if !ok {
		return "", errors.Errorf("Decrypt GCache Get %s is Not String",
			plaintextInterface)
	}
	return plaintext, nil
}

func decrypt(ciphertext string) (plaintext string, err error) {
	req := curl.NewRequest()
	req.SetC("/security")
	req.SetM("/decrypt")
	req.SetP(map[string]interface{}{"ciphertext": []string{ciphertext}})
	rsp := curl.NewResponse()
	if err := curl.Curl(constUrl, req, rsp); err != nil {
		return "", errors.Wrap(err, "CURL")
	}
	if errNo := rsp.GetErrNo(); errNo != 0 {
		return "", errors.Errorf("decrypt ErrNo:%d ErrMsg:%s",
			errNo, rsp.GetErrMsg())
	}
	plaintextes := make(map[string]string)
	err = json.Unmarshal(rsp.GetResults(), &plaintextes)
	if err != nil {
		return "", errors.Wrap(err, "decrypt Unmarshal")
	}
	if len(plaintextes) == 0 {
		return "", errors.Errorf("decrypt Reply is Empty Array")
	}
	plaintext = plaintextes[ciphertext]
	if len(plaintext) == 0 {
		return "", errors.Errorf("decrypt Reply is Empty")
	}
	return plaintext, nil
}
