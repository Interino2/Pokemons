FROM golang:1.27.1-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/pokemon-app .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /out/pokemon-app ./pokemon-app
COPY templates ./templates
COPY static ./static

EXPOSE 8080
CMD ["./pokemon-app"]