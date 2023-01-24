package conf

import "testing"

func TestGetNodes(t *testing.T) {
	if err := LocalConfig(); err != nil {
		t.Fatal(err)
	}

	nodes := GetNodes()
	for name, rpcAddr := range nodes {
		t.Log(name, rpcAddr)
	}
}
