package repository

import (
	ackdevgit "github.com/aws-controllers-k8s/dev-tools/pkg/git"
	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
)

type Option func(m *Manager)

func WithSigner(privateKeyPath string) Option {
	return func(m *Manager) {
		signer, err := util.NewSigner(privateKeyPath)
		if err != nil {
			panic(err)
		}
		m.git = ackdevgit.New(signer)
	}
}
