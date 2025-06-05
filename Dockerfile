FROM golang:1.24.3-bullseye AS build

WORKDIR /build

COPY ./cmd /build/cmd
COPY ./internal /build/internal
COPY ./go.mod /build/go.mod
COPY ./go.sum /build/go.sum

RUN go mod vendor
RUN go build -o go-data-collector ./cmd/main.go

FROM golang:1.24.3-bullseye AS final

WORKDIR /app

COPY --from=build /build/go-data-collector .

CMD ["/app/go-data-collector"]