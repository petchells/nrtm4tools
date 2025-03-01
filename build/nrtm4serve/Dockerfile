FROM alpine:3.21.3

RUN apk add --no-cache tzdata
ENV TZ=UTC
COPY ./app /app

EXPOSE 8080
ENTRYPOINT ["/app/nrtm4serve_unix"]
CMD ["--webdir", "/app/dist"]
