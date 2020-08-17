FROM golang@sha256:1e9c36b3fd7d7f9ab95835fb1ed898293ec0917e44c7e7d2766b4a2d9aa43da6 AS builder
# Ensure ca-certficates are up to date
RUN update-ca-certificates

WORKDIR /usr/local/app

# use modules
COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/app .

FROM gcr.io/distroless/static@sha256:c6d5981545ce1406d33e61434c61e9452dad93ecd8397c41e89036ef977a88f4
EXPOSE 9094
COPY --from=builder /go/bin/app /go/bin/app

ENTRYPOINT ["/go/bin/app"]