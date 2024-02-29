FROM golang:1.22 as build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY cmd/ cmd/
WORKDIR /workspace/cmd/openmetadata-initializer
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /workspace/cmd/openmetadata-initializer/openmetadata-initializer /home/nonroot/openmetadata-initializer
USER 65532:65532
CMD ["/home/nonroot/openmetadata-initializer"]
