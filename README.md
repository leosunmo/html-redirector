# html-redirector
Simple service that replies with an html snippet that causes a redirect. Good for Organizr + Traefik Auth redirects for example.

## Usage
It requires a REDIRECT_URL environment variable, the rest works using URL query parameters (following Organizr format).

```
$ REDIRECT_URL=organizr.${DOMAIN} ./html-redirector &
Serving on HTTP port: 3000

$ curl 'localhost:3000/?error=401&return=https://myservice.${DOMAIN}'

<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Refresh" content="0; url=organizr.${DOMAIN}/?error=401&return=https://myservice.${DOMAIN}" />
  </head>
</html>
```

Working Organizr examples:

Docker-compose:
```yaml
  sonarr:
    image: myservice
    ...
    labels:
      - "traefik.enable=true"
      - "traefik.docker.network=web"
      - "traefik.http.frontend.rule=Host:myservice.${DOMAIN}"
      - "traefik.http.protocol=http"
      - "traefik.http.port=8989"
      - "traefik.frontend.auth.forward.address=https://${ORGANIZR}.${DOMAIN}/api/?v1/auth&group=1"
      - "traefik.frontend.errors.all-error-pages.backend=html-redirector"
      - "traefik.frontend.errors.all-error-pages.status=401"
      - "traefik.frontend.errors.all-error-pages.query=/?error={status}&return=https://myservice.${DOMAIN}"
```
Assuming you have the `html-redirector` as a configured backend in Traefik, this will redirect to Organizr for authentication when it received a 401 and then return to `https://myservice.${DOMAIN}` after you've signed in.


Kubernetes manifest:
`ingress.yaml`
```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: myservice
  annotations:
    kubernetes.io/ingress.class: "traefik"
    ingress.kubernetes.io/auth-type: forward
    ingress.kubernetes.io/auth-url: http://organizr.${DOMAIN}/api/?v1/auth&group=1
    traefik.ingress.kubernetes.io/error-pages: |-
      auth:
        status:
        - "401"
        backend: html-redirector
        query: "/?error=401&return=https://myservice.${DOMAIN}"
```