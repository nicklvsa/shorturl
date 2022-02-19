# ShortURL API

## Design

### Summary

The `ShortURL API` provides clear endpoints to shorten, view, and delete short urls. It also provides a convenient endpoint to retrieve metrics on a specific short url (including total views, past 24 hours, and past week). For installation instructions, see the `FAQs` section.

#### Tech used:

- __Langauge__: Golang
- __Database__: Redis
- __Containerization__: Docker / docker-compose

#### Libraries used:

- [Gin-Gonic](https://github.com/gin-gonic/gin) (Web framework)
- [Go-Redis](https://github.com/go-redis/redis) (Redis client for go)
- [redismock](https://github.com/go-redis/redismock) (Go-Redis client meant for mocking in unit tests)
- [xid](https://github.com/rs/xid) (Used to generate short url identifiers)

#### Assumptions

- Short URLs must __always__ be unique, and may never collide.
- API users must be able to retrieve short url metrics (these metrics can be easily modified, by changing the [Metrics Configuration](./metrics-config.json))
- Short URL ownership is controlled by the `employee_id`. This stands as an example how a real production environment could implement access control on certain routes that manage a short url. (Like metrics collection or deletion)
- Optional short url expiration, by using an `expires` query parameter when creating a new short url.
- All short urls are public, regardless of the creator. (via the `/v/:short_id` route)


### Endpoints
- **GET** /healthcheck
- **GET** /v/:short_id
- **GET** /short/new/:employee_id
  - Query params:
    - url=http://... (required, note: if no protocol is provided, `https` is used)
    - expires=5 (optional, in minutes)
- **GET** /short/delete/:employee_id/:short_id
- **GET** /short/metrics/:employee_id/:short_id


#### Endpoint examples:
- Create a short url that never expires:
    `http://localhost:8080/short/new/abc123?url=http://example.com`
- Create a short url that expires after 30 minutes:
    `http://localhost:8080/short/new/abc123?url=http://example.com&expires=30`
- Delete a short url:
    `http://localhost:8080/short/delete/abc123/short_url_id`
- Fetch metrics for a short url:
    `http://localhost:8080/short/metrics/abc123/short_url_id`
- View short url's long url:
    `http://localhost:8080/v/short_url_id`

#### Endpoint explanation

These endpoints are all defined as `GET` endpoints. This makes it easy to use in the browser, but also useful for a developer to integrate. In some of the above routes, `employee_id` is listed as part of path. This path variable is a unique string that represents a specific user, allowing short urls to be owned by an `employee_id`. When creating a new short url with the `/short/new/:employee_id` route, the `employee_id` may be any string. This will then be needed to reference metrics about the short url, and to delete the short url.

## FAQs

<details>
<summary>What are the required prerequisites?</summary>
<br/>
The following items are required to use this API locally:
<ul>
<li>
<i>
<a href="https://docker.com">docker</a>
</i>
</li>
<li>
<i>
<a href="https://docs.docker.com/compose/install/">docker-compose</a>
</i>
</li>
</ul>

If you are attempting to run this API outside of a containerized environment, the following items are required:
<ul>
<li>
<i>
<a href="https://golang.org">golang</a>
</i>
</li>
<li>
<i>
<a href="https://redis.io/">redis-server</a>
</i>
</li>
</ul>
</details>

<details>
<summary>How do I run it locally?</summary>
<br/>
For local use, the following command can be run to spin up the API using <strong>docker-compose</strong>:
<pre>
docker-compose up
</pre>
</details>

<details>
<summary>How do I run tests?</summary>
<br/>
To run the built in unit tests, the following script can be run from the root directory:
<pre>
./run-tests.sh
</pre>

If you are unable to run the above script, you may use the following to run the tests:
<pre>
cd tests
go test -v
</pre>
</details>

<details>
<summary>How do I collect more metrics?</summary>
<br/>
1. Open the <a href="./metrics-config.json">Metrics Configuration</a> json file.<br/>
2. Modify the <strong><i>periods</i></strong> object to collect customized metrics.

<h4>Example:</h4>
Example addition which collects 1 hour of short url link counts
<pre>
{
    "1h": "past_hour"
}
</pre>
</details>
