# httputil

`httputil` is a simple Go program deployed as a Docker image that is heavily inspired by the [`httpenv` from Bret Fisher](https://hub.docker.com/r/bretfisher/httpenv) with some additional functionality.

This image is for deploying and verifying container based systems and provides the following endpoints:

| Endpoint | Status Codes | Usage | Environment Variables |
|--|--|--|--|
| `/` | `200/400` | Displays all environment variables as JSON | `n/a` |
| `/redis` | `200/400` | Tries to ping a Redis database, returns PONG as JSON when successful | `REDIS_URL`, `REDIS_PORT`,`REDIS_DB` |
| `/database` | `200/400` | Tries to ping a database, returns PONG as JSON when successful.  | `DB_SERVER`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_DRIVER` |

> 💡 **Tip:** This image should **not be used in production** since it directly exposes container environment variables which can contain sensitive data.
