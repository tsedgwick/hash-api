package api

import (
	"crypto/sha512"
	"encoding/base64"
	"sync"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected func() string
	}{
		{
			name:  "Success",
			input: []byte("angryMonkey"),
			expected: func() string {
				i := sha512.Sum512([]byte("angryMonkey"))
				return base64.StdEncoding.EncodeToString(i[:])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				keys:       map[int]bool{},
				waitEncode: 5 * time.Millisecond,
			}
			actual := c.Encode(tt.input)
			if actual != tt.expected() {
				t.Errorf("failed at %s : expected %v : actual : %v", tt.name, tt.expected(), actual)
			}
		})
	}
}

func TestToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		save     func(c *client)
	}{
		{
			name:  "Success",
			input: "123",
			save: func(c *client) {
				c.tokens.Store("123", "successToken")
			},
			expected: "successToken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				keys: map[int]bool{},
			}
			tt.save(c)
			actual := c.Token(tt.input)
			if actual != tt.expected {
				t.Errorf("failed at %s : expected %v : actual : %v", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestConcurrentSavePassword(t *testing.T) {
	c := &client{keys: map[int]bool{}}

	input := []string{
		"angryMonkey",
		"funkyMonkey",
		"sadMonkey",
		"happyMonkey",
		"justMonkey",
	}

	ch := make(chan (string))
	var wg sync.WaitGroup
	wg.Add(len(input))
	for _, v := range input {
		go func(input string, ch chan (string)) {
			key := c.Save([]byte(input))
			_, ok := c.tokens.Load(key)
			if ok {
				t.Errorf("failed at %s : expected empty value in cache", key)
			}
			ch <- key
			wg.Done()
		}(v, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	time.Sleep(1 * time.Second)
	for key := range ch {
		token, ok := c.tokens.Load(key)
		if !ok {
			t.Errorf("failed at %s : expected value in cache", key)
		}
		t.Logf("Key %s : Token : %s", key, token)
	}

}
