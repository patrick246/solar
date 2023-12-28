FROM golang:1.21 as builder
WORKDIR /src/solar
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app -ldflags "-s -w" -trimpath ./cmd/statistics/

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /src/solar/app /app

ENTRYPOINT ["/app"]