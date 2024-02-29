FROM golang:1.22 as build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY cmd/ cmd/
WORKDIR /workspace/cmd
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o /go/bin/renewer

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /go/bin/renewer /home/nonroot/
COPY confs/ /home/nonroot/
USER 65532:65532
CMD ["/home/nonroot/renewer"]
