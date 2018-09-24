.PHONY: api deploy-infra deploy-api start-api dbstream create-infra token notification deploy-notification update-lib deploy-dbstream
.SILENT: api deploy-infra deploy-api start-api dbstream create-infra token notification deploy-notification update-lib deploy-dbstream

# staging
ENV=staging
DOMAIN_NAME=knncreative.com
CLOUDFRONT_DISTRIBUTION_ID?=EN255N14EJZLA

# prod
# ENV=production
# DOMAIN_NAME=networth.app
# CLOUDFRONT_DISTRIBUTION_ID?=E1N6WQQH3K4M1R

LAMBDA_BUCKET=lambda.${DOMAIN_NAME}
LANDING_S3_BUCKET?=${DOMAIN_NAME}
WEBAPP_S3_BUCKET?=webapp.${DOMAIN_NAME}
REGION=us-east-1
APP_NAME=networth

api:
	rm -rf bin/* && cd api && env GOOS=linux go build -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o ../bin/${APP_NAME}-api .

dbstream:
	rm -rf bin/* && cd api/dbstream && env GOOS=linux go build -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o ../../bin/${APP_NAME}-dbstream .

notification:
	rm -rf bin/* && cd api/notification && env GOOS=linux go build -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o ../../bin/${APP_NAME}-notification .

validate-template:
	aws cloudformation validate-template --template-body file://cfn/${APP_NAME}.yml --region ${REGION}

create-infra:
	aws cloudformation create-stack --template-body file://cfn/${APP_NAME}.yml --stack-name ${APP_NAME}-infra --capabilities CAPABILITY_IAM --region ${REGION}
	aws cloudformation wait stack-create-complete

start-api:
	cd api && gin --appPort 8000
	# make api
	# cd api && sam local start-api --env-vars .env.json

start-web:
	cd web && npm run start

start-dbstream:
	make dbstream
	sam local generate-event dynamodb | sam local invoke "NetWorthDbStreamFunction" --template cfn/dbstream.yml

token:
	# reset pass
	# aws cognito-idp admin-respond-to-auth-challenge --user-pool-id ${REGION}_5cJz62UiG --challenge-name NEW_PASSWORD_REQUIRED --client-id 2tam11a22g38in2vqcd5kge3cu --challenge-responses USERNAME=demo@networth.app,NEW_PASSWORD=Testing!!1234. --session "kvdceEkVdbV6hI-_8Gpk4--kgIIu89Fic5GYW-F5BiL94WvWrgVB4_ZDjOUWGmQRSpNdC7qxjjgfThHaRAxspPIYTnOUql9RIiSVBtIyzc3sMU8gtQYjyjkJs04zy4gZROp8FA6GUe41Se_J9I5s_J9zWa7g36OWIpFRivJ2KGmVVpxAZnePYA3NzG01sTOBDd2XXWP_j4wb-c27FzCBlUSba3dlBZORysRhtEg-mtoGeRxu3KMc5RzSJELw-56fqcWZOu22aFV0qbH4sixWfhdZQmSkdAce50CYQt5UQCN6ZF-IgLU2Itghb2LaYO_rcICo5WKLnhF61iQ9Vqn9d9_VA7rljVsGuCofXNAwnqXhKfSMxbJT8hHVcYuX90XUnT--CvcU9tY51F_pa_KYsulrKRWvrkizbVRcIvnvIo8tInI6YQE5yWOzHN6UAB3rb9RK_6bqAcLO6yTAgYjJ-sLBqsozN7DLijYnIJ52-Ch74mSJATbfaBHAu22mT8hsMvhxS0NvbpY_NjJl5sJAAoaW4mea4p4i0LDfI3k8hKkp_JUav_R6yz8LZ01jW7Cco10qHbHO6RYoiPD6V8xPF-CfZaOiWQw9LceoJA_RbOiFKkgTGuoSR1p6YdtTE0G_5PNqsO4oF2bt6Yo6uAnTYb5fZYYVHDfD-kC5uN3lQ_tIp6ojEI70xhm1XfYh0FwpAYuHyje9ucfp2vlNFI6sXyyJBFD7rInp6dFwTia0KgUfxNvpFJrCa9ls3VS92qh5Hbv6lo_2arf0zbemq_oS2xP2-Rflyscbi1qPRjUdCu4vkT9zYg6tnAxMGWmEzoFwZ2Ju6-4Kcv3dkyImAw-DqQ"
	aws cognito-idp initiate-auth --client-id 5a1pls13n4igqenffk3s8cnb00 --auth-flow USER_PASSWORD_AUTH --auth-parameters USERNAME=demo@networth.app,PASSWORD=Testing!!1234. | jq -r .AuthenticationResult.IdToken

update-lib:
	cd api/notification && go get github.com/networth-app/networth/api/lib && go mod download && go mod vendor && cd ../
	cd api/dbstream && go get github.com/networth-app/networth/api/lib && go mod download && go mod vendor && cd ../

test:
	cd api && go test
	cd api/notification && go test
	cd api/dbstream && go test
	cd web && npm test

deploy-infra:
	aws cloudformation deploy --template-file cfn/infra.yml --stack-name ${APP_NAME}-infra --capabilities CAPABILITY_IAM --region ${REGION} --no-fail-on-empty-changeset
	aws cloudformation wait stack-update-complete

deploy-api:
	make api
	aws cloudformation package --template-file cfn/api.yml --s3-bucket ${LAMBDA_BUCKET} --output-template-file /tmp/api.yml --s3-prefix ${APP_NAME}-api
	aws cloudformation deploy --template-file /tmp/api.yml --stack-name ${APP_NAME}-api --capabilities CAPABILITY_IAM --region ${REGION} --no-fail-on-empty-changeset

deploy-dbstream:
	make dbstream
	aws cloudformation package --template-file cfn/dbstream.yml --s3-bucket ${LAMBDA_BUCKET} --output-template-file /tmp/dbstream.yml --s3-prefix ${APP_NAME}-dbstream
	aws cloudformation deploy --template-file /tmp/dbstream.yml --stack-name ${APP_NAME}-dbstream --capabilities CAPABILITY_IAM --region ${REGION} --no-fail-on-empty-changeset

deploy-notification:
	make notification
	aws cloudformation package --template-file cfn/notification.yml --s3-bucket ${LAMBDA_BUCKET} --output-template-file /tmp/notification.yml --s3-prefix ${APP_NAME}-notification
	aws cloudformation deploy --template-file /tmp/notification.yml --stack-name ${APP_NAME}-notification --capabilities CAPABILITY_IAM --region ${REGION} --no-fail-on-empty-changeset

deploy-landing:
	aws s3 sync landing s3://${LANDING_S3_BUCKET}
	aws cloudfront create-invalidation --paths '/*' --distribution-id ${CLOUDFRONT_DISTRIBUTION_ID}

deploy-webapp:
	cd web && npx env-cmd .env.${ENV} npm run build
	aws s3 sync web/build s3://${WEBAPP_S3_BUCKET}/app --delete --acl public-read
	aws cloudfront create-invalidation --paths '/app*' --distribution-id ${CLOUDFRONT_DISTRIBUTION_ID}
