package main

import (
	"crypto"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	osutil "github.com/snapcore/snapd/osutil"

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
	Algorithm string `json:"alg"`
	BytesHash int64  `json:"bytes"`
	FileHash  int64  `json:"file"`
}

func main() {
	// Get the file from the flag
	fileStr := flag.String("file", "", "file to hash")

	randomMode := flag.Bool("random", true, "whether to generate a random file or not")

	randomSizeMB := flag.Int("size", 10, "size of generated random file")

	unitStr := flag.String("unit", "s", "units to use (possible values : ns, us, ms, s)")

	formatStr := flag.String("format", "user", "format to output, user, json, or yaml (not implemented yet)")

	algStr := flag.String("alg", "sha3_512", "algorithm to use, sha3_512, sha3_384, sha3_256, or sha3_224")

	// Parse command line flags
	flag.Parse()

	// Check the units to use
	var timeVal time.Duration
	switch strings.ToLower(*unitStr) {
	case "ns":
		// use Nanoseconds
		timeVal = time.Nanosecond
	case "us":
		timeVal = time.Microsecond
	case "ms":
		timeVal = time.Millisecond
	case "s":
		timeVal = time.Second
	default:
		fmt.Println("wrong value for units")
		os.Exit(1)
	}

	var fileExistsQ bool
	if _, err := os.Stat(*fileStr); os.IsNotExist(err) {
		fileExistsQ = false
	} else {
		fileExistsQ = true
	}

	if *randomMode && !fileExistsQ {
		// then don't use a file - generate one randomly
		size := (*randomSizeMB) * 1048576
		randomBytes := make([]byte, size)
		_, err := rand.Read(randomBytes)

		// now write the file out
		tmpfile, err := ioutil.TempFile("", "sha3sum_example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(randomBytes); err != nil {
			log.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}

		var tempfileName string
		tempfileName = tmpfile.Name()
		fileStr = &tempfileName
	}

	fmt.Printf("file is %s\n", *fileStr)

	switch strings.ToLower(*formatStr) {
	case "user":
		// run a single run and print out human readable form
		// check what algorithm to use
		var hasherToUse crypto.Hash
		switch strings.ToLower(*algStr) {
		case "sha3_512":
			hasherToUse = crypto.SHA3_512
		case "sha3_384":
			hasherToUse = crypto.SHA3_384
		case "sha3_256":
			hasherToUse = crypto.SHA3_256
		case "sha3_224":
			hasherToUse = crypto.SHA3_224
		default:
			fmt.Println("not implemented yet")
			os.Exit(1)
		}
		_, timeElapsedFile := timeFileHash(hasherToUse, *fileStr)
		fmt.Printf("system: %3f%s\n", float64(timeElapsedFile)/float64(timeVal), *unitStr)
	case "yaml":
		fmt.Println("not implemented yet")
		os.Exit(1)
	case "json":
		// run the normal test suite and then output it as json
		runJSON(initializeHasherTable(), *fileStr, timeVal)
	default:
		fmt.Println("wrong value for format")
		os.Exit(1)
	}
}

func runJSON(hasherTable []HasherType, file string, timeVal time.Duration) {
	// Read in the file as a byte array
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error reading a file : %+v\n", err)
	}

	var results []HashResult
	results = make([]HashResult, len(hasherTable))

	// For each hasher run the hash and print off how long it took
	for index, hasher := range hasherTable {
		// Compute the hash for each one
		_, timeElapsedBytes := timeBytesHash(hasher.Hasher, bytes)
		_, timeElapsedFile := timeFileHash(hasher.HashType, file)

		results[index] = HashResult{
			Algorithm: hasher.Name,
			BytesHash: int64(timeElapsedBytes / timeVal),
			FileHash:  int64(timeElapsedFile / timeVal),
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
		{"sha3_224", sha3.New224(), crypto.SHA3_224},
		{"sha3_256", sha3.New256(), crypto.SHA3_256},
		{"sha3_384", sha3.New384(), crypto.SHA3_384},
		{"sha3_512", sha3.New512(), crypto.SHA3_512},
	}
}
