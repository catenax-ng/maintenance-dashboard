# maintenance-dashboard
DevSecOps team maintenance dashboard

notes:

wget https://go.dev/dl/go1.18.4.linux-amd64.tar.gz

rm -rf /usr/local/go

tar -C /usr/local -xzf go1.18.4.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin

go env -w GO111MODULE="on" #or "auto"

go mod init maintenance-dashboard

go mod tidy

go get github.com/prometheus/client_golang/prometheus

go get github.com/prometheus/client_golang/prometheus/promauto

go get github.com/prometheus/client_golang/prometheus/promhttp

go get k8s.io/client-go@latest

go build maintenance-dashboard.go

##############################################################

# newreleases.io api calls

NR_PROJECTS=$(curl -H "X-Key: $NR_API_KEY" https://api.newreleases.io/v1/projects | jq -r '.projects[].id')

for id in $NR_PROJECTS

  do

    curl -H "X-Key: $NR_API_KEY" https://api.newreleases.io/v1/projects/$id/releases | jq -r '.releases[0].version'

  done
