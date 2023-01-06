package qfsign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"

	"github.com/pkg/errors"
)

type Handler interface {
	Sign(data []byte) (string, error)
	Verify(data []byte, base64Sign string) error
}

type ConfigParams struct {
	Algorithm          AlgorithmType
	PublicKeyCertPath  string
	PrivateKeyCertPath string
}

type AlgorithmType string

const (
	RSASign AlgorithmType = "rsa"
)

var (
	defaultRSAHandler = &rsaHandler{}
)

func NewHandler(c ConfigParams) (Handler, error) {
	switch c.Algorithm {
	case RSASign:
		return defaultRSAHandler.SetPublicKeyPath(c.PublicKeyCertPath).SetPriVateKeyPath(c.PrivateKeyCertPath), nil
	default:
		return nil, errors.Errorf("unsupported sign algorithm type: %s", c.Algorithm)
	}
}

type rsaHandler struct {
	Config ConfigParams
}

func (h *rsaHandler) SetPublicKeyPath(pubKeyPath string) *rsaHandler {
	h.Config.PublicKeyCertPath = pubKeyPath
	return h
}

func (h *rsaHandler) SetPriVateKeyPath(priKeyPath string) *rsaHandler {
	h.Config.PrivateKeyCertPath = priKeyPath
	return h
}

func (h *rsaHandler) Sign(data []byte) (string, error) {
	if h.Config.PrivateKeyCertPath == "" {
		return "", errors.New("private key path is empty, can't generate sign")
	}

	// 1.select sha256 as algorithm of hasg, and generate sign for the data
	myhash := crypto.SHA256
	hashInstance := myhash.New()
	hashInstance.Write(data)
	hashed := hashInstance.Sum(nil)

	// 2.read private key file, parse the key object
	privateKey, err := h.readParsePrivateKey()
	if err != nil {
		return "", err
	}

	// 3.rsa sing(param is random, private key, hash of the file to be signed)
	bytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, myhash, hashed)
	if err != nil {
		return "", err
	}

	// 4.use base64 to encode bytes
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (h *rsaHandler) Verify(data []byte, base64Sign string) error {
	if h.Config.PublicKeyCertPath == "" {
		return errors.New("public key path is empty, can't verify sign")
	}

	// 1.the Base64-encoded signature content is decoded to return signature bytes
	bytes, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}

	// 2.select sha256 as algorithm of hasg, and generate sign for the data
	myhash := crypto.SHA256
	hashInstance := myhash.New()
	hashInstance.Write(data)
	hashed := hashInstance.Sum(nil)

	// 3、read the public key file and parse out the public key object
	publicKey, err := h.readParsePublicKey()
	if err != nil {
		return err
	}
	// 4. rsa verifies  ign (param are public key object, hash type,
	// hash string of sign file, and bytes after sign)
	return rsa.VerifyPKCS1v15(publicKey, myhash, hashed, bytes)
}

func (h *rsaHandler) readParsePublicKey() (*rsa.PublicKey, error) {
	publicKeyBytes, err := ioutil.ReadFile(h.Config.PublicKeyCertPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("public key info err")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}

func (h *rsaHandler) readParsePrivateKey() (*rsa.PrivateKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(h.Config.PrivateKeyCertPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("privte key info err")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
