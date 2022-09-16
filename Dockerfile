FROM python:3.7.14-alpine3.16

LABEL Maintainer="catenax-ng"

ENV NR_API_KEY="valuefromdockerfile"
ENV UWSGI_PROFILE="core"

WORKDIR /usr/app/src

COPY maintenance-dashboard-app.py ./
COPY config.json ./
COPY requirements.txt ./

RUN apk add --update python3-dev build-base linux-headers pcre-dev uwsgi-python

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

CMD [ "uwsgi","--plugin","/usr/lib/uwsgi/python_plugin.so","--http-socket",":5000","--wsgi-file","maintenance-dashboard-app.py","--callable","app" ]
