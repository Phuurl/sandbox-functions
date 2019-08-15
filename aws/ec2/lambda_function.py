import json
import boto3


def lambda_handler(event, context):
  ec2 = boto3.client("ec2", region_name="eu-west-1")
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

