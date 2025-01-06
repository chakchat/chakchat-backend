package main

import "testing"

func TestMain(m *testing.M) {
	m.Run()
}

func TestPass(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("1 + 1 = 2 but it sayed no")
	}
}
