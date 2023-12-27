FROM golang:1.21.5-alpine3.18 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main . && \
    go build -o migrate ./migrations/migrate.go && \
    if [ ! -f env/config ]; then cp env/sample.config env/config ; fi

FROM alpine:3.18

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/migrate .
COPY --from=build /app/env ./env
COPY --from=build /app/migrations/sql ./migrations/sql

EXPOSE 8080

ENTRYPOINT ["sh", "-c", "./migrate up && ./main"]
