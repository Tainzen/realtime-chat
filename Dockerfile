FROM alpine:3

RUN adduser --disabled-password --home /realtime-chat --shell /bin/sh realtime-chat

WORKDIR /realtime-chat

RUN mkdir bin conf

COPY bin bin/.

COPY conf conf/.

ENTRYPOINT ["sh","conf/export.sh"]