FROM alpine
CMD ["/factorizer"]
RUN apk update && apk add gcc musl-dev
ADD factorizer.c /factorizer.c
RUN gcc -O2 /factorizer.c -o /factorizer &&\
  apk del gcc musl-dev &&\
  rm -rf /var/cache/apk/
