FROM alpine:3.12
ARG CED_VERSION
ARG TARGETARCH
ARG TARGETVARIANT
RUN apk update && apk add --no-cache ca-certificates
ADD https://github.com/blmhemu/ced/releases/download/v$CED_VERSION/ced\_$CED_VERSION\_linux\_${TARGETARCH}${TARGETVARIANT} /usr/bin/ced
ADD ced.properties /etc/ced/ced.properties
ENTRYPOINT ["/usr/bin/ced"]
CMD ["-cfg", "/etc/ced/ced.properties"]
