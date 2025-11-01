### Generate RSA keys

```shell
mkdir -p keys && \
cd keys && \
openssl genrsa -out access_private.pem 4096 && \
openssl rsa -in access_private.pem -pubout -out access_public.pem && \
openssl genrsa -out refresh_private.pem 4096 && \
openssl rsa -in refresh_private.pem -pubout -out refresh_public.pem
```

### Add hosts to `/etc/hosts`

- 127.0.0.1 api.smarthome.hipahopa.io
- 127.0.0.1 smarthome.hipahopa.io

### RUN
```
docker compose up -d
```

[//]: # (### Apply migrations from [db]&#40;db&#41; folder)

[//]: # ()

[//]: # (Can be done with [migrate]&#40;https://github.com/golang-migrate/migrate.git&#41;  tool)
