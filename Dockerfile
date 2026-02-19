FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /flock-api .

FROM gcr.io/distroless/static-debian12

COPY --from=build /flock-api /flock-api
EXPOSE 8080 9090
ENTRYPOINT ["/flock-api"]
