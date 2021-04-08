package ssh

import (
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

var (
	defaultPrivateKeys = []string{
		"id_dsa",
		"id_rsa",
		"id_ecdsa",
		"id_ecdsa_sk",
		"id_ed25519",
		"id_ed25519_sk",
	}
)

func DefaultSigner(parentDir string) []ssh.Signer {
	var signers []ssh.Signer
	for _, pk := range defaultPrivateKeys {
		signer, err := NewSigner(filepath.Join(parentDir, pk), false)
		if err == nil {
			signers = append(signers, signer)
		}
	}
	return signers
}
