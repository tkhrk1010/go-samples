FROM node:14

RUN npm install -g dynamodb-admin

ENV DYNAMO_ENDPOINT=http://localstack:4566
ENV AWS_REGION=us-east-1
ENV AWS_ACCESS_KEY_ID=your-access-key-id
ENV AWS_SECRET_ACCESS_KEY=your-secret-access-key

CMD ["dynamodb-admin"]