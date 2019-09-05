# our-data-go/nws

Interact with the United States National Weather Service (NWS) API Web Service in Go. 

## Introduction

Although many companies provide web APIs for weather data, most of them are limited or unavailable without a fee. The NWS analyze a full range of weather data and is available for free at [weather.gov](https://www.weather.gov) via several APIs. The goal of this package is to provide a convenient way to access those data in Go, with a focus on use cases centered around a single point on Earth, such as home automation. 

The "API Web Service" used in this package is (somewhat) documented [here](https://www.weather.gov/documentation/services-web-api) and [here](https://forecast-v3.weather.gov/documentation). 

## State of the API and this package

**WARNING!:** This project is a work in progress and absolutely not production ready. The exported interface may be unstable.

The API that this package interacts with is itself a work in progress. According to the specification tab [here](https://www.weather.gov/documentation/services-web-api), only `alerts/*` endpoints are currently considered operational by the NWS.

The client is built around around a specific point on Earth and all methods operate based on that point. See comments on the code or the `godoc` for more. 

## License

Please see the `LICENSE` file in this directory.


## Todo

- [ ] write tests
