FROM golang:latest AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main . && \
    if [ ! -f env/config ]; then cp env/sample.config env/config ; fi

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/env ./env

EXPOSE 8080

CMD ["./main"]
