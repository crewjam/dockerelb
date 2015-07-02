FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD dockerelb /
CMD dockerelb