import boto3
import json
import os


def lambda_handler(event, context):
    if "AWS_REGION" in os.environ:
        ec2 = boto3.client("ec2", region_name=os.environ["AWS_REGION"])
    else:
        ec2 = boto3.client("ec2")
    response = ec2.describe_instances()
    terminate_list = []
    for res in response["Reservations"]:
        for i in res["Instances"]:
            if i["State"]["Code"] == 16 or i["State"]["Code"] == 80:  # Running or stopped
                if "Tags" in i:
                    preserve = False
                    for tag in i["Tags"]:
                        if tag["Key"].lower() == "preserve":
                            preserve = True
                    if not preserve:
                        terminate_list.append(i["InstanceId"])
                else:
                    terminate_list.append(i["InstanceId"])
    print(json.dumps(terminate_list))
    if len(terminate_list) > 0:
        response = ec2.terminate_instances(InstanceIds=terminate_list)
        return json.dumps(response)
    else:
        return None


if __name__ == "__main__":
    # Not running in a Lambda environment
    lambda_handler(None, None)
