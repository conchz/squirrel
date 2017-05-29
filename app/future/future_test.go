package future

import "testing"

func TestSimpleFuture(t *testing.T) {
	future := SimpleFuture("https://api.github.com/users/octocat/orgs")

	// not block
	t.Log("Not block")

	body, err := future() // block
	if err == nil {
		t.Logf("The response length is %d\n", len(body))
	} else {
		t.Errorf("Error occurred: %v\n", err)
	}
}
