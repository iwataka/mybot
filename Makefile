build:
	go generate
	go build

test:
	go test . ./lib ./models ./worker ./utils
