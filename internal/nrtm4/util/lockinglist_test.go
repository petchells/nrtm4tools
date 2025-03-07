package util

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"math/rand/v2"
)

func TestLockingList(t *testing.T) {

	ll := NewLockingList[string](50)

	printBatch := func(b []string) {
		if len(b) > 0 {
			t.Logf("Batch of %v strings. %v", len(b), strings.Join(b, ","))
		}
	}
	var wg sync.WaitGroup
	listClient := func(pfx string) {
		ll.Add(pfx + "-" + RandStringBytes(10))
		b := ll.GetBatch(15)
		printBatch(b)
	}
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			listClient(fmt.Sprintf("go%02d", i))
		}()
	}
	wg.Wait()
	b := ll.GetAll()
	t.Log("Remainder:")
	printBatch(b)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.IntN(len(letterBytes))]
	}
	return string(b)
}
