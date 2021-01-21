package git

import "testing"

func TestCheckRepo(t *testing.T) {
	sut := DefaultAgent{}
	sut.CheckRepo("https://github.com/redsailtechnologies/boatswain", "", "")
}

func TestGetFile(t *testing.T) {
	sut := DefaultAgent{}
	bytes := sut.GetFile("https://github.com/redsailtechnologies/boatswain", "main", "docs/example.yaml", "", "")
	t.Log(string(bytes))
}
