FROM golang:alpine as builder

WORKDIR /tmp/build

COPY go.* ./

RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o theater .

# using static nonroot image
# user:group is nobody:nobody, uid:gid = 65534:65534
FROM gcr.io/distroless/static as final

COPY --from=builder /tmp/build/theater /go/bin/theater

ENTRYPOINT ["/go/bin/theater"]