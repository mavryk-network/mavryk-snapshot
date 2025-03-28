ARG VERSION=latest
FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
# To use the libc functions for net and os/user, and still get a static binary (for containers)
# https://github.com/remotemobprogramming/mob/issues/393
RUN cd cmd/photographer && go build -ldflags "-linkmode 'external' -extldflags '-static'" -o /main

FROM mavrykdynamics/mavryk:${VERSION}
COPY --from=builder /main ./

RUN sudo apk add lz4

USER mavryk
ENV USER=mavryk

ENTRYPOINT ["./main"]
CMD [""]
