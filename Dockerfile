FROM golang:1.21.5-alpine3.18 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main . && \
    go build -o migrate ./migrations/migrate.go && \
    if [ ! -f env/config ]; then cp env/sample.config env/config ; fi

FROM busybox:1.36.1 AS runner

RUN adduser -D -u 1001 appuser

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/migrate .
COPY --from=build /app/env ./env
COPY --from=build /app/migrations/sql ./migrations/sql

USER appuser

ENTRYPOINT ["sh", "-c", "./migrate up && ./main"]
