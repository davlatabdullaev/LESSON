package main

import "testing"



func TestHello(t *testing.T) {
  
	msg := Hello()
	want := "Hello World!"
	if want != msg {
		t.Errorf("expected %q, but got: %q", want, msg)
	}

}