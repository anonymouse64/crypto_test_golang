package main

import (
	"crypto"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	osutil "github.com/snapcore/snapd/osutil"
	_ "golang.org/x/crypto/sha3"
)

func timeFileHash(hasher crypto.Hash, file string) ([]byte, time.Duration) {
	// Get start time
	start := time.Now()

	// Compute the hash of the file
	hashRes, _, _ := osutil.FileDigest(file, hasher)

	// Return the hash result and the time since the start of the function
	return hashRes, time.Since(start)
}

func main() {
	// Setup flags
	fileStr := flag.String("file", "", "file to hash")
	randomSizeMB := flag.Int64("size", 10, "size of generated random file")
	unitStr := flag.String("unit", "s", "units to use (possible values : ns, us, ms, s)")
	algStr := flag.String("alg", "sha3_512", "algorithm to use, sha3_512, sha3_384, sha3_256, or sha3_224")

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
		tmpfile, err := ioutil.TempFile("", "sha3sum_example")
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

	// Run a single run and print out human readable form
	// First check what algorithm to use
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
		log.Fatalf("error: algorithm %s not supported\n", *algStr)
	}

	// Run the hash
	hashBytes, timeElapsedFile := timeFileHash(hasherToUse, *fileStr)

	// Print the hash and the file name
	fmt.Printf("%x %s\n", hashBytes, *fileStr)

	// Print the stats
	fmt.Printf("Calculated in %3f sec, %5.2f MBps\n", float64(timeElapsedFile)/float64(timeVal), float64(fileSize)/1048576/(float64(timeElapsedFile)/float64(time.Second)))
}
