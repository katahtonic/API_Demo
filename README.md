# API_Demo
API Demo for validating location of IP using MaxMind API

## Endpoints
All API endpoints have been organized into appropriate directories. 

**handlers/customer.go** - These endpoints are meant to represent manipulating the database for a customer table, with generic keys such as Name, Email, Phone, ID, and their restriction settings (restricted and the list of restrictions). List, Add, Edit, Delete, Gets, were all created to simplify "DB" data communication.

**handlers/user.go** - These endpoints are the primary functions for validating the location of the end user to log in. 

## User.go Endpoints
**login:** This endpoint utilizes the dummy data available in the models.go file. It represents a more integrated usage of the API with data from a DB caller, such as a customer ID or name. 
Execution: 
`POST 127.0.0.1:8080/login`
body: `{"ip": "101.167.184.6", "customer": "27d8e79f-c3c2-4a3c-a07f-8a6a019725da"}`

**loginValidation:** This endpoint offers a more "import-able" functionality, where the body expects the IP as well as a list of countries allowed. It is intended to be the more universal flow, allowing you to execute the API utilizing the raw data fields, instead of performing a customer lookup to retrieve country restriction rules. 
Execution: 
`POST 127.0.0.1:8080/loginValidation`
body: `{"ip": "101.167.184.6", "countries": ["US", "GB", "AU"]}`


## Possible Responses

Both API endpoints for user Login will respond with the same data. 
`Verified User location is approved: US` - This text will be in the HTTP response, with status OK, when the user is in a valid location to log in.

`Access denied from this location: US` - This text will be in the HTTP response, when the user is not in a valid location to log in.

`IP parameter is required` - This response will be in the HTTP response body, when there is no "ip" key in the payload body. It returns with status "Bad Request".

`At least one country is required` - This response will be in the HTTP response body, when there is no "countries" key in the payload body for the `loginValidation` API endpoint only. It returns with status "Bad Request".

`Customer parameter is required` - This response will be in the HTTP response body, when there is no "customer" key in the payload body for the `login` API endpoint only. It returns with status "Bad Request".

## Models

Models.go contains the type definitions of Customer, GeoIPResponse, as well as dummy data for testing purposes to be used with no DB or for reference of some valid IPs and their country codes. 