FROM golang:1.20-alpine AS build

# allow build-time override of GOPROXY for environments with restricted network
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

WORKDIR /app
RUN apk add --no-cache git ca-certificates

# cache modules: copy go.mod and optionally go.sum first
COPY go.mod go.sum* ./
# ensure Go uses desired GOPROXY and pull dependencies early
RUN go env -w GOPROXY=${GOPROXY} && go mod download

# copy rest of source and build
COPY . .
RUN go build -o /bin/server ./cmd/server

FROM alpine:latest
COPY --from=build /bin/server /bin/server
EXPOSE 8080
ENTRYPOINT ["/bin/server"]
