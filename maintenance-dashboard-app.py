from flask import Flask
from werkzeug.middleware.dispatcher import DispatcherMiddleware
from prometheus_client import make_wsgi_app
import json
import requests
from os import environ
from prometheus_client.core import GaugeMetricFamily, REGISTRY, CounterMetricFamily
from apscheduler.schedulers.background import BackgroundScheduler

app = Flask(__name__)

nr_api_key = environ['NR_API_KEY']

def query_and_register_collectors():
  configjson = json.load(open('config.json'))
  apps = configjson['apps']
  results = []
  for app in apps:
    name = app['name']
    deployed = app['deployed']
    project = app['project']
    prefix = app['prefix']
    kps_releases = []
    nr_api_url = configjson['newreleases_api_url']
    projects_url = f'{nr_api_url}/projects'
    search_url = f'{projects_url}/search?q={project}'
    response_search = requests.get(search_url,
      headers = {'Content-Type':'application/json',
                 'X-Key':nr_api_key})
    response_search.raise_for_status()
    json_response_search = response_search.json()
    projectid = json_response_search['projects'][0]['id']
    releases_url = f'{projects_url}/{projectid}/releases'
    response_releases = requests.get(releases_url,
      headers = {'Content-Type':'application/json',
                 'X-Key':nr_api_key})
    response_search.raise_for_status()
    json_response_releases = response_releases.json()
    latest = json_response_releases['releases'][0]
    if 'prometheus' in app['project']:
      for rel in json_response_releases['releases']:
        if 'kube-prometheus-stack' in rel['version']:
          kps_releases.append(rel)
      latest = kps_releases[0]
    latestversion = latest.get('version').replace(prefix,'') or latest.get('version')
    results.append({'name': name, 'deployed': deployed, 'latest': latestversion})

  class VersionsCollector(object):
    def __init__(self):
      pass
    def collect(self):
      gauge = GaugeMetricFamily('appversions', 'A gauge for software versions', labels=['name','deployed','latest'])
      for r in results:
        gauge.add_metric([r['name'],r['deployed'],r['latest']], 1.0)
      yield gauge

  collectors = list(REGISTRY._collector_to_names.keys())
  for collector in collectors:
    REGISTRY.unregister(collector)
  REGISTRY.register(VersionsCollector())
  return collectors

app.wsgi_app = DispatcherMiddleware(app.wsgi_app, {
    '/metrics': make_wsgi_app()
})

if __name__ == '__main__':
    scheduler = BackgroundScheduler()
    scheduler.add_job(query_and_register_collectors, 'interval', seconds=15)
    scheduler.start()
    app.run()
