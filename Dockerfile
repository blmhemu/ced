FROM alpine:3.12
ARG CED_VERSION
RUN apk update && apk add --no-cache ca-certificates
ADD https://github.com/blmhemu/ced/releases/download/v0.1.0/ced_0.1.0_linux${TARGETPLATFORM} /usr/bin/ced
ADD ced.properties /etc/ced/ced.properties
# EXPOSE 9998 9999
ENTRYPOINT ["/usr/bin/ced"]
CMD ["-cfg", "/etc/ced/ced.properties"]