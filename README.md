# Pub DB

## tasks

### apply

directory: stack
environment: AWS_PROFILE=pubdb

```shell
terraform apply -var-file=tfvars/pro.auto.tfvars -auto-approve
```

### format

requires: format-tf

### format-tf

directory: stack

```shell
terraform fmt --recursive
```

### init

directory: stack
environment: AWS_PROFILE=pubdb

```shell
terraform init -backend-config=tfvars/pro.backend.tfvars
```

### plan

directory: stack
environment: AWS_PROFILE=pubdb

```shell
terraform plan -var-file=tfvars/pro.auto.tfvars
```

### push-config

environment: AWS_PROFILE=pubdb
directory: stack

```shell
aws s3 cp tfvars/pro.auto.tfvars s3://pub-db-pro-state-20260619093551840000000001
aws s3 cp tfvars/pro.backend.tfvars s3://pub-db-pro-state-20260619093551840000000001
```

### pull-config

environment: AWS_PROFILE=pubdb
directory: stack

```shell
mkdir -p tfvars
aws s3 cp s3://pub-db-pro-state-20260619093551840000000001/pro.auto.tfvars tfvars/
aws s3 cp s3://pub-db-pro-state-20260619093551840000000001/pro.backend.tfvars tfvars/
```
