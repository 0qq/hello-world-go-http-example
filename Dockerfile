FROM golang:1.16-alpine AS build

WORKDIR /wd

# Downdload and cache deps
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code and build app
COPY pkg ./pkg
COPY cmd ./cmd
RUN go build -o /wd/hello-world ./cmd/hello-world 

FROM alpine:3.14

COPY --from=build /wd/hello-world /hello-world

ENTRYPOINT ["/hello-world"]
