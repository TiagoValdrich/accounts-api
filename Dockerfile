FROM golang:1.25.6 AS builder

WORKDIR /app

ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY . .

RUN make install && make build

########################################################

FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates tzdata

COPY --from=builder /app/accounts-api .

COPY --from=builder /app/db/migrations ./db/migrations

ENV DB_HOST=""
ENV DB_USER=""
ENV DB_PASSWORD=""
ENV DB_NAME=""
ENV DB_SSLMODE="disable"
ENV DB_DEBUG=false

EXPOSE 8889

CMD ["./accounts-api"]
