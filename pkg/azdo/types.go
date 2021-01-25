package azdo

// GitRepo represents an azdo git repo
type GitRepo struct {
	ID   string
	Name string
}

// PullRequest represents an azdo pr
type PullRequest struct {
	ID    int
	Title string
}
