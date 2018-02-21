package main

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"strings"
	"time"

	osutil "github.com/snapcore/snapd/osutil"

	md4 "golang.org/x/crypto/md4"
	ripemd160 "golang.org/x/crypto/ripemd160"
	sha3 "golang.org/x/crypto/sha3"
)

func timeBytesHash(hasher hash.Hash, bytes []byte) ([]byte, time.Duration) {
	// Get start time
	start := time.Now()

	// Compute the hash of the bytes
	hashRes := hasher.Sum(bytes)

	// Return the hash result and the time since the start of the function
	return hashRes, time.Since(start)
}

func timeFileHash(hasher crypto.Hash, file string) ([]byte, time.Duration) {
	// Get start time
	start := time.Now()

	// Compute the hash of the file
	hashRes, _, _ := osutil.FileDigest(file, hasher)

	// Return the hash result and the time since the start of the function
	return hashRes, time.Since(start)
}

type HasherType struct {
	Name     string
	Hasher   hash.Hash
	HashType crypto.Hash
}

type HashResult struct {
	Algorithm   string `json:"alg"`
	BytesHashNs int64  `json:"bytes_ns"`
	FileHashNs  int64  `json:"file_ns"`
}

func main() {
	// Get the file from the flag
	fileStr := flag.String("file", "file.txt", "file to hash")

	// Parse command line flags
	flag.Parse()

	// Make all the hasher objects
	hasherTable := initializeHasherTable()

	// Read in the file as a byte array
	bytes, err := ioutil.ReadFile(*fileStr)
	if err != nil {
		log.Fatalf("Error reading a file : %+v\n", err)
	}

	var results []HashResult
	results = make([]HashResult, len(hasherTable))

	// For each hasher run the hash and print off how long it took
	for index, hasher := range hasherTable {
		// Compute the hash for each one
		_, timeElapsedBytes := timeBytesHash(hasher.Hasher, bytes)
		_, timeElapsedFile := timeFileHash(hasher.HashType, *fileStr)

		results[index] = HashResult{
			Algorithm:   hasher.Name,
			BytesHashNs: int64(timeElapsedBytes / time.Nanosecond),
			FileHashNs:  int64(timeElapsedFile / time.Nanosecond),
		}
	}

	jsonbytes, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("failed to encode json")
	}

	fmt.Println(strings.Replace(string(jsonbytes[:]), "\"", "\\\"", -1))
}

func initializeHasherTable() []HasherType {
	return []HasherType{
		{"md4", md4.New(), crypto.MD4},
		{"md5", md5.New(), crypto.MD5},
		{"sha1", sha1.New(), crypto.SHA1},
		{"sha256", sha256.New(), crypto.SHA256},
		{"sha256_224", sha256.New224(), crypto.SHA224},
		{"sha512", sha512.New(), crypto.SHA512},
		{"sha512_224", sha512.New512_224(), crypto.SHA512_224},
		{"sha512_256", sha512.New512_256(), crypto.SHA512_256},
		{"sha384", sha512.New384(), crypto.SHA384},
		{"ripemd160", ripemd160.New(), crypto.RIPEMD160},
		{"sha3_224", sha3.New224(), crypto.SHA3_224},
		{"sha3_256", sha3.New256(), crypto.SHA3_256},
		{"sha3_384", sha3.New384(), crypto.SHA3_384},
		{"sha3_512", sha3.New512(), crypto.SHA3_512},
		// Disabled for now...
		// {"blake2s128", blake2s.New128(nil)},
		// {"blake2s256", blake2s.New256(nil)},
		// {"blake2b", blake2b.New(128, nil)},
		// {"blake2b256", blake2b.New256(nil)},
		// {"blake2b384", blake2b.New384(nil)},
		// {"blake2b512", blake2b.New512(nil)},
	}
}
