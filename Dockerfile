FROM golang:1.14 AS build

WORKDIR /usr/src/app
RUN mkdir -p /usr/local/bin/allbctl

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /allbctl

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /allbctl /allbctl

USER nonroot:nonroot

ENTRYPOINT ["/allbctl"]
CMD ["help"]