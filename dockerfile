# MULTI STAGE BUILD

# builder stage
FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz


FROM alpine:3.17 as RUNNER
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY app.env .
COPY wait-for.sh .
COPY start.sh .

EXPOSE 3000
CMD [ "/app/start.sh","/app/main" ]
