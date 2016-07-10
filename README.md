# My bot

## Usage

### Local Environment

```
go get -d ./...
go build
./mybot
```

### Using [Docker](https://www.docker.com/)

```
cd path/to/mybot
docker build -t mybot .
docker run -v $(pwd):/mybot -it sh
cd /mybot
```

then run commands in the above
