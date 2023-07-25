package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	Charset             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	MinLengthOfStrings  = 2
	MaxLengthOfStrings  = 20000
	CountOfEnoughHashes = 498
)

// A function to generate random strings and send them to a channel
func randomStringGenerator(ch chan string) {
	for {
		length := rand.Intn(MaxLengthOfStrings-MinLengthOfStrings+1) + MinLengthOfStrings
		b := make([]byte, length)

		for i := range b {
			b[i] = Charset[rand.Intn(len(Charset))]
		}

		ch <- string(b)
	}
}

func hashString(s string) string {
	h := sha256.Sum256([]byte(s))

	return hex.EncodeToString(h[:])
}

func sumHash(h string) int {
	sum := 0

	for _, c := range h {
		if c >= '0' && c <= '9' {
			n, _ := strconv.Atoi(string(c))
			sum += n
		}
	}

	return sum
}

func checkAndAppendHash(s string, hashes *[]hashSum, mu *sync.Mutex) {
	h := hashString(s)

	if strings.HasSuffix(h, "000") {
		fmt.Println("Found a hash ending with 000:", h)
		sum := sumHash(h)

		mu.Lock()
		*hashes = append(*hashes, hashSum{
			hash: h,
			sum:  sum,
		})
		mu.Unlock()
	}
}

type hashSum struct {
	hash string
	sum  int
}

// A type to implement sort.Interface for []hashSum based on the sum field
type bySum []hashSum

func (s bySum) Len() int {
	return len(s)
}

func (s bySum) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s bySum) Less(i, j int) bool {
	return s[i].sum < s[j].sum
}

func main() {
	hashes := make([]hashSum, 0)

	wg := sync.WaitGroup{}

	mu := sync.Mutex{}

	ch := make(chan string)

	go randomStringGenerator(ch)

	for {

		// Receive a random string from the channel
		s := <-ch

		wg.Add(1)

		go func(s string) {
			defer wg.Done()

			checkAndAppendHash(s, &hashes, &mu)

			//This is called an immediately-invoked function expression (IIFE)
		}(s)

		//Break the loop if we have enough hashes
		if len(hashes) >= CountOfEnoughHashes {
			break
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All goroutines finished")

	fmt.Println("Unsorted hashes:")

	for _, hs := range hashes {
		fmt.Printf("%s (sum = %d)\n", hs.hash, hs.sum)
	}

	fmt.Println("Sorted hashes:")
	sort.Sort(bySum(hashes))

	for _, hs := range hashes {
		fmt.Printf("%s (sum = %d)\n", hs.hash, hs.sum)
	}
}
