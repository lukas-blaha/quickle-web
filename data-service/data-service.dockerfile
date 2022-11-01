FROM alpine:latest

RUN mkdir /app

COPY dataApp /app

CMD [ "/app/dataApp" ]
