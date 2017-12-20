build:
	go build

test:
	go test . ./lib ./mocks ./models ./worker
