# Project Name

This is a CRM Backend Project that provides a set of APIs for managing customer data. This is a final project for Go Language (Golang) course at Udacity. 

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/gkj/golang.git
    ```

2. Change to the project directory:

    ```bash
    cd golang
    ```

3. Install the dependencies:

    ```bash
    go mod download
    ```

4. Run the service:

    ```bash
    go run main.go
    ```
    

## Usage

To use the APIs, you can perform CURL requests to the following endpoints:

### Getting all customers
Endpoint: `/customers`

Method: `GET`

Description: Retrieves all customers.
```bash
curl --location 'http://localhost:3000/customers' \
--data ''
```

### Getting a single customer
Endpoint: `/customers/{id}`

Method: `GET`

Description: Retrieves a single customer by their ID.
```bash
curl --location 'http://localhost:3000/customers/1' \
--data ''
```

### Creating a customer
Endpoint: `/customers`

Method: `POST`

Description: Creates a new customer.
```bash
curl --location 'http://localhost:3000/customers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": 4,
    "name": "Gilang",
    "role": "User",
    "email": "gilang@gmail.com",
    "phone": "098763454321"
}'
```

### Updating a customer
Endpoint: `/customers/{id}`

Method: `PUT`

Description: Updates an existing customer by their ID.
```bash
curl --location --request PUT 'http://localhost:3000/customers/3' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": 3,
    "name": "Gilang Kusuma",
    "role": "User",
    "email": "gilang@gmail.com",
    "phone": "098763454321"
}'
```

### Deleting a customer
Endpoint: `/customers/{id}`

Method: `DELETE`

Description: Deletes a customer by their ID.
```bash
curl --location --request DELETE 'http://localhost:3000/customers/3' \
--data ''
```