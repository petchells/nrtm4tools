FROM alpine:3.21.3

RUN apk add --no-cache tzdata
ENV TZ=UTC
COPY ./app /app

EXPOSE 8080
ENTRYPOINT ["/app/nrtm4serve_unix"]
CMD ["--port", "8080", "--webroot", "/app/dist"]
