FROM alpine
EXPOSE 8040
WORKDIR /blog/
ENV GOPATH=/go
RUN apk --no-cache add go nodejs git &&\
	mkdir -p /blog/frontend/ && \
	npm install -g npm &&\
	npm install -g ember-cli bower

ADD backend/ /go/src/github.com/thriqon/backend/
RUN go build -o /blog/blog github.com/thriqon/backend/

ADD frontend/package.json /frontend/
ADD frontend/npm-shrinkwrap.json /frontend/
ADD frontend/bower.json /frontend/
RUN cd /frontend/ && npm install && bower --allow-root install
ADD frontend/ /frontend/
RUN cd /frontend/ && ember build -prod --output-path /blog/frontend/
