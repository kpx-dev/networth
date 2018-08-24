.PHONY: api deploy-infra deploy-demo deploy-landing deploy-api start-api
.SILENT: api deploy-infra deploy-demo deploy-landing deploy-api start-api

deploy-infra:
	aws cloudformation update-stack --stack-name networth --capabilities CAPABILITY_IAM --template-body file://cloud/aws.infra.yml

api:
	cd api && env GOOS=linux GOARCH=amd64 go build -o ../bin/networth

deploy-api:
	make api
	aws cloudformation package --template-file api/template.yml --s3-bucket lambda.networth.app --output-template-file /tmp/networth-api.yml --s3-prefix networth-api
	aws cloudformation deploy --template-file /tmp/networth-api.yml --stack-name networth-api --capabilities CAPABILITY_IAM --region us-east-1

start-api:
	make api
	cd api && sam local start-api

# TODO: pass in args
# validate-cfn-template:
# aws cloudformation validate-template --template-body file://$1
