# aws-elasticsearch-proxy
A proxy for AWS's ElasticSearch

## Install
`go get github.com/oremj/aws-elasticsearch-proxy`

## Usage
Make sure your credentials are set using the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` or
ensure that the instance running the proxy has an instance profile. Ensure that the user or account
that the credentials are assigned to have access to your ElasticSearch endpoint.

Run:
```
aws-elasticsearch-proxy \
  -region us-east-1 \
  -es-endpoint yourendpoint.us-east-1.es.amazonaws.com \
  -addr 127.0.0.1:8080
```

Test it out:
```
curl http://127.0.0.1:8080
```

