REST API Shopping Cart in GO with Tracing and Metrics
---

### Project Features
* This project is an implementation of a Cart for an e-commerce and has 2 endpoints:
    - One to add items
    - One to retrieve the cart status
* The project has been implemented with DDD and Hexagonal Arquitecture, isolating domain, application and infrastructure
* In the root of internal the package cart has domain entities and Value objects of the aggregate 
* There are two application services, in the package creator and getter
* In storage there is an inMemoryRepository for the two entities that could be easily modified to a Real DB
* In the io/rest package there is the http server and the handler 
* The only two libraries used that are not from the standard library are:
    - Gorilla Mux
    - Wire as the dependency injector
* The application is Dockerized and can be launched from a makefile

* The Server is launched in the port :8888
* Metrics are pulled by Prometheus and can be seen in http://localhost:9090
* Traces are retrieved by Jaeger and can be seen in http://localhost:16686


### How to use the application in your local machine
##### Launch the application server
```
make up
```
##### Stop the application
```
make down
```

##### Rebuild the application if changes are done to the code
```
make build
```

* Only 3 kinds of items can be added to the cart:
    - "book"
    - "dvd"
    - "casette"


##### To insert new items with a Curl Command 
```
curl -i \
-H "Accept: application/json" \
-H "Content-Type:application/json" \
-X POST --data '{"item_id":"book","quantity":4}' http://localhost:8888/
```

##### To get the current cart status
```
curl -i -X GET http://localhost:8888/
```

There are unit tests and acceptance, in future iterations also integration should be added.
To launch the tests:
```
make test
```