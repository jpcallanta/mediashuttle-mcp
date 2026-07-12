FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /mediashuttle-mcp ./cmd/mediashuttle-mcp

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /mediashuttle-mcp /usr/local/bin/mediashuttle-mcp
EXPOSE 8080
ENTRYPOINT ["mediashuttle-mcp"]
CMD ["serve"]