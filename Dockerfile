FROM node:current-alpine as ui-builder

WORKDIR /tmp/ui

COPY ui/package*.json ./

RUN npm ci

COPY ui/ ./

RUN npm run build

FROM golang:alpine as builder

WORKDIR /tmp/build

COPY go.* ./

RUN go mod download
RUN go mod verify

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o theater .

FROM gcr.io/distroless/static:debug as final

COPY --from=builder /tmp/build/theater /theater
COPY --from=ui-builder /tmp/ui/build /web

ENTRYPOINT ["/theater"]