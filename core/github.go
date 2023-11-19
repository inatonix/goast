package core

type GithubManipulator interface {
}

type initParams struct {
	Model string
}

type githubManipulator struct {
}

func NewGithubManipulator() (*githubManipulator, error) {
	return nil, nil
}
