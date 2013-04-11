// Package udp implements UDP test helpers. It lets you assert that certain
// strings must or must not be sent to a given local UDP listener.
package udp

import (
	"net"
	"strings"
	"testing"
	"time"
)

var (
	addr     *string
	listener *net.UDPConn
)

type fn func()

// SetAddr sets the UDP port that will be listened on.
func SetAddr(a string) {
	addr = &a
}

func start(t *testing.T) {
	resAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		t.Fatal(err)
	}
	listener, err = net.ListenUDP("udp", resAddr)
	if err != nil {
		t.Fatal(err)
	}
}

func stop(t *testing.T) {
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
}

func getMessage(t *testing.T, body fn) string {
	start(t)
	defer stop(t)

	result := make(chan string)

	go func() {
		message := make([]byte, 1024*32)
		n, _, _ := listener.ReadFrom(message)
		result <- string(message[0:n])
	}()

	body()

	select {
	case text := <-result:
		return text
	case <-time.After(time.Millisecond):
	}

	return ""
}

func get(t *testing.T, match string, body fn) (got string, equals bool, contains bool) {
	got = getMessage(t, body)
	equals = got == match
	contains = strings.Contains(got, match)
	return got, equals, contains
}

// ShouldReceiveOnly will fire a test error if the given function doesn't send
// exactly the given string over UDP.
func ShouldReceiveOnly(t *testing.T, expected string, body fn) {
	got, equals, _ := get(t, expected, body)
	if !equals {
		t.Errorf("Expected %#v but got %#v instead", expected, got)
	}
}

// ShouldNotReceiveOnly will fire a test error if the given function sends
// exactly the given string over UDP.
func ShouldNotReceiveOnly(t *testing.T, notExpected string, body fn) {
	_, equals, _ := get(t, notExpected, body)
	if equals {
		t.Errorf("Expected not to get %v but did", notExpected)
	}
}

// ShouldReceive will fire a test error if the given function doesn't send the
// given string over UDP.
func ShouldReceive(t *testing.T, expected string, body fn) {
	got, _, contains := get(t, expected, body)
	if !contains {
		t.Errorf("Expected to find %#v but got %#v instead", expected, got)
	}
}

// ShouldNotReceive will fire a test error if the given function sends the
// given string over UDP.
func ShouldNotReceive(t *testing.T, expected string, body fn) {
	got, _, contains := get(t, expected, body)
	if contains {
		t.Errorf("Expected not to find %#v but got %#v", expected, got)
	}
}

// ShouldReceiveAll will fire a test error unless all of the given strings are
// sent over UDP.
func ShouldReceiveAll(t *testing.T, expected []string, body fn) {
	got := getMessage(t, body)
	for _, str := range expected {
		if !strings.Contains(got, str) {
			t.Errorf("Expected to find %#v but got %#v instead", str, got)
		}
	}
}

// ShouldNotReceiveAny will fire a test error if any of the given strings are
// sent over UDP.
func ShouldNotReceiveAny(t *testing.T, unexpected []string, body fn) {
	got := getMessage(t, body)
	for _, str := range unexpected {
		if strings.Contains(got, str) {
			t.Errorf("Expected not to find %#v but got %#v", str, got)
		}
	}
}

func ShouldReceiveAllAndNotReceiveAny(t *testing.T, expected []string, unexpected []string, body fn) {
	got := getMessage(t, body)
	for _, str := range expected {
		if !strings.Contains(got, str) {
			t.Errorf("Expected to find %#v but got %#v instead", str, got)
		}
	}
	for _, str := range unexpected {
		if strings.Contains(got, str) {
			t.Errorf("Expected not to find %#v but got %#v", str, got)
		}
	}
}
