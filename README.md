# PiePay

**Tech Used**
* Golang 
* Sentry
* Elasticsearch

**Postman Collection**
* [Import from the link](https://www.getpostman.com/collections/09fbd18d2eee5e446a64)

**How to start the project?**
* ```go run main.go```

**Basic Functionalities**
* Goroutine to fetch data from Youtube API every 10 second and upload on Elasticsearch.
* ```/get``` to get paginated data from Elasticsearch sorted on basis of dates(latest first).
* ```/search``` to search data on basis of description and title ( ```A video with 'title How to make tea?'  matches for the search query 'tea how'```)  
* Multiple key support (in case one expires)
