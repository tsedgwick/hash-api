package api

import (
	"crypto/sha512"
	"encoding/base64"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//Client handles the encoded and saving of passwords
type Client interface {
	Encode(input []byte) string
	Tokenize(input []byte) string
	Token(key string) string
	Save(input []byte) string
}

type client struct {
	mu         sync.RWMutex
	tokens     sync.Map
	keys       map[int]bool
	nextKey    int
	waitEncode time.Duration
}

//New returns a new instance of client
func New() Client {
	return &client{
		keys:       map[int]bool{},
		waitEncode: 5 * time.Second,
	}
}

//Encode takes in the input []byte and encodes to sha512
func (c *client) Encode(input []byte) string {
	i := sha512.Sum512(input)
	return base64.StdEncoding.EncodeToString(i[:])
}

//Tokenize takes in the password and returns a token
func (c *client) Tokenize(input []byte) string {
	time.Sleep(c.waitEncode)
	return c.Encode(input)
}

//Token returns the encoded password
func (c *client) Token(key string) string {
	val, ok := c.tokens.Load(key)
	if !ok {
		return ""
	}
	return val.(string)
}

//Save the encoded password and returns the retrieval key
func (c *client) Save(input []byte) string {
	key := c.newKey()
	go func() {
		c.tokens.Store(key, c.Tokenize(input))
	}()

	return key
}

func (c *client) newKey() string {
	c.mu.Lock()
	key := rand.Intn(10000)
	exists := c.keys[key]
	if exists {
		c.mu.Unlock()
		c.newKey()

	}
	c.keys[key] = true
	c.mu.Unlock()
	return strconv.Itoa(key)
}
