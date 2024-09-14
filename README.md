# HTTP to AWS Proxy

This is a simple proxy that takes an HTTP request and forwards it to an AWS Lambda function that is configured to receive AWS Lambda Function URL proxy events. This is helpful in local environments to emulate an AWS Lambda function URL.

## Usage

Pipe an example proxy event to the program. Pass the lambda endpoint as the first argument.

```bash
sam local generate-event apigateway aws-proxy | http-to-aws-proxy http://127.0.0.1:3001/2015-03-31/functions/FunctionName/invocations
```

Then, send HTTP requests to the proxy endpoint, `http://127.0.0.1:5544/`, and it will be forwarded to the lambda function.

```bash
cat test-event.json | curl -v -X POST --json @- "http://127.0.0.1:5544/"
```
