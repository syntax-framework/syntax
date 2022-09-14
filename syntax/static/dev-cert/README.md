The existing certificate in that directory is only used for the development environment (localhost) and when the mode is
dev.

Isso para simplificar o processo de desenvolvimento

Ele foi gerado com o comando abaixo:

```
openssl req -x509 -out localhost.crt -keyout localhost.key -newkey rsa:2048 -nodes -sha256 \
-subj '/CN=localhost' -extensions EXT -config <( \
 printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
```

## References

- https://letsencrypt.org/docs/certificates-for-localhost/
- https://www.contextualcode.com/Blog/Using-self-signed-SSL-certificates-in-local-development
- https://gitlab.com/contextualcode/selfsigned-ssl-certificates/-/blob/master/generate.sh
