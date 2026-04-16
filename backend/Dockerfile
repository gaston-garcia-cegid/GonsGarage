# API GonsGarage — imagen mínima (CGO desactivado; Postgres vía driver puro).
# Build: docker build -t gonsgarage-api ./backend
# Run:  ver docker-compose.prod.yml o docker-compose.smoke.yml

FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/api .
USER nobody:nobody
EXPOSE 8080
ENV GIN_MODE=release
ENTRYPOINT ["/app/api"]
