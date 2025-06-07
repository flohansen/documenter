package scraper

import (
	"context"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitScraper struct {
	repoURL string
	sshKey  *string
}

func NewGitScraper(repoURL string, opts ...GitScraperOption) *GitScraper {
	gs := &GitScraper{
		repoURL: repoURL,
	}

	for _, opt := range opts {
		opt(gs)
	}

	return gs
}

func (s *GitScraper) Scrape(ctx context.Context) ([]byte, error) {
	cloneOptions, err := s.cloneOptionsForSection()
	if err != nil {
		return nil, fmt.Errorf("clone options setup error: %w", err)
	}

	repo, err := git.CloneContext(ctx, memory.NewStorage(), nil, &cloneOptions)
	if err != nil {
		return nil, fmt.Errorf("clone error: %s", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("repository head error: %s", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("repository commit error: %s", err)
	}

	file, err := commit.File("README.md")
	if err != nil {
		return nil, fmt.Errorf("commit file error: %s", err)
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, fmt.Errorf("file reader error: %s", err)
	}

	b, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		return nil, fmt.Errorf("file reader error: %s", err)
	}

	return b, nil
}

func (s *GitScraper) cloneOptionsForSection() (git.CloneOptions, error) {
	cloneOptions := git.CloneOptions{
		URL:   s.repoURL,
		Depth: 1,
	}

	if s.sshKey != nil {
		publicKeys, err := ssh.NewPublicKeysFromFile("git", *s.sshKey, "")
		if err != nil {
			return git.CloneOptions{}, fmt.Errorf("could not read ssh key: %w", err)
		}

		cloneOptions.Auth = publicKeys
	}

	return cloneOptions, nil
}

type GitScraperOption func(*GitScraper)

func WithSSHKey(sshKey string) GitScraperOption {
	return func(gs *GitScraper) {
		if len(sshKey) > 0 {
			gs.sshKey = &sshKey
		}
	}
}
