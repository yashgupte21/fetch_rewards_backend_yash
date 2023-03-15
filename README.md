# fetch_rewards_backend_yash
Fetch Rewards Backend Take Home - Yash Pradeep Gupte

Name: Yash Pradeep Gupte \
Email: ygupte@hawk.iit.edu

## Table of Contents
* [Technologies ](#technologies)
* [Installation Guide ](#installation-guide)
* [Assumptions ](#assumptions)

---
## Technologies
Project is developed with:
* Mac Machine
* Go: go1.20.1 darwin/arm64
--- 
## Installation Guide 

I have used GoLang to implement my solution 

To install golang on your system follow the steps provided below:

1. Go to https://go.dev/dl/ and download golang for your system in your root or HOME directory

2. Verify if golang has been installed. Open terminal and type the following command

``` 
go version 
```

3. Clone this repository in your root location 

4. The module files have been already added to the respository, but incase you want to install / update modules execute the following command 

```
go mod init github.com/fetch_rewards_backend_yash
go mod tidy
```

5. Next step is to run main.go program 

```
go run main.go
```

6. Once the code is executed, open the following url on your local browser

`http://127.0.0.1:8000/receipts/process`

Output:
```json
{"id":"3cbec19d-7515-4d40-96ae-8eeab4c4d996"} 
```

The above API call will generate a unique id for a particular receipt 

Next, in a new tab open the following url to get points received by the  corresponding receipt id 

`http://127.0.0.1:8000/receipts/{id}/points`

Example: \
`http://127.0.0.1:8000/receipts/3cbec19d-7515-4d40-96ae-8eeab4c4d996/points`

Outputs points:
```json
{"points":109}
```

Executing this program will Extract AWS SQS Queue Messages from Localstack, convert them into JSON format, Transform the data as required and Load it on the Postgres database docker container.

--- 

## Assumptions
 * In memory solution has been presented utilizing go cache to store json data fetched from API and mapped with corresponding id
 * id : To generate id for each receipt, I have implemented UUID generation. A UUID is a 128 bit or 16 byte Universal Unique IDentifier. Each receipt is assigned with UUID and stored in cache maintaing mapping between id and json object of the receipt
 * points : generated based on rules sepcified and is implemented in a method which takes the id from the URL, fetches json object, and returns points for a receipt