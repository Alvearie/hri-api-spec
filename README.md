# Health Record Ingestion Service API Specification
The Alvearie Health Record Ingestion service: a common 'Deployment Ready Component' designed to serve as a “front door” for data for cloud-based solutions. See our [documentation](https://alvearie.io/HRI/) for more details.

This repo contains the API definition for the HRI service, which consists of two parts, the management API, and Event Streams (IBM Cloud managed Apache Kafka service). Kafka already has a defined [API](https://kafka.apache.org/documentation/), so this repo only documents the format of the Notification messages. The format of the health care data written to Kafka is not restricted by the HRI, but there is an Alvearie [FHIR Implementation Guide](https://github.com/Alvearie/alvearie-fhir-ig) project for using FHIR as the common format.  

## Communication
* Please [join](https://alvearie.io/contributions/requestSlackAccess) our Slack channel for further questions: [#health-record-ingestion](https://alvearie.slack.com/archives/C01GM43LFJ6)
* Please see recent contributors or [maintainers](MAINTAINERS.md)

## Management API
The HRI Management API is defined using the OpenAPI 3.0 specification: [management-api/management.yml](management-api/management.yml).

### Viewing or Modifying the API
The Swagger [viewer](https://swagger.io/tools/swagger-ui/) or [editor](https://editor.swagger.io/) is an easy way to view and/or modify the API. You can also run the editor locally in a docker image, and it will be available at http://localhost:80.
```
docker run -i --rm -p 80:8080 swaggerapi/swagger-editor
```

### Authentication & Authorization
All endpoints (except the health check) require an OAuth 2.0 JWT bearer access token per [RFC8693](https://tools.ietf.org/html/rfc8693) in the `Authorization` header field. The Tenant and Stream endpoints require IBM Cloud IAM tokens, but the Batch endpoints require a token with HRI and Tenant scopes for authorization. The Batch token issuer is configurable via a bound parameter, and must be OIDC compliant because the code dynamically uses the OIDC defined well know endpoints to validate tokens. Integration and testing have already been completed with App ID, the standard IBM Cloud solution.

In the API specification, the required tokens for a given endpoint are specified in the `security` section. `bearerAuth` is used to specify an IAM token. `hri_auth` is used to specify a JWT token from the configured issuer, and additionally it will list the required scopes. If multiple combinations of scopes allow access, each combination is listed. For example, the `getBatches` (GET '/tenants/{tenantId}/batches') endpoint allows two different combinations of scopes: (`tenant_tenantId` and `hri_data_integrator`) or (`tenant_tenantId` and `hri_consumer`). The specification appears as
```yaml
      security:
        - hri_auth: [tenant_tenantId, hri_data_integrator]
        - hri_auth: [tenant_tenantId, hri_consumer]
```

**NOTE:** the `tenant_tenantId` scope is relative to the tenant being accessed, where `tenantId` must match the endpoint's path parameter. For example, when calling the `getBatches` endpoint for tenant 'provider1234' (GET '/tenants/provider1234/batches'), the scopes must include `tenant_provider1234`.

#### Data Integrator Filtering
Data Integrators are only allowed to access Batches that they created. So if only the `hri_data_integrator` scope is used to access an endpoint, the returned data will be filtered. If a request includes the `hri_consumer` scope, results will not be filtered.

### Running a Sandbox API
You can use docker-compose to run a sandbox API that meets the specification and returns default values. Run `docker-compose up -d` from the `management-api` directory to start a sandbox API at `localhost:8000`. The sandbox endpoint is created using the [APISprout](https://github.com/danielgtaylor/apisprout) docker image.

## CI/CD
The GitHub actions build verifies the specification is valid and generates a Java library that defines a BatchNotification POJO. This library is available in **TBD** as a convenience for consumers of HRI Kafka notification messages. Since this is just the specification, there are no other tests or deployments. See [Alvearie/hri-mgmt-api](https://github.com/Alvearie/hri-mgmt-api) for implementation details of the API.

## Releases
Releases are created by creating Git tags, which trigger a build that publishes a release version in **TBD**, see [Overall strategy](https://github.com/Alvearie/HRI/wiki/Overall-Project-Branching,-Test,-and-Release-Strategy) for more details.

## Notification Messages
Notifications are messages written to the notification topic to notify clients of the HRI about any state changes. Notification messages are encoded in json and each type has a defined json schema. Currently, there is only one kind of notification message, batch notifications. 

### Batch Notifications
Batch notifications indicate a status change for a batch. For example, when a new batch is started or completed. The json schema is defined here: [batchNotification.json](notifications/batchNotification.json). For convenience, a Java library is available in **TBD**, which defines a BatchNotification POJO, and can be used when parsing notification messages using a JSON parsing library such as [Jackson]() or [Gson](). The generated POJO contains Jackson2 annotations, but is compatible with both Jackson and Gson readers. The generated POJO represents the notification message's date-time fields as Java8 OffsetDateTime fields, so some care must be taken when configuring either the Jackson or Gson readers. See the examples below. Also, note that schema evolution is not supported by default when using Jackson readers, unless the **FAIL_ON_UNKNOWN_PROPERTIES** feature is disabled, as shown below.

Jackson Reader:  
Requires additional dependency: *com.fasterxml.jackson.datatype:jackson-datatype-jsr310*
```Java
ObjectMapper mapper = JsonMapper.builder()
        .addModule(new JavaTimeModule())
        .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
        .build();
```

Gson Reader:  
Requires additional dependency: *com.fatboyindustrial.gson-javatime-serialisers:gson-javatime-serialisers*
```Java
Gson gson = Converters.registerOffsetDateTime(new GsonBuilder()).create();
```

### Invalid Records
When validation is enabled, the Flink validation processing will output 'invalid' records to a `*.invalid` topic. These records contain a reference to the original record (topic, partition, offset) and a failure message. It does not contain the original record. The json schema is defined here: [invalidRecord.json](notifications/invalidRecord.json).

## Contribution Guide
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.
