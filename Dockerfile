FROM golang:1.22.2-alpine3.19 as go-dependencies-cache
WORKDIR /src

COPY go.mod go.sum /src/

RUN go mod download

FROM go-dependencies-cache as builder

WORKDIR /src

COPY . .
RUN go build -o /fv-server cmd/server/main.go

FROM alpine:3.19.1

COPY --from=builder /fv-server /bin

CMD ["sh", "-c", "fv-server"]