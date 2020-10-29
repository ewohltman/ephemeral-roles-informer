FROM alpine:latest as base

FROM scratch

COPY --from=base /etc/ssl/certs /etc/ssl/certs
COPY passwd.nobody /etc/passwd
COPY ephemeral-roles-informer .

USER nobody

ENTRYPOINT ["./ephemeral-roles-informer"]
