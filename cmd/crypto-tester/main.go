package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
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
	md4 "golang.org/x/crypto/md4"
	ripemd160 "golang.org/x/crypto/ripemd160"
	sha3 "golang.org/x/crypto/sha3"
	yaml "gopkg.in/yaml.v2"
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
	fileStr := flag.String("file", "", "file to hash (random file is generated if this is omitted)")
	randomSizeMB := flag.Int64("size", 10, "size of generated random file")
	unitStr := flag.String("unit", "ns", "units to use (possible values : ns, us, ms, s)")
	escapteQuotes := flag.Bool("escape-quote", false, "whether to escape quotes in the json output")
	formatStr := flag.String("format", "yaml", "format to output (json or yaml)")

	// Parse command line flags
	flag.Parse()

	// Check the units to use for output
	var timeVal time.Duration
	switch strings.ToLower(*unitStr) {
	case "ns":
		timeVal = time.Nanosecond
	case "us":
		timeVal = time.Microsecond
	case "ms":
		timeVal = time.Millisecond
	case "s":
		timeVal = time.Second
	default:
		log.Fatalf("error : invalid units specification %s\n", *unitStr)
	}

	// Check whether the file exists or not, if it doesn't that might be okay as we
	// might be generating a random file
	var fileExistsQ bool
	var fileSize int64
	if file, err := os.Stat(*fileStr); os.IsNotExist(err) {
		fileExistsQ = false
	} else {
		fileExistsQ = true
		// Also save the file size now that we have a file that exists
		fileSize = file.Size()
	}

	// Now check the type of file handling
	switch {
	case !fileExistsQ && *fileStr != "":
		// File was specified but doesn't exist - we can use err as it won't have been cleared yet
		log.Fatalf("error : file %s doesn't exist\n", *fileStr)
	case !fileExistsQ:
		// then don't use a file - generate one randomly
		fileSize = (*randomSizeMB) * 1048576
		randomBytes := make([]byte, fileSize)

		// Read this many bytes from the OS's random number generator
		_, err := rand.Read(randomBytes)

		// Make a new temp file
		tmpfile, err := ioutil.TempFile("", "crypto_tester_example")
		if err != nil {
			log.Fatal(err)
		}

		// Clean up automatically
		defer os.Remove(tmpfile.Name())

		// Write out all the random bytes to the file
		if _, err := tmpfile.Write(randomBytes); err != nil {
			log.Fatal(err)
		}

		// Close the file as we want osutil.FileDigest to read the file
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}

		// Can't take the address of .Name() method, so save it in a variable first
		var tempfileName string
		tempfileName = tmpfile.Name()
		fileStr = &tempfileName
	}

	// Make all the hasher objects
	hasherTable := initializeHasherTable()

	// Read in the file as a byte array for the bytes test
	bytes, err := ioutil.ReadFile(*fileStr)
	if err != nil {
		log.Fatalf("error reading file: %+v\n", err)
	}

	// For each hasher run the hash and save the results
	results := make([]HashResult, len(hasherTable))
	for index, hasher := range hasherTable {
		// Compute the hash for each type, the bytes and the file
		_, timeElapsedBytes := timeBytesHash(hasher.Hasher, bytes)
		_, timeElapsedFile := timeFileHash(hasher.HashType, *fileStr)

		results[index] = HashResult{
			Algorithm: hasher.Name,
			BytesHash: int64(timeElapsedBytes / timeVal),
			FileHash:  int64(timeElapsedFile / timeVal),
		}
	}

	// Finally, handle the output format, json or yaml
	var resultBytes []byte
	switch *formatStr {
	case "json":
		resultBytes, err = json.Marshal(results)
	case "yaml":
		resultBytes, err = yaml.Marshal(results)
	default:
		log.Fatalf("error: invalid format %s (supported formats are yaml or json)\n", *formatStr)
	}

	// Check on the encoding error
	if err != nil {
		log.Fatalf("error: failed to encode output: %+v\n", err)
	}

	// Check if we should escapte the quote character or not
	if *escapteQuotes {
		fmt.Println(strings.Replace(string(resultBytes[:]), "\"", "\\\"", -1))
	} else {
		fmt.Println(string(resultBytes[:]))
	}
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
