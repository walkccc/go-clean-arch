# AWS Deployment

## AWS ECR - Create repository to store images built by GitHub Actions.

```bash
REPOSITORY_NAME=microservice
```

```bash
aws ecr create-repository --repository-name $REPOSITORY_NAME
```

## AWS Role - Access AWS ECR by GitHub Actions

```bash
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
REPO_URI_PATTERN="repo:walkccc/go-clean-arch:*"
GITHUB_ACTIONS_ROLE_GO_CLEAN_ARCH="GitHubActionsRoleGoCleanArch"
```

```bash
aws iam create-role \
    --role-name $GITHUB_ACTIONS_ROLE_GO_CLEAN_ARCH \
    --assume-role-policy-document '{
        "Version": "2012-10-17",
        "Statement": [
          {
            "Effect": "Allow",
            "Principal": {
              "Federated": "arn:aws:iam::'"$AWS_ACCOUNT_ID"':oidc-provider/token.actions.githubusercontent.com"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
              "StringEquals": {
                "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
              },
              "StringLike": {
                "token.actions.githubusercontent.com:sub": "'$REPO_URI_PATTERN'"
              }
            }
          }
        ]
      }'
```

```bash
aws iam attach-role-policy \
    --role-name $GITHUB_ACTIONS_ROLE_GO_CLEAN_ARCH \
    --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser
```

To delete the role, you need to detach the policy first:

```bash
aws iam detach-role-policy \
    --role-name $GITHUB_ACTIONS_ROLE_GO_CLEAN_ARCH \
    --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser
```

```bash
aws iam delete-role --role-name $GITHUB_ACTIONS_ROLE_GO_CLEAN_ARCH
```

Then, store `AWS_ROLE_TO_ASSUME` in GitHub secrets.

```bash
gh secret set AWS_ROLE_TO_ASSUME
# Paste the arn:aws:iam:"$AWS_ACCOUNT_ID":role/$GITHUB_ACTIONS_ROLE_GO_CLEAN_ARCH
```

## AWS RDS - Create db instance

Create a new security group

```bash
VPC_SECURITY_GROUP_ID=$(aws ec2 create-security-group \
    --group-name AccessPostgresAnywhere \
    --description "Access Postgres anywhere" \
    --query "SecurityGroups[*].GroupId" \
    --output text)
```

```bash
aws ec2 authorize-security-group-ingress \
    --group-name AccessPostgresAnywhere \
    --protocol tcp \
    --port 5432 \
    --cidr 0.0.0.0/0
```

or read the existing security group

```bash
VPC_SECURITY_GROUP_ID=$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=AccessPostgresAnywhere" \
    --query "SecurityGroups[*].GroupId" \
    --output text)
```

```bash
aws rds create-db-instance \
    --engine postgres \
    --engine-version 15.2 \
    --db-instance-identifier microservice \
    --master-username root \
    --master-user-password password \
    --db-instance-class db.t3.micro \
    --allocated-storage 20  \
    --publicly-accessible \
    --vpc-security-group-ids $VPC_SECURITY_GROUP_ID \
    --enable-performance-insights \
    --db-name microservice \
    --backup-retention-period 7 \
    --auto-minor-version-upgrade
```

To delete the security group, you need to delete the DB instance first:

```bash
aws rds delete-db-instance --db-instance-identifier microservice --skip-final-snapshot
```

```bash
aws ec2 delete-security-group --group-name AccessPostgresAnywhere
```

## AWS Secrets Manager

```bash
AWS_SECRET_ID=MyAwsSecretId
PASSWORD=password
DB_SOURCE=$(aws rds describe-db-instances \
    --db-instance-identifier microservice --no-paginate | jq -r \
    '.DBInstances[0] | "postgresql://" + .MasterUsername + ":" + "'$PASSWORD'" + "@" + .Endpoint.Address + ":" + (.Endpoint.Port | tostring) + "/microservice"')
```

```bash
aws secretsmanager create-secret \
    --name $AWS_SECRET_ID \
    --description "Environment variables used in microservice db" \
    --secret-string '{
        "DB_DRIVER": "postgres",
        "DB_SOURCE": "'$DB_SOURCE'",
        "GRPC_SERVER_ADDRESS": "0.0.0.0:50051",
        "HTTP_SERVER_ADDRESS": "0.0.0.0:8080",
        "TOKEN_SYMMETRIC_KEY": "'$(openssl rand -hex 64 | head -c 32)'",
        "ACCESS_TOKEN_DURATION": "15m",
        "REFRESH_TOKEN_DURATION": "24h"}'
```

```bash
aws iam attach-role-policy \
    --role-name $GITHUB_ACTIONS_ROLE_GO_CLEAN_ARCH \
    --policy-arn arn:aws:iam::aws:policy/SecretsManagerReadWrite
```

Then, store `AWS_SECRET_ID` in GitHub secretss

```bash
gh secret set AWS_SECRET_ID
# Paste the $AWS_SECRET_ID
```

## (Optional) Pull Docker image from AWS ECR to run locally

```bash
aws ecr get-login-password | docker login --username AWS \
    --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
```

```bash
docker pull $IMAGE_URI
docker run -p 50051:50051 -p 8080:8080 $IMAGE_URI
```

## AWS ECS - Create cluster to run containerized task

```bash
CLUSTER_NAME=MyCluster
TASK_DEFINITION=microservice
```

```bash
aws ecs create-cluster \
    --cluster-name $CLUSTER_NAME
```

```bash
aws ecs register-task-definition \
    --cli-input-json '{
      "requiresCompatibilities": ["FARGATE"],
      "family": "'$TASK_DEFINITION'",
      "containerDefinitions": [
        {
          "name": "'$TASK_DEFINITION'",
          "image": "'$AWS_ACCOUNT_ID'.dkr.ecr.us-east-1.amazonaws.com/microservice:latest",
          "essential": true
        }
      ],
      "volumes": [],
      "networkMode": "awsvpc",
      "memory": "1 GB",
      "cpu": ".5 vCPU",
      "executionRoleArn": "arn:aws:iam::'$AWS_ACCOUNT_ID':role/ecsTaskExecutionRole"
    }'
```

Before running the task on our cluster, we need to retrieve required subnets and
security groups.

For subnets,

```bash
SUBNET_IDS=$(aws ec2 describe-subnets \
    | jq '.Subnets[]' \
    | jq '.SubnetId' | jq -srR 'split("\n") | .[:-1] | join(",")')
```

Or filter out subnets with specific `AvailabilityZone` if got error like this:

```
An error occurred (UnsupportedAvailabilityZoneException) when calling the
CreateCluster operation: Cannot create cluster 'bank' because us-east-1e, the
targeted availability zone, does not currently have sufficient capacity to
support the cluster. Retry and choose from these availability zones: us-east-1a,
us-east-1b, us-east-1c, us-east-1d, us-east-1f
```

```bash
SUBNET_IDS=$(aws ec2 describe-subnets \
    | jq '.Subnets | map(select(.AvailabilityZone != "us-east-1e")) []' \
    | jq '.SubnetId' | jq -srR 'split("\n") | .[:-1] | join(",")')
```

For security group, create two security groups to access port 8080 and port
50051 from anywhere, respectively,

```bash
HTTP_VPC_SECURITY_GROUP_ID=$(aws ec2 create-security-group \
    --group-name AccessHTTPAnywhere \
    --description "Access Postgres anywhere" \
    --query "SecurityGroups[*].GroupId" \
    --output text)
```

```bash
aws ec2 authorize-security-group-ingress \
    --group-name AccessHTTPAnywhere \
    --protocol tcp \
    --port 8080 \
    --cidr 0.0.0.0/0
```

```bash
GRPC_VPC_SECURITY_GROUP_ID=$(aws ec2 create-security-group \
    --group-name AccessGRPCAnywhere \
    --description "Access Postgres anywhere" \
    --query "SecurityGroups[*].GroupId" \
    --output text)
```

```bash
aws ec2 authorize-security-group-ingress \
    --group-name AccessGRPCAnywhere \
    --protocol tcp \
    --port 50051 \
    --cidr 0.0.0.0/0
```

or read the existing security groups

```bash
HTTP_VPC_SECURITY_GROUP_ID=$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=AccessHTTPAnywhere" \
    --query "SecurityGroups[*].GroupId" \
    --output text)
GRPC_VPC_SECURITY_GROUP_ID=$(aws ec2 describe-security-groups \
    --filters "Name=group-name,Values=AccessGRPCAnywhere" \
    --query "SecurityGroups[*].GroupId" \
    --output text)
```

Finally, we can create a service running the task:

```bash
SERVICE_NAME=MyMicroservice
```

```bash
aws ecs create-service \
    --service-name $SERVICE_NAME \
    --cluster $CLUSTER_NAME \
    --task-definition $TASK_DEFINITION \
    --launch-type FARGATE \
    --platform-version LATEST \
    --desired-count 1 \
    --network-configuration '{
        "awsvpcConfiguration": {
          "subnets": ['$SUBNET_IDS'],
          "securityGroups": ["'$HTTP_VPC_SECURITY_GROUP_ID'", "'$GRPC_VPC_SECURITY_GROUP_ID'"],
          "assignPublicIp": "ENABLED"
        }
      }'
```

Finally, we can run the task on our cluster:

```bash
aws ecs run-task \
    --cluster $CLUSTER_NAME \
    --task-definition $TASK_DEFINITION \
    --launch-type FARGATE \
    --platform-version LATEST \
    --network-configuration '{
        "awsvpcConfiguration": {
          "subnets": ['$SUBNET_IDS'],
          "securityGroups": ["'$HTTP_VPC_SECURITY_GROUP_ID'", "'$GRPC_VPC_SECURITY_GROUP_ID'"],
          "assignPublicIp": "ENABLED"
        }
      }'
```

To clean up all created running services, running tasks, and clusters,

```bash
aws ecs delete-service --cluster $CLUSTER_NAME --service $SERVICE_NAME --force
aws ecs delete-cluster --cluster $CLUSTER_NAME
```

To clean up all task definitionsr,

```bash
aws ecs deregister-task-definition --task-definition $TASK_DEFINITION:1
aws ecs delete-task-definitions --task-definitions $TASK_DEFINITION:1
# ...
```
