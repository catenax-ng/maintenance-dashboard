FROM python:3.7.14-alpine3.16

LABEL Maintainer="catenax-ng"

ENV NR_API_KEY="valuefromdockerfile"
ENV UWSGI_PROFILE="core"

WORKDIR /usr/app/src

COPY maintenance-dashboard-app.py ./
COPY config.json ./
#COPY requirements.txt ./
COPY crontab ./
#COPY startup.sh ./

#RUN chmod u+x startup.sh

RUN apk add python3-dev build-base linux-headers pcre-dev

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir requests argparse Flask prometheus_client uwsgi pyjson

RUN pip show uwsgi

#RUN crontab crontab

#CMD ["/bin/sh","-c","/usr/app/src/startup.sh"]

#CMD ["crond","-f"]

#CMD [ "uwsgi","--http",":8000","--wsgi-file","maintenance-dashboard-app.py","--callable","app" ]
