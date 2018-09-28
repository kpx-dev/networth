# data "aws_caller_identity" "current" {}

variable "region" {
  default = "us-east-1"
}

provider "aws" {
  region = "${var.region}"
}

variable "AppName" {
  default = "networth"
}

variable "DomainName" {
  default = "knncreative.com"
}

resource "aws_s3_bucket" "landing-bucket" {
  bucket = "${var.DomainName}"
}

resource "aws_s3_bucket" "lambda" {
  bucket = "lambda.${var.DomainName}"
}

resource "aws_ssm_parameter" "SLACK_CHANNEL" {
  name  = "/${var.AppName}/SLACK_CHANNEL"
  type  = "String"
  value = "sns"
}

resource "aws_ssm_parameter" "PLAID_ENV" {
  name  = "/${var.AppName}/PLAID_ENV"
  type  = "String"
  value = "sandbox"
}

resource "aws_ssm_parameter" "PLAID_CLIENT_ID" {
  name  = "/${var.AppName}/PLAID_CLIENT_ID"
  type  = "String"
  value = " "
}

resource "aws_ssm_parameter" "PLAID_SECRET" {
  name  = "/${var.AppName}/PLAID_SECRET"
  type  = "String"
  value = " "
}

resource "aws_ssm_parameter" "PLAID_PUBLIC_KEY" {
  name  = "/${var.AppName}/PLAID_PUBLIC_KEY"
  type  = "String"
  value = " "
}

resource "aws_ssm_parameter" "SLACK_WEBHOOK_URL" {
  name  = "/${var.AppName}/SLACK_WEBHOOK_URL"
  type  = "String"
  value = " "
}

resource "aws_sns_topic" "SNSTopic" {
  name = "${var.AppName}"
}

resource "aws_dynamodb_table" "db_table" {
  name             = "${var.AppName}"
  read_capacity    = 1
  write_capacity   = 1
  hash_key         = "id"
  range_key        = "sort"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  server_side_encryption {
    enabled = true
  }

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "sort"
    type = "S"
  }
}

resource "aws_cognito_user_pool" "UserPool" {
  name = "${var.AppName}"
  username_attributes = ["email"]
  auto_verified_attributes = ["email"]
}

resource "aws_cognito_user_pool_client" "UserPoolClient" {
  name         = "${var.AppName}-webapp"
  user_pool_id = "${aws_cognito_user_pool.UserPool.id}"
}

resource "aws_cognito_identity_pool" "CognitoIdentityPool" {
  identity_pool_name      = "${var.AppName}"
  developer_provider_name = "${var.AppName}"

  cognito_identity_providers = [
    {
      client_id     = "${aws_cognito_user_pool_client.UserPoolClient.id}"
      provider_name = "cognito-idp.${var.region}.amazonaws.com/${aws_cognito_user_pool.UserPool.id}"
    },
  ]

  allow_unauthenticated_identities = true
}

data "aws_iam_policy_document" "CognitoIdentityPoolPolicyDoc" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["cognito-idp.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "CognitoIdentityPoolRole" {
  name               = "CognitoIdentityPoolRole"
  assume_role_policy = "${data.aws_iam_policy_document.CognitoIdentityPoolPolicyDoc.json}"
}

resource "aws_cognito_identity_pool_roles_attachment" "CognitoIdentityPoolRoleAttachment" {
  identity_pool_id = "${aws_cognito_identity_pool.CognitoIdentityPool.id}"

  roles {
    authenticated   = "${aws_iam_role.CognitoIdentityPoolRole.arn}"
    unauthenticated = "${aws_iam_role.CognitoIdentityPoolRole.arn}"
  }
}

resource "aws_acm_certificate" "cert" {
  domain_name = "${var.DomainName}"

  subject_alternative_names = [
    "*.${var.DomainName}",
  ]

  validation_method = "EMAIL"
}

data "aws_iam_policy_document" "KMSKeyPolicyDoc" {
  statement {
    resources = ["*"]

    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
    ]

    principals {
      type        = "AWS"
      identifiers = ["aws_iam_role.LambdaRole.arn"]
    }
  }
}

resource "aws_kms_key" "KMSKey" {
  # TODO: enable
  # policy = "${data.aws_iam_policy_document.KMSKeyPolicyDoc.json}"
}

resource "aws_kms_alias" "KMSAlias" {
  name          = "alias/${var.AppName}"
  target_key_id = "${aws_kms_key.KMSKey.key_id}"
}

data "aws_iam_policy_document" "LambdaAssumeRolePolicyDoc" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "DynamoDBPolicyDoc" {
  statement {
    actions = [
      "dynamodb:*",
    ]

    resources = [
      "${aws_dynamodb_table.db_table.arn}",
    ]
  }
}

resource "aws_iam_role" "LambdaRole" {
  name               = "LambdaRole"
  assume_role_policy = "${data.aws_iam_policy_document.LambdaAssumeRolePolicyDoc.json}"
}

resource "aws_iam_policy" "DynamoDBPolicy" {
  name   = "DyanmoDBFullAccessForLambda"
  policy = "${data.aws_iam_policy_document.DynamoDBPolicyDoc.json}"
}

resource "aws_iam_role_policy_attachment" "LambdaRoleAttachDynamoDB" {
  role       = "${aws_iam_role.LambdaRole.name}"
  policy_arn = "${aws_iam_policy.DynamoDBPolicy.arn}"
}

resource "aws_iam_role_policy_attachment" "LambdaRoleAttachLambdaBasic" {
  role       = "${aws_iam_role.LambdaRole.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_cloudfront_origin_access_identity" "id" {}

resource "aws_s3_bucket" "S3BucketLambda" {
  bucket = "lambda.${var.DomainName}"
}

resource "aws_s3_bucket" "LandingS3Bucket" {
  bucket = "${var.DomainName}"
}

data "aws_iam_policy_document" "LandingS3BucketPublicReadPolicyDoc" {
  statement {
    actions = ["s3:GetObject"]

    principals {
      type        = "AWS"
      identifiers = ["${aws_cloudfront_origin_access_identity.id.iam_arn}"]
    }

    resources = ["${aws_s3_bucket.LandingS3Bucket.arn}/*"]
  }
}

resource "aws_s3_bucket_policy" "LandingS3BucketPublicReadPolicy" {
  bucket = "${aws_s3_bucket.LandingS3Bucket.id}"
  policy = "${data.aws_iam_policy_document.LandingS3BucketPublicReadPolicyDoc.json}"
}

resource "aws_s3_bucket" "WebAppS3BucketResource" {
  bucket = "webapp.${var.DomainName}"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "app/index.html"

    # TODO: use advanced routing so we might not the above? ^
    # routing_rules
  }
}

data "aws_iam_policy_document" "WebAppS3BucketPolicyPolicyDoc" {
  statement {
    actions = ["s3:GetObject"]

    principals {
      type        = "*"
      identifiers = ["*"]
    }

    resources = ["${aws_s3_bucket.WebAppS3BucketResource.arn}/*"]
  }
}

resource "aws_s3_bucket_policy" "WebAppS3BucketPolicy" {
  bucket = "${aws_s3_bucket.WebAppS3BucketResource.id}"
  policy = "${data.aws_iam_policy_document.WebAppS3BucketPolicyPolicyDoc.json}"
}

resource "aws_s3_bucket" "LoggingS3Bucket" {
  bucket = "log.${var.DomainName}"
}

resource "aws_api_gateway_rest_api" "api" {
  name = "${var.AppName}"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  parent_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_resource.proxy.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_method" "proxy_root" {
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  resource_id = "${aws_api_gateway_method.proxy.resource_id}"
  http_method = "${aws_api_gateway_method.proxy.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.LambdaAPIFunction.invoke_arn}"
}

resource "aws_api_gateway_integration" "lambda_root" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  resource_id = "${aws_api_gateway_method.proxy_root.resource_id}"
  http_method = "${aws_api_gateway_method.proxy_root.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.LambdaAPIFunction.invoke_arn}"
}

resource "aws_api_gateway_deployment" "api" {
  depends_on = [
    "aws_api_gateway_integration.lambda",
    "aws_api_gateway_integration.lambda_root",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  stage_name  = "latest"
}

resource "aws_api_gateway_authorizer" "auth" {
  name          = "cognito"
  type          = "COGNITO_USER_POOLS"
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  provider_arns = ["${aws_cognito_user_pool.UserPool.arn}"]
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.LambdaAPIFunction.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_deployment.api.execution_arn}/*/*"
}

resource "aws_cloudfront_distribution" "CloudFrontResource" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = "${var.DomainName}"
  default_root_object = "index.html"
  aliases             = ["${var.DomainName}", "www.${var.DomainName}"]
  price_class         = "PriceClass_100"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  logging_config {
    include_cookies = false
    bucket          = "${aws_s3_bucket.LoggingS3Bucket.bucket_domain_name}"
    prefix          = "${var.DomainName}"
  }

  origin {
    domain_name = "${aws_s3_bucket.LandingS3Bucket.bucket_regional_domain_name}"
    origin_id   = "landing"

    s3_origin_config {
      origin_access_identity = "${aws_cloudfront_origin_access_identity.id.cloudfront_access_identity_path}"
    }
  }

  origin {
    domain_name = "${aws_s3_bucket.WebAppS3BucketResource.website_endpoint}"
    origin_id   = "webapp"

    custom_origin_config {
      origin_protocol_policy = "http-only"
      http_port              = 80
      https_port             = 443
      origin_ssl_protocols   = ["TLSv1.2", "TLSv1.1"]
    }
  }

  origin {
    domain_name = "${replace(replace(aws_api_gateway_deployment.api.invoke_url, "https://", ""), "/latest", "") }"
    origin_id   = "api"
    origin_path = "/latest"

    custom_origin_config {
      origin_protocol_policy = "https-only"
      http_port              = 80
      https_port             = 443
      origin_ssl_protocols   = ["TLSv1.2", "TLSv1.1", "TLSv1"]
    }
  }

  default_cache_behavior {
    compress               = true
    target_origin_id       = "landing"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  ordered_cache_behavior {
    path_pattern     = "/app*"
    target_origin_id = "webapp"
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]

    forwarded_values {
      query_string = false
      headers      = ["Authorization", "Origin"]

      cookies {
        forward = "none"
      }
    }

    compress               = true
    viewer_protocol_policy = "redirect-to-https"
  }

  ordered_cache_behavior {
    path_pattern           = "/api*"
    target_origin_id       = "api"
    allowed_methods        = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    cached_methods         = ["GET", "HEAD"]
    min_ttl                = 0
    max_ttl                = 0
    default_ttl            = 0
    compress               = true
    viewer_protocol_policy = "https-only"

    forwarded_values {
      query_string = false
      headers      = ["Authorization", "Origin"]

      cookies {
        forward = "none"
      }
    }
  }

  viewer_certificate {
    acm_certificate_arn      = "${aws_acm_certificate.cert.arn}"
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.1_2016"
  }
}

data "aws_route53_zone" "root" {
  name = "${var.DomainName}."
}

resource "aws_route53_record" "ARecordLanding" {
  type    = "A"
  name    = "${var.DomainName}"
  zone_id = "${data.aws_route53_zone.root.zone_id}"

  alias {
    name                   = "${aws_cloudfront_distribution.CloudFrontResource.domain_name}"
    zone_id                = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}

resource "aws_lambda_function" "LambdaAPIFunction" {
  filename         = "../bin/${var.AppName}-api.zip"
  function_name    = "${var.AppName}-api"
  role             = "${aws_iam_role.LambdaRole.arn}"
  handler          = "${var.AppName}-api"
  source_code_hash = "${base64sha256(file("../bin/${var.AppName}-api.zip"))}"
  runtime          = "go1.x"

  environment {
    variables = {
      DB_TABLE          = "${aws_dynamodb_table.db_table.id}"
      SNS_TOPIC_ARN     = "${aws_sns_topic.SNSTopic.arn}"
      KMS_KEY_ALIAS     = "${aws_kms_alias.KMSAlias.id}"
      SLACK_WEBHOOK_URL = "${aws_ssm_parameter.SLACK_WEBHOOK_URL.value}"
      PLAID_ENV         = "${aws_ssm_parameter.PLAID_ENV.value}"
      PLAID_PUBLIC_KEY  = "${aws_ssm_parameter.PLAID_PUBLIC_KEY.value}"
      PLAID_CLIENT_ID   = "${aws_ssm_parameter.PLAID_CLIENT_ID.value}"
      PLAID_SECRET      = "${aws_ssm_parameter.PLAID_SECRET.value}"
    }
  }
}

resource "aws_lambda_function" "dbstream" {
  filename         = "../bin/${var.AppName}-dbstream.zip"
  function_name    = "${var.AppName}-dbstream"
  role             = "${aws_iam_role.LambdaRole.arn}"
  handler          = "${var.AppName}-dbstream"
  source_code_hash = "${base64sha256(file("../bin/${var.AppName}-dbstream.zip"))}"
  runtime          = "go1.x"

  environment {
    variables = {
      SNS_TOPIC_ARN     = "${aws_sns_topic.SNSTopic.arn}"
      KMS_KEY_ALIAS     = "${aws_kms_alias.KMSAlias.id}"
      SLACK_WEBHOOK_URL = "${aws_ssm_parameter.SLACK_WEBHOOK_URL.value}"
      PLAID_ENV         = "${aws_ssm_parameter.PLAID_ENV.value}"
      PLAID_PUBLIC_KEY  = "${aws_ssm_parameter.PLAID_PUBLIC_KEY.value}"
      PLAID_CLIENT_ID   = "${aws_ssm_parameter.PLAID_CLIENT_ID.value}"
      PLAID_SECRET      = "${aws_ssm_parameter.PLAID_SECRET.value}"
    }
  }
}

resource "aws_lambda_function" "notification" {
  filename         = "../bin/${var.AppName}-notification.zip"
  function_name    = "${var.AppName}-notification"
  role             = "${aws_iam_role.LambdaRole.arn}"
  handler          = "${var.AppName}-notification"
  source_code_hash = "${base64sha256(file("../bin/${var.AppName}-notification.zip"))}"
  runtime          = "go1.x"

  environment {
    variables = {
      SLACK_WEBHOOK_URL = "${aws_ssm_parameter.SLACK_WEBHOOK_URL.value}"
      SLACK_CHANNEL     = "${aws_ssm_parameter.SLACK_CHANNEL.value}"
    }
  }
}

output "user_pool_client_id" {
  value = "${aws_cognito_user_pool_client.UserPoolClient.id}"
}

output "user_pool_id" {
  value = "${aws_cognito_user_pool.UserPool.id}"
}

output "cloudfront_distribution_id" {
  value = "${aws_cloudfront_distribution.CloudFrontResource.id}"
}

output "aws_dynamodb_table_name" {
  value = "${aws_dynamodb_table.db_table.id}"
}
