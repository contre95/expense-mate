# Create secret for cloudflare apikey

Move the `secrets.json.example` to `secrets.json`

```shell
mv secrets.json.example secrets.json
podman secret create mysql-secret secrets.json
```

Note: secret values should be in base64
