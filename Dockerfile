FROM alpine:3.12
RUN apk update && apk add --no-cache ca-certificates
COPY ced /usr/bin
ADD ced.properties /etc/ced/ced.properties
# EXPOSE 9998 9999
ENTRYPOINT ["/usr/bin/ced"]
CMD ["-cfg", "/etc/ced/ced.properties"]