FROM golang:1.24.5-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /out/app /app

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
