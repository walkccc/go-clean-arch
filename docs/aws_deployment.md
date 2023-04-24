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
