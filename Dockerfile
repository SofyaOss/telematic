FROM golang:1.21.3

WORKDIR /app

# Copy the source code
COPY . /app

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o api ./cmd
#RUN go test -o api ./internal/generator

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api"]