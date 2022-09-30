FROM scratch
# copy the ca-certificate.crt from the build stage
COPY ca-certificates.crt /etc/ssl/certs/
COPY main /
CMD ["/main"]