create-aws-infra:
	aws cloudformation create-stack --stack-name networth --capabilities CAPABILITY_IAM --template-body file://cloud/aws.infra.yml

update-aws-infra:
	aws cloudformation update-stack --stack-name networth --capabilities CAPABILITY_IAM --template-body file://cloud/aws.infra.yml

validate-cfn-template:
	# TODO: pass in args
	aws cloudformation validate-template --template-body file://$1

deploy-demo:
	aws s3 cp web/index.html s3://demo.networth.app/
	aws s3 cp web/assets s3://demo.networth.app/ --recursive

deploy-landing:
	aws s3 cp --recursive landing s3://networth.app/

start-api:
	cd api && gin --appPort=8000
