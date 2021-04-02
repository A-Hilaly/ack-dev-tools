// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package git

import (
	"context"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var _ OpenCloner = (*Git)(nil)

const (
	defaultUser = "git"
)

// Cloner is the interface that wraps the Clone method.
//
// Clone clones a remote git repository into a destination path.
type Cloner interface {
	Clone(
		ctx context.Context,
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

// OpenCloner is the interface that wraps the Open and Clone methods.
type OpenCloner interface {
	Opener
	Cloner
}

// New instanciate a new Git struct. remote defaults to 'origin'
// signer can be nil.
func New(signer ssh.Signer) *Git {
	return &Git{
		signer: signer,
	}
}

// Git represents the components reponsible for cloning and
// opening git repositories. It is supposed to hide the authentication
// mechanisms used to clone repositories.
// Git implements OpenCloner interface.
type Git struct {
	signer ssh.Signer
	remote string
}

// Clone clones a remote git repository into a destination path.
func (g *Git) Clone(ctx context.Context, url, dest string) error {
	auth := &gitssh.PublicKeys{
		User: defaultUser,
	}
	if g.signer != nil {
		auth.Signer = g.signer
	}
	_, err := git.PlainCloneContext(ctx, dest, false, &git.CloneOptions{
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
