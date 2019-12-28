FROM golang

# install package 
RUN wget -O gagocli https://github.com/MarkEdmondson1234/gago/releases/download/v0.2.1/gagocli-v0.2.1-linux-amd64 \
    && mv gagocli /usr/local/bin/gagocli \
    && chmod 755 /usr/local/bin/gagocli

ENTRYPOINT gagocli

CMD ["reports", "-h"]
