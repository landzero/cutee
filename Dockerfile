FROM alpine

RUN apk add --no-cache ca-certificates

EXPOSE 3000

ENV HOST=0.0.0.0
ENV PORT=3000
ENV WEB_ENV=production

HEALTHCHECK CMD /cutee-maint --health-check

COPY cutee.out          /cutee
COPY cutee-maint.out    /cutee-maint

CMD ["/cutee"]
