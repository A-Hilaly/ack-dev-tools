package ssh

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func promptPassphrase() ([]byte, error) {
	fmt.Printf("type your ssh key passphrase")
	passphrase, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	// TODO(hilalymh): retry + validation mechanisms.
	return passphrase, nil
}

// NewSigner returns a new ssh.Signer. If the PEM file is encrypted
// it will try to read the passphrase from a terminal without local echo.
func NewSigner(sshKeyPath string) (ssh.Signer, error) {
	pemBytes, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid ssh certificate")
	}

	if x509.IsEncryptedPEMBlock(block) {
		passphrase, err := promptPassphrase()
		if err != nil {
			return nil, err
		}
		return ssh.ParsePrivateKeyWithPassphrase(pemBytes, passphrase)
	}

	return ssh.ParsePrivateKey(pemBytes)
}
