# Build stage
FROM golang:1.19-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go      

# Run stage
FROM alpine:3.17
RUN apk add --no-cache mysql-client
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
RUN ["chmod", "+x", "./start.sh"]
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]