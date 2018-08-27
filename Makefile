.PHONY: api deploy-infra deploy-api start-api
.SILENT: api deploy-infra deploy-api start-api

api:
	cd api && env GOOS=linux GOARCH=amd64 go build -o ../bin/networth

deploy-infra:
	aws cloudformation deploy --template-file cloud/aws.infra.yml --stack-name networth --capabilities CAPABILITY_IAM --region us-east-1

deploy-api:
	make api
	aws cloudformation package --template-file api/template.yml --s3-bucket lambda.networth.app --output-template-file /tmp/networth-api.yml --s3-prefix networth-api
	aws cloudformation deploy --template-file /tmp/networth-api.yml --stack-name networth-api --capabilities CAPABILITY_IAM --region us-east-1

start-api:
	make api
	cd api && sam local start-api
