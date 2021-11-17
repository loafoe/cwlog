# cwlog

A PoC for shipping Cloudwatch logs to HSDP Logging

## How it works

We use the AWS Go SDK to fetch Cloudwatch logs from a Group and Stream and then turn it
into a `CustomLogEvent`. Currently we rely on a Cloud foundry logdrainer to ship the logs
to HSPD Logging, but adding direct API calls should be trivial using [go-hsdp-api](https://github.com/philips-software/go-hsdp-api)

## Configuration

Use environment variables

| Variable | Description | Required |
|----------|-------------|----------|
| AWS_ACCESS_KEY_ID | AWS Access Key | Y |
| AWS_SECRET_ACCESS_KEY | AWS Secret Key | Y |
| CWLOG_GROUP | The Cloudwatch Group to use | Y |
| CWLOG_STREAM | The Cloudwatch Stream from the Group to fetch | Y |

## Docker

Docker builds are available: [philipslabs/cwlog](https://hub.docker.com/r/philipslabs/cwlog/tags)

### Docker example

```shell
docker run --rm \
  -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
  -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
  -e CWLOG_GROUP=/aws/sagemaker/ProcessingJobs \
  -e CWLOG_STREAM=DRS-xxx/algo-yyy \
  cwlog:latest
```

### Cloud foundry example

Use a `manifest.yml` like below:

```yaml
---
applications:
- name: cwlog
  env:
    AWS_SECRET_ACCESS_KEY: XXX
    AWS_ACCESS_KEY_ID: YYY
    CWLOG_GROUP: "/aws/sagemaker/ProcessingJobs"
    CWLOG_STREAM: "DRS-uuid/algo-1-nnnn"
  docker:
    image: philipslabs/cwlog:latest
  services:
  - logdrainer
  processes:
  - type: web
    instances: 1
    memory: 64M
    disk_quota: 1024M
    health-check-type: process
```

## License

License is MIT
