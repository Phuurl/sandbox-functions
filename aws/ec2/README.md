# AWS EC2

Shuts down any running EC2s that aren't tagged with `preserve`.

## Setup

- Runtime: `Python 3.7`
- Handler: `lambda_function.lambda_handler`
- Memory: `128 MB`

## Trigger

This could be executed on a cron schedule via CloudWatch Events, such as `cron(30 0 * * ? *)`, or another suitable trigger method.

## IAM Policy

In addition to the basic Lambda execution permissions, the following policy is also required to allow listing and terminating EC2 instances:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": "ec2:TerminateInstances",
      "Resource": "arn:aws:ec2:*:*:instance/*"
    },
    {
      "Sid": "VisualEditor1",
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeTags",
        "ec2:DescribeInstanceStatus"
      ],
      "Resource": "*"
    }
  ]
}
```
