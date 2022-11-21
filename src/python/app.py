from flask import Flask
from werkzeug.middleware.dispatcher import DispatcherMiddleware
from prometheus_client import make_wsgi_app
import json, yaml, requests, re
from os import environ, listdir
from prometheus_client.core import GaugeMetricFamily, REGISTRY, CounterMetricFamily
from apscheduler.schedulers.background import BackgroundScheduler
from github import Github

newreleases_api_key = environ["NEWRELEASES_API_KEY"]
github_token = environ["GITHUB_TOKEN"]

with open("config/github_repo.yaml") as stream:
    try:
        github_config = yaml.safe_load(stream)
    except yaml.YAMLError as exc:
        print(exc)

with open("config/newreleases.yaml") as stream:
    try:
        newreleases_config = yaml.safe_load(stream)
    except yaml.YAMLError as exc:
        print(exc)

github_repo = github_config["github_repo"]
github = Github(github_token)
repo = github.get_repo(github_repo)

def read_github_file(path):
    try:
        fc = yaml.safe_load(repo.get_contents(path).decoded_content.decode())
    except yaml.YAMLError as exc:
        print(exc)
    return fc

def read_config_file(filename):
    with open(filename) as stream:
        try:
            deployment_config = yaml.safe_load(stream)
        except yaml.YAMLError as exc:
            print(exc)
        return deployment_config

def query_github(name, path):
    filecontent = read_github_file(path)
    if 'kind' in filecontent:
        if filecontent['kind'] == 'Kustomization':
            for item in filecontent['resources']:
                if 'argoproj/argo-cd' in item:
                    deployed = item.split('/')[5].replace('v','')
        if filecontent['kind'] == 'Deployment':
            for initcontainer in filecontent['spec']['template']['spec']['initContainers']:
                for envvar in initcontainer['env']:
                    if envvar['name'] == 'AVP_VERSION':
                        deployed = envvar['value']
        if filecontent['kind'] == 'ApplicationSet':
            if 'reflector' in filecontent['metadata']['name']:
                for overlay in filecontent['spec']['generators'][0]['list']['elements']:
                    if overlay['cluster'] == 'core':
                        deployed = overlay['targetRevision']
    elif 'Chart.yaml' in path:
        try:
            dependencies = (filecontent['dependencies'])
        except KeyError as kerr:
            dependencies = 'dependencies'
        if dependencies != 'dependencies':
            for dependency in dependencies:
                if dependency['name'] == name:
                    deployed = dependency['version']
        else:
            deployed = filecontent['version']
    return deployed

def select_latest_semver(vl,p):
    versionlist = []
    pattern = re.compile("^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$")
    for release in vl:
        semver = release['version'].replace(p,'')
        if pattern.match(semver):
            versionlist.append(semver)
    versionlist.sort(key=lambda s: [int(u) for u in s.split('.')],reverse=1)
    l = versionlist[0]
    return l

def query_newreleases(project, prefix):
    kps_releases = []
    newreleases_api_url = newreleases_config['newreleases_api_url']
    projects_url = f'{newreleases_api_url}/projects'
    search_url = f'{projects_url}/search?q={project}'
    response_search = requests.get(search_url,
      headers = {'Content-Type':'application/json',
                 'X-Key':newreleases_api_key})
    response_search.raise_for_status()
    json_response_search = response_search.json()
    projectid = json_response_search['projects'][0]['id']
    releases_url = f'{projects_url}/{projectid}/releases'
    response_releases = requests.get(releases_url,
      headers = {'Content-Type':'application/json',
                 'X-Key':newreleases_api_key})
    response_search.raise_for_status()
    json_response_releases = response_releases.json()
    latest = select_latest_semver(json_response_releases['releases'],prefix)
    if 'prometheus' in project:
      for rel in json_response_releases['releases']:
        if 'kube-prometheus-stack' in rel['version']:
          kps_releases.append(rel)
      latest = select_latest_semver(kps_releases,prefix)
    return latest

def run_queries():
    results = []
    for configfile in listdir("config/deployment"):
        filepath = f'config/deployment/{configfile}'
        config = read_config_file(filepath)
        name = config["name"]
        deployedversion = query_github(name,config["path"])
        latestversion = query_newreleases(config["project"],config["prefix"])
        results.append({'name': name, 'deployed': deployedversion, 'latest': latestversion})
    return results

def register_collectors():
    results = run_queries()
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

register_collectors()

app = Flask(__name__)

app.wsgi_app = DispatcherMiddleware(app.wsgi_app, {
    '/metrics': make_wsgi_app()
})

scheduler = BackgroundScheduler()
scheduler.add_job(register_collectors, 'interval', hours=1)
scheduler.start()
