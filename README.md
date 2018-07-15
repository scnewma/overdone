# Overdone
A simple todo web server written in Go. Currently, the tasks created are only stored in memory so restarting the server will cause all of the tasks to be lost.

## Development

### Install Dependencies
This project uses [Glide](https://glide.readthedocs.io/en/latest/). After cloning the project, you will need to install the dependencies:

```
glide install
```

### Building

```
go build -o bin/overdone cmd/overdone/main.go
```

### Testing

```
go test -v ./...
```

### Running

```
bin/overdone
```
