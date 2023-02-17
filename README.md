# URL-Shortener

URL-Shortener is a web application which performs the simple task of mapping a shortened url (or id) to a valid full url and redirecting the user.

<img width="494" alt="Screenshot 2022-10-23 at 12 49 41" src="https://user-images.githubusercontent.com/88666178/197388434-4d5fcf25-bfd9-4343-a663-77341e481886.png">

_The irony that the "shortened" url is significantly longer than the original is not lost on me_

## Run locally

Requirements:

- [golang](https://go.dev/)
- [aws-cli](https://aws.amazon.com/cli/)
- [sam-cli](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [docker](https://www.docker.com/)
- [npm](https://www.npmjs.com/)

### Run the API

Depending on your operating system you might have to update a configuration file to handle docker networking (see To Dos for fix)

Currently configured for mac, in `.aws-sam/development-params`, the value `DDBEndpoint=http://docker.for.mac.localhost:8000/` should be updated for:

- Linux: `http://127.0.0.1:8000`
- Windows: `http://docker.for.windows.localhost:8000/`

To run the application:

```shell
make local-start-api
```

This command runs an instance of `amazon/dynamodb-local` and the `sam` application api.

### API Endpoints

`/shorten`

POST with `url` field in body, test with:

```shell
curl --location --request POST 'http://127.0.0.1:5000/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url": "https://example.com/"
}'
```

`/short/{id}`

GET with `id` path parameter, (after populating the db with the above) test with:

```shell
curl --location --request GET 'http://127.0.0.1:5000/short/c984d06aaf'
```

### Run unit tests

Run unit tests with `make test` and get coverage report with `make test-cover`

The main `happy path` functionality is covered by these tests, but it is important to extend these cases to better coverage

### Run the Web Application

_As this is a back-end challenge, minimal effort was applied to the front-end_. Run the react app with:

```shell
npm --prefix front run start
```

## Notes

### "[...] as short as possible" vs "unique"

I chose to use a very simple id generator, the first 8 characters (configurable in `.aws-sam/{environment}-params` with `PathLength`) of an md5 hash of the sanitized url. This mechanism does not guarantee zero collision (see https://en.wikipedia.org/wiki/Birthday_attack), however given that there are no real security concerns in this use-case I believe it is a suitable solution. In the case of 10,000 entries there is a roughly 1.1% chance of collision (`1 - e^(-10000^2 / (2 * 16^8))`), if we increase this from 8 characters to 10, this chance reduces to 0.0045%.
The reason for choosing this approach is a fast, simple application that can handle concurrent requests without bloat.
To avoid any collision we could check the database, shift the slice from the original hash until a unique id is found.

### AWS SAM (Serverless Application Model)

I chose to use aws-sam to handle a lot of the application functionality, this was largely because I wanted to familiarize myself with the technology. It is very nice to be able to quickly and easily spin-up http endpoints on golang lambdas and the system is perfect for an application of this complexity.

### Hexagonal Architecture

I tried to align the design of the application to the concept of hexagonal architecture. So it should be relatively painless to switch out the storage provider (to redis for example), as well as the shortener (for a better algorithm or other solution).
See https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3

## To do

- **Introduce TTL**, in order to protect against collision we could introduce a time to live on each entry. ✅
- **Improve id generator**, either with a better hashing algorithm, database sanity checks (will add to the response time), or for fun with a word dictionary (`verb`-`noun`, e.g. `running`-`monkey`)
- **Extend unit-test coverage**, due to time constraints I only wrote `happy path` test cases for the most part
- **Acceptance tests**, to ensure the database is configured correctly and the application is working end-to-end I would like to write acceptance tests
- **Improve deployment / local**, implement automated deployments with a CI/CD process, using `docker-compose` make running locally easier
- **Handle 404s**, due to time contraints the application doesn't currently handle cases where a short url is `not found`, we could add a 404 route on the redirect function
- **Sort query parmeters**, we could reduce the number of entries in the database by sorting query parmeters
