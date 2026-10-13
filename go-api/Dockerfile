FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /go-api .

FROM alpine:3.20

WORKDIR /app
COPY --from=build /go-api .

EXPOSE 8080

CMD ["./go-api"]