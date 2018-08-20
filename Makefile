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

deploy-api:
	export GOOS=linux
	cd api && go build -o ../bin/networth .
	sam package --template-file api/template.yml --s3-bucket lambda.networth.app --output-template-file /tmp/networth-api.yml --s3-prefix networth-api
	sam deploy --template-file /tmp/networth-api.yml --stack-name networth-api --capabilities CAPABILITY_IAM

start-api:
	cd api && gin --appPort=8000
