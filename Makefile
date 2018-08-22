.PHONY: api deploy-infra deploy-demo deploy-landing deploy-api start-api
.SILENT: api deploy-infra deploy-demo deploy-landing deploy-api start-api

deploy-infra:
	aws cloudformation update-stack --stack-name networth --capabilities CAPABILITY_IAM --template-body file://cloud/aws.infra.yml

deploy-demo:
	aws s3 cp web/index.html s3://demo.networth.app/
	aws s3 cp web/assets s3://demo.networth.app/ --recursive
	# aws cloudfront create-invalidation --distribution-id E21OGDJ6NKWTTA --paths '/*'

deploy-landing:
	aws s3 cp --recursive landing s3://networth.app/
	aws cloudfront create-invalidation --distribution-id E21OGDJ6NKWTTA --paths '/*'

api:
	cd api && env GOOS=linux GOARCH=amd64 go build -o ../bin/networth .

deploy-api:
	cd api && env GOOS=linux GOARCH=amd64 go build -o ../bin/networth .
	sam package --template-file api/template.yml --s3-bucket lambda.networth.app --output-template-file /tmp/networth-api.yml --s3-prefix networth-api
	sam deploy --template-file /tmp/networth-api.yml --stack-name networth-api --capabilities CAPABILITY_IAM

start-api:
	make api
	 cd api && sam local start-api

# TODO: pass in args
# validate-cfn-template:
# aws cloudformation validate-template --template-body file://$1
