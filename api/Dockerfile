FROM golang:alpine as go-builder

# grab the needed certs
RUN apk --no-cache add ca-certificates

WORKDIR /app

# build the go app
COPY go.mod go.sum ./
RUN go mod download
COPY src ./src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/*.go

FROM scratch

# copy the certs and the binary
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /app/main .

EXPOSE 8080

# Command to run the executable
CMD ["/main"]
