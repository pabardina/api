FROM scratch
ADD deployment/ca-certificates.crt /etc/ssl/certs/
ADD hirondelle /
ENTRYPOINT ["/hirondelle"]
