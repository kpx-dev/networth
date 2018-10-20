.PHONY: api deploy-infra-staging deploy-infra-prod deploy-api start-api dbstream notification deploy-notification update-lib deploy-dbstream sync deploy-sync sync
.SILENT: api deploy-infra-staging deploy-infra-prod deploy-api start-api dbstream notification deploy-notification update-lib deploy-dbstream deploy-sync sync

# staging
# ENV=staging
# DOMAIN_NAME=knncreative.com
# CLOUDFRONT_DISTRIBUTION_ID?=EN255N14EJZLA

# prod
ENV=production
DOMAIN_NAME=networth.app
CLOUDFRONT_DISTRIBUTION_ID?=E2777SQLYCBXKE

LAMBDA_BUCKET=lambda.${DOMAIN_NAME}
LANDING_S3_BUCKET?=${DOMAIN_NAME}
WEBAPP_S3_BUCKET?=webapp.${DOMAIN_NAME}
REGION=us-east-1
APP_NAME=networth
TIMESTAMP=$(shell date +%s)
GO_BUILD_OPS=-ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o ../bin/${APP_NAME}

api:
	cd api && env GOOS=linux go build ${GO_BUILD_OPS}-api .
	cd bin && zip ${APP_NAME}-api.zip ${APP_NAME}-api

dbstream:
	cd dbstream && env GOOS=linux go build ${GO_BUILD_OPS}-dbstream .
	cd bin && zip ${APP_NAME}-dbstream.zip ${APP_NAME}-dbstream

sync:
	cd sync && env GOOS=linux go build ${GO_BUILD_OPS}-sync .
	cd bin && zip ${APP_NAME}-sync.zip ${APP_NAME}-sync

notification:
	cd notification && env GOOS=linux go build ${GO_BUILD_OPS}-notification .
	cd bin && zip ${APP_NAME}-notification.zip ${APP_NAME}-notification

start-api:
	cd api && gin --appPort 8000

start-web:
	cd web && npm run start

update-lib:
	cd sync && go get && go mod tidy && go mod vendor
	cd notification && go get && go mod tidy && go mod vendor
	cd dbstream && go get && go mod tidy && go mod vendor

test:
	cd api && go test
	cd notification && go test
	cd dbstream && go test
	cd sync && go test
	# cd web && npm test

deploy-infra-staging:
	ln -s -f ~/.aws/credentials.staging.networth ~/.aws/credentials
	cd cloud && terraform workspace select staging && terraform apply -auto-approve

deploy-infra-prod:
	ln -s -f ~/.aws/credentials.prod.networth ~/.aws/credentials
	cd cloud && terraform workspace select prod && terraform apply -auto-approve

deploy-api:
	make api
	aws s3 cp bin/${APP_NAME}-api.zip s3://${LAMBDA_BUCKET}/api/${TIMESTAMP}.zip > /dev/null
	aws lambda update-function-code --function-name ${APP_NAME}-api --zip-file fileb://bin/${APP_NAME}-api.zip --publish > /dev/null

deploy-dbstream:
	make dbstream
	aws s3 cp bin/${APP_NAME}-dbstream.zip s3://${LAMBDA_BUCKET}/dbstream/${TIMESTAMP}.zip > /dev/null
	aws lambda update-function-code --function-name ${APP_NAME}-dbstream --zip-file fileb://bin/${APP_NAME}-dbstream.zip --publish > /dev/null

deploy-sync:
	make sync
	aws s3 cp bin/${APP_NAME}-sync.zip s3://${LAMBDA_BUCKET}/sync/${TIMESTAMP}.zip > /dev/null
	aws lambda update-function-code --function-name ${APP_NAME}-sync --zip-file fileb://bin/${APP_NAME}-sync.zip --publish > /dev/null

deploy-notification:
	make notification
	aws s3 cp bin/${APP_NAME}-notification.zip s3://${LAMBDA_BUCKET}/notification/${TIMESTAMP}.zip > /dev/null
	aws lambda update-function-code --function-name ${APP_NAME}-notification --zip-file fileb://bin/${APP_NAME}-notification.zip --publish > /dev/null

deploy-landing:
	aws s3 sync landing s3://${LANDING_S3_BUCKET}
	aws cloudfront create-invalidation --paths '/*' --distribution-id ${CLOUDFRONT_DISTRIBUTION_ID}

deploy-webapp:
	cd web && npx env-cmd .env.${ENV} npm run build
	aws s3 sync web/build s3://${WEBAPP_S3_BUCKET}/app --delete --acl public-read
	aws cloudfront create-invalidation --paths '/app*' --distribution-id ${CLOUDFRONT_DISTRIBUTION_ID}
