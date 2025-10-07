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

- 127.0.0.1 ca.internal.hopahome.io

### Apply migrations from [db](db) folder

Can be done with [migrate](https://github.com/golang-migrate/migrate.git)  tool
___
You can pull all repositories as I do with
