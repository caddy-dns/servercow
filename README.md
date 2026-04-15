# caddy-dns/servercow

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy).
It can be used to manage DNS records with [Servercow](https://servercow.de).

## Caddy module name

```
dns.providers.servercow
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "servercow",
        "username": "YOUR_SERVERCOW_USERNAME",
        "password": "YOUR_SERVERCOW_PASSWORD"
      }
    }
  }
}
```

Or with the Caddyfile:

```
# globally
{
    acme_dns servercow <username> <password>
}
```

```
# one site, inline credentials
tls {
    dns servercow <username> <password>
}
```

```
# one site, block form with env-var placeholders
tls {
    dns servercow {
        username {env.SERVERCOW_USERNAME}
        password {env.SERVERCOW_PASSWORD}
    }
}
```

## Building with xcaddy

```bash
xcaddy build --with github.com/caddy-dns/servercow
```
