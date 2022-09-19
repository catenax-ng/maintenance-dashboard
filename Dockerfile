FROM alpine:3.16.2

LABEL Maintainer="catenax-ng"

ENV NR_API_KEY="valuefromdockerfile"
ENV UWSGI_PROFILE="core"
ENV PYTHONUNBUFFERED=1

WORKDIR /usr/app/src

COPY maintenance-dashboard-app.py ./
COPY config.json ./
COPY requirements.txt ./

RUN apk add --update --no-cache python3 uwsgi uwsgi-python3 tzdata && ln -sf python3 /usr/bin/python
RUN ln -fs /usr/share/zoneinfo/Etc/CET /etc/localtime
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools
RUN pip3 install --no-cache --upgrade Flask prometheus-client pyjson requests APScheduler

CMD [ "uwsgi","--enable-threads","--http-socket",":5000","--plugin","/usr/lib/uwsgi/python_plugin.so","--plugins-list","--wsgi-file","maintenance-dashboard-app.py","--callable","app" ]
