package git

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var _ OpenCloner = &Git{}

const (
	defaultUser = "git"
)

// Cloner is the interface that wraps the Clone method.
//
// Clone clones a remote git repository into a destination path.
type Cloner interface {
	Clone(
		url string,
		dest string,
	) error
}

// Open is the interface that wraps the Open method.
//
// Open opens a git repository from the given path.
type Opener interface {
	Open(path string) (*git.Repository, error)
}

// Open is the interface that wraps the Open and Clone methods.
type OpenCloner interface {
	Opener
	Cloner
}

// New instanciate a new Git struct. remote defaults to 'origin'
// signer can be nil.
func New(remote string, signer ssh.Signer) *Git {
	return &Git{
		signer: signer,
		remote: remote,
	}
}

// Git represents the components reponsible for cloning and
// opening git repositories.
type Git struct {
	signer ssh.Signer
	remote string
}

// Clone clones a remote git repository into a destination path.
func (g *Git) Clone(url, dest string) error {
	auth := &gitssh.PublicKeys{
		User: defaultUser,
	}
	if g.signer != nil {
		auth.Signer = g.signer
	}
	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		Auth:       auth,
		URL:        url,
		RemoteName: g.remote,
		Progress:   nil,
	})
	if err != nil {
		return err
	}
	return nil
}

// Open opens a git repository from the given path.
func (g *Git) Open(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}
