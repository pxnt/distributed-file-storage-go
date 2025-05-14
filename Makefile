.PHONY: run

run:
	go run .

# run server s1
s1:
	go run . s1

# run server s2
s2:
	go run . s2

.PHONY: test

test:
	go test -v .
