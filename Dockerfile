FROM golang:alpine AS server-build
WORKDIR /build
ENV CGO_ENABLED=0
ENV GO_VERSION=1.18.1
ENV GO_OS=linux
ENV GO_ARCH=amd64
ENV GO_BUILD_TAGS=netgo
COPY  go.mod .
COPY  go.sum .
RUN go mod download
COPY . .
RUN go build -o /goftp main.go

FROM scratch
COPY --from=server-build goftp /goftp
ENV HOST=0.0.0.0
ENV PORT=3090
ENTRYPOINT ["/goftp"]