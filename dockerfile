FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base && go build -o forum cmd/web/main.go

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app .

EXPOSE 4000
LABEL name="FORUM" \
      authors="dyelesho, aishagul, aladik" \
      release_date="05.06.2023"
CMD ["./forum"]
