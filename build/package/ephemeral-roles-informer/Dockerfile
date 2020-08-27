FROM alpine:latest as base

FROM scratch

COPY --from=base /etc/ssl/certs /etc/ssl/certs
COPY ephemeral-roles-informer .

ENTRYPOINT ["./ephemeral-roles-informer"]
