FROM golang

RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/go-martini/martini
RUN go get github.com/martini-contrib/render
RUN go get github.com/PuerkitoBio/goquery

ADD . /go/src/test/crawler_martini
RUN go install test/crawler_martini

EXPOSE 8888

ENTRYPOINT ["/go/bin/crawler_martini"]