# Health Record Ingestion Service API Specification
The IBM Watson Health, Health Record Ingestion service is an open source project designed to serve as a “front door”, receiving health data for cloud-based solutions. See our [documentation](https://alvearie.github.io/hri/) for more details.

This repo contains the API definition for the HRI service, which consists of two parts, the management API, and Event Streams (IBM Cloud managed Kafka service). Kafka already has a defined [API](https://kafka.apache.org/documentation/), so this repo only documents the format of the Notification messages. The format of the health care data written to Kafka is not restricted by the HRI.

This is an initial publish of the code, which is being transitioned to open source, but is not yet completed. To use this code, you will need to update the CI/CD for your environment. There is a TravisCI build and integration tests, but all the references to the CD Process that deployed to internal IBM Cloud environments were removed. At some future date, a more comprehensive CI/CD pipeline will be published, as part of Watson Health's continuing Open Source initiative.

## Communication
* Please TBD
* Please see [MAINTAINERS.md](MAINTAINERS.md)

## Management API
The HRI Management API is defined using the OpenAPI 3.0 specification: [management-api/management.yml](management-api/management.yml).

### Viewing or Modifying the API
The Swagger [viewer](https://swagger.io/tools/swagger-ui/) or [editor](https://editor.swagger.io/) is an easy way to view and/or modify the API. You can also run the editor locally in a docker image and it will be available at http://localhost:80.
```
docker run -i --rm -p 80:8080 swaggerapi/swagger-editor
```

### Running a Sandbox API
You can use docker-compose to run a sandbox API that meets the specification and returns default values. Run `docker-compose up -d` from the `management-api` directory to start a sandbox API at `localhost:8000`. The sandbox endpoint is created using the [APISprout](https://github.com/danielgtaylor/apisprout) docker image.

### CI/CD
There is a TravisCI build that verifies the specification is valid and generates a Java library that defines a BatchNotification POJO. This library is available in Artifactory as a convenience for consumers of HRI Kafka notification messages. Since this is just the specification, there are no other tests or deployments. See [IBM/hri-mgmt-api](https://github.com/Alvearie/hri-mgmt-api) for implementation details of the API.

## Notification Messages
Notifications are messages written to the notification topic to notify clients of the HRI about any state changes. Notification messages are encoded in json and each type has a defined json schema. Currently, there is only one kind of notification message, batch notifications. 

### Batch Notifications
Batch notifications indicate a status change for a batch (for example, when a new batch is started or completed). The json schema is definded here: [batchNotification.json](notifications/batchNotification.json). For convenience, a Java library is available in Artifactory, which defines a BatchNotification POJO, and can be used when parsing notification messages using a JSON parsing library such as [Jackson](https://github.com/FasterXML/jackson) or [Gson](https://github.com/google/gson). The generated POJO contains Jackson2 annotations, but is compatible with both Jackson and Gson readers.
 
The generated POJO represents the notification message's date-time fields as Java8 OffsetDateTime fields, so some care must be taken when configuring either the Jackson or Gson readers. See the examples below. Also, note that schema evolution is not supported by default when using Jackson readers, unless the **FAIL_ON_UNKNOWN_PROPERTIES** feature is disabled, as shown below.

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

#### Schema Validation
There is a simple Golang program in the `notifications` directory that will validate an example notification message against the schema.  It was build with Go version 1.13 and uses Go Modules for dependency management.

##### Build
In the notification folder run:
```go build validate.go```

##### Run 
Then execute in the notification folder.  It expects the `batchNotification.json` schema to be in the current working directory.  
```
$ ./validate example1.json
The document is valid
```

## Contribution Guide
Since we have not completely moved our development into the open yet, external contributions are limited. If you would like to make contributions, please create an issue detailing the change. We will work with you to get it merged in. 
