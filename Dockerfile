FROM debian

MAINTAINER zakharov.andrii@gmail.com

ENV TOKEN=mMovLkapqLxAMseLQ4oBwZ6aq

ENV PORT=80
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379
ENV REDIS_POOL_SIZE=10
ENV REDIS_CONNECTION_TYPE=tcp

ADD gokiezen /opt/gokiezen/gokiezen
ADD start.sh /opt/gokiezen/start.sh
ADD stats.html /opt/gokiezen/stats.html

RUN ["chmod", "+x", "/opt/gokiezen/start.sh"]

RUN ["apt-get", "update"]
RUN ["apt-get", "install", "-y", "ca-certificates"]

ENTRYPOINT /opt/gokiezen/start.sh
