project(crypto-tester)
cmake_minimum_required(VERSION 3.2 FATAL_ERROR)

include("cmake/GoLangSimple.cmake")

# Add the main target executable
ADD_GO_INSTALLABLE_PROGRAM(TARGET crypto-tester
							MAIN_SOURCE cmd/crypto-tester/main.go
							IMPORT_PATH github.com/anonymouse64/crypto_tester
							SOURCE_DIRECTORIES cmd
	)

ADD_GO_INSTALLABLE_PROGRAM(TARGET sha3sum
							MAIN_SOURCE cmd/sha3sum/main.go
							IMPORT_PATH github.com/anonymouse64/crypto_tester
							SOURCE_DIRECTORIES cmd
	)

