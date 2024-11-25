# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

COPY ./go.mod ./go.sum ./

COPY ./ ./

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

# Run the outyet command by default when the container starts.

ENTRYPOINT go run ./main_app/cmd/app/main.go

# Document that the service listens on port 8080.
EXPOSE 8082 8080