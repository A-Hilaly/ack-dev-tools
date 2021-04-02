package testutil

import (
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// NewInMemoryGitRepository returns a in-memory git repository containing one commit.
func NewInMemoryGitRepository() (*git.Repository, error) {
	fs := memfs.New()
	store := memory.NewStorage()
	repo, err := git.Init(store, fs)
	if err != nil {
		return nil, err
	}

	file, err := fs.Create("ramanujan_serie.txt")
	if err != nil {
		return nil, err
	}
	_, err = file.Write([]byte("1 + 2 + 3 + 4 + ... = -1/12"))
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	_, err = w.Add("ramanujan_serie.txt")
	if err != nil {
		return nil, err
	}

	commit, err := w.Commit("first commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Srinivasa Ramanujan",
			Email: "sramanujan@1729",
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = repo.CommitObject(commit)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
