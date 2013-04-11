// Package udp implements UDP test helpers. It lets you assert that certain
// strings must or must not be sent to a given local UDP listener.
package udp

import (
	"net"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	addr     *string
	listener *net.UDPConn
	Timeout  time.Duration = time.Millisecond
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
		listener.SetReadDeadline(time.Now().Add(Timeout))
		n, _, _ := listener.ReadFrom(message)
		result <- string(message[0:n])
	}()

	body()

	return <-result
}

func get(t *testing.T, match string, body fn) (got string, equals bool, contains bool) {
	got = getMessage(t, body)
	equals = got == match
	contains = strings.Contains(got, match)
	return got, equals, contains
}

func printLocation(t *testing.T) {
	_, file, line, _ := runtime.Caller(2)
	t.Errorf("At: %s:%d", file, line)
}

// ShouldReceiveOnly will fire a test error if the given function doesn't send
// exactly the given string over UDP.
func ShouldReceiveOnly(t *testing.T, expected string, body fn) {
	got, equals, _ := get(t, expected, body)
	if !equals {
		printLocation(t)
		t.Errorf("Expected: %s", expected)
		t.Errorf("But got: %s", got)
	}
}

// ShouldNotReceiveOnly will fire a test error if the given function sends
// exactly the given string over UDP.
func ShouldNotReceiveOnly(t *testing.T, notExpected string, body fn) {
	_, equals, _ := get(t, notExpected, body)
	if equals {
		printLocation(t)
		t.Errorf("Expected not to get: %s", notExpected)
	}
}

// ShouldReceive will fire a test error if the given function doesn't send the
// given string over UDP.
func ShouldReceive(t *testing.T, expected string, body fn) {
	got, _, contains := get(t, expected, body)
	if !contains {
		printLocation(t)
		t.Errorf("Expected to find: %s", expected)
		t.Errorf("But got: %s", got)
	}
}

// ShouldNotReceive will fire a test error if the given function sends the
// given string over UDP.
func ShouldNotReceive(t *testing.T, expected string, body fn) {
	got, _, contains := get(t, expected, body)
	if contains {
		printLocation(t)
		t.Errorf("Expected not to find: %s", expected)
		t.Errorf("But got: %s", got)
	}
}

// ShouldReceiveAll will fire a test error unless all of the given strings are
// sent over UDP.
func ShouldReceiveAll(t *testing.T, expected []string, body fn) {
	got := getMessage(t, body)
	failed := false

	for _, str := range expected {
		if !strings.Contains(got, str) {
			if !failed {
				printLocation(t)
				failed = true
			}
			t.Errorf("Expected to find: %s", str)
		}
	}

	if failed {
		t.Errorf("But got: %s", got)
	}
}

// ShouldNotReceiveAny will fire a test error if any of the given strings are
// sent over UDP.
func ShouldNotReceiveAny(t *testing.T, unexpected []string, body fn) {
	got := getMessage(t, body)
	failed := false

	for _, str := range unexpected {
		if strings.Contains(got, str) {
			if !failed {
				printLocation(t)
				failed = true
			}
			t.Errorf("Expected not to find: %s", str)
		}
	}

	if failed {
		t.Errorf("But got: %s", got)
	}
}

func ShouldReceiveAllAndNotReceiveAny(t *testing.T, expected []string, unexpected []string, body fn) {
	got := getMessage(t, body)
	failed := false

	for _, str := range expected {
		if !strings.Contains(got, str) {
			if !failed {
				printLocation(t)
				failed = true
			}
			t.Errorf("Expected to find: %s", str)
		}
	}
	for _, str := range unexpected {
		if strings.Contains(got, str) {
			if !failed {
				printLocation(t)
				failed = true
			}
			t.Errorf("Expected not to find: %s", str)
		}
	}

	if failed {
		t.Errorf("but got: %s", got)
	}
}
