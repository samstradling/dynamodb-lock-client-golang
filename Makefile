
test:
	go test -v

install:
	dep ensure

coverage:
	go test -cover
