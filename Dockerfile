FROM alpine:3.14

WORKDIR /app

COPY discovery /app/
COPY etc/discovery.yaml /app/etc/
COPY etc/discovery-api.yaml /app/etc/

RUN chmod +x /app/discovery

EXPOSE 9999

CMD ["/app/discovery", "-f", "/app/etc/discovery-api.yaml"]
