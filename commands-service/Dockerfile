FROM golang:alpine as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o commands-service .

FROM scratch

WORKDIR /app

COPY --from=builder /app/commands-service /usr/bin/

ENTRYPOINT ["commands-service"]