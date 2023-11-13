ARG VERSION_GO=1.21
ARG VERSION_ALPINE=3.18

FROM golang:${VERSION_GO}-alpine${VERSION_ALPINE} AS build

WORKDIR /usr/src

COPY go.* .
RUN go mod download

# required for go-sqlite3
RUN apk add --no-cache gcc musl-dev

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build \
    -ldflags "-s -w -extldflags '-static'" \
    -buildvcs=false \
    -tags osusergo,netgo \
    -o /usr/local/bin/ ./...

FROM alpine:${VERSION_ALPINE} AS runtime

WORKDIR /opt

ARG LITEFS_CONFIG="litefs.yml"
ENV LITEFS_DIR="/litefs"
ENV DATA_SOURCE_IAM="${LITEFS_DIR}/iam"
ENV DATA_SOURCE_ROOMS="${LITEFS_DIR}/rooms"
ENV GITHUB_CLIENT_ID="Iv1.7603b1f7d60b536f"
ENV INTERNAL_PORT=8080
ENV PORT=8081

COPY --from=build /usr/local/bin/yap ./a
COPY --from=build /usr/local/bin/sqlite3 ./b

# prepare litefs
RUN apk add --no-cache fuse3 sqlite ca-certificates
COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs
ADD litefs/etc/${LITEFS_CONFIG} /etc/litefs.yml

CMD ["litefs", "mount"]