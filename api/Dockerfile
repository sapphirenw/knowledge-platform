FROM golang:1.22.5-bookworm as go-builder

# grab the needed certs
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

# build the go app
COPY go.mod go.sum ./
RUN go mod download
COPY src ./src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/*.go

FROM alpine:latest

# copy the certs and the binary
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /app/main .

# install dependencies
RUN apk add --no-cache poppler-utils wv unrtf tidyhtml perl wget

# install pdflatex
RUN wget http://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz
RUN tar -xzf install-tl-unx.tar.gz
RUN cd install-tl-20*
RUN ./install-tl

EXPOSE 8080

# Command to run the executable
CMD ["/main"]
