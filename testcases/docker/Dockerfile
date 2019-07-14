FROM node:10.15.0

WORKDIR /usr/bin/app

COPY ./.babelrc .
COPY ./.env .
COPY ./.eslintrc.json .
COPY ./package.json .
COPY ./yarn.lock .
COPY ./docker/.env .
COPY ./docker/entrypoint.sh .
COPY ./test ./test

RUN yarn install

ENTRYPOINT [ "./entrypoint.sh" ]
