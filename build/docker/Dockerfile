FROM alpine:3.21.3

RUN apk add --no-cache tzdata
ENV TZ=UTC
COPY ./imageroot /

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/nrtm4serve"]
CMD ["--webdir", "/www"]
