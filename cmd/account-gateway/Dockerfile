FROM uhub.service.ucloud.cn/entropypool/alpine:3.16

WORKDIR /app

COPY account-gateway /app

ENTRYPOINT ["/app/account-gateway"]

EXPOSE 50480 50481 50482 50483

CMD ["run"]
