include .test.env

tests:
	go clean $(DIR)
	go test -v $(DIR)
