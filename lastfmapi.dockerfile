# Since, we are using Makefile, we don't need the above dockerfile configuration anymore
# Because, we are building linux executable of our go application while using
# "make build_lastfm" command. So, we only need to copy this
# executable file(brokerApp) to the /app workdir of our container(alpine image). No need to build
# multi stage dockerfile and this will take much less time than the multi stage build
FROM alpine:latest

RUN mkdir /app

COPY lastFmApp /app

EXPOSE 80

CMD [ "/app/lastFmApp" ]