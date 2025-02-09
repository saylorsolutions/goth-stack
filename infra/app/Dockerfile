FROM golang:1.23 AS builder

WORKDIR /app
COPY ../.. .

RUN go run ./modmake build

FROM gcr.io/distroless/static-debian12:latest

# Used to set the connection URL to interact with the database.
ENV DBURL=""
# Used to set the hash key for session values.
ENV SESSION_HASHKEY=""

COPY --from=builder /app/build/yourapp/yourapp /yourapp

ENTRYPOINT ["/yourapp"]
