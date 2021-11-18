FROM golang:alpine as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux  go build -ldflags="-w -s" -o platform-service .

FROM scratch

WORKDIR /app

COPY --from=builder /app/platform-service /usr/bin/

ENTRYPOINT ["platform-service"]