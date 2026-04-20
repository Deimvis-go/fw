package fakefw

import (
	"hash/fnv"
	"math/rand/v2"
)

// NOTE: idea is to use like this:
//
// 	m1 := String("message-1")
// 	m2 := String("message-2")
//  n1 := String("note-1")
//  n2 := String("note-2")
//  // ^ all variables above are different
//
//  Send(m1)
//  ...
//  CheckEqual(db.Fetch(msg1), String("message-1"))
//  // ^ may be used to do not pass context and retrieve value this way
//  // (even though not recommended)

// All above variables are different.

// TODO: support option with length
// TODO: support option with charset
func String(s seed) string {
	const length = 16
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rng := rand.New(pcg(s))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.IntN(len(charset))]
	}
	return string(b)
}

// TODO: support option with length
// TODO: support option with charset
func Bytes(s seed, n int) []byte {
	rng := rand.New(pcg(s))
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rng.IntN(256))
	}
	return b
}

type seed string

func pcg(s seed) *rand.PCG {
	h := fnv.New64a()
	h.Write([]byte(s))
	seed1 := h.Sum64()

	h.Reset()
	h.Write([]byte(s + "salt"))
	seed2 := h.Sum64()

	return rand.NewPCG(seed1, seed2)
}
