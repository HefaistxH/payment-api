# MNC Tech-Test

## Table of Contents

- [About](#about)
- [Getting Started](#getting-started)
- [Documentation](#documentation

## About 

> I'm focusing on listed required features which are login, payment, and logout in this project. The table for customer and merchant is separated for effectiveness on data management as the users base will likely keep growing. JWT is being used for authentication purpose, although it's only a simple implementation, it still works in this type of scenario. As for now, payment is the only activity that can be logged into history table and i'm not yet to make an endpoint for it. I do log errors from server-side/repository too using logrus, you can check the logrus.log on root branch for error checking. I hope this meets the required criteria, even if it's not i do really hope to get some feedback for me to reflect and learn from my own mistakes. It is recomended to open this .md file using relevant text editor (you can copy this file to readme.so for more convenience).

## Getting Started 

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 
### Prerequisites
GO Language
Postgres
Postman

### Installation
Run the ddl and dml query below
Use json data on documentation section for functionality checking
#
#### DDL query:
```sh
CREATE DATABASE mnctechtestdb;
```
```sh
CREATE TYPE role_enum AS ENUM ('merchant', 'customer');
CREATE TABLE credentials (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARHCAR(36) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role role_enum NOT NULL
);
```
```sh
CREATE TABLE customers (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL
);
```
```sh
CREATE TABLE merchants (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL
);
```
```sh
CREATE TABLE history (
    id VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    merchant_id VARCHAR(36) NOT NULL,
    activity VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    message TEXT,
    timestamp TIMESTAMP NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);
```

#### DML Query:
```sh
INSERT INTO customers (id, name, balance) VALUES 
('cus-001', 'Test User', 1000.00);
```
```sh
INSERT INTO merchants (id, name, balance) VALUES 
('mer-001', 'Test Merchant', 5000.00);
```
```sh
UPDATE credentials set password = '$2a$12$7fb6RY3cGve07FAjf51gsOlfw0bbowS2JZ5WSCDmJKDtkwZ/JoOx6'
WHERE user_id = 'cus-001'
```
#
#
## Documentation
#
#### Login: 
> Get the accessToken and refreshToken for accesing the payment feature
#### Endpoint: 
> POST localhost:8080/api/v1/login
#### Request Body:
```sh
{
 "email": "testuser@example.com",
 "password": "Password123"
}
```
#### Response Body:
```sh
{
 "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
 "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
#
#
#### Payment: 
> Transfer some balances from customer to merchant
#### Endpoint: 
> POST localhost:8080/api/v1/payment
#### Header: 
> Authorization: Bearer <refreshToken>
#### Request Body:
```sh
{
 "customer_id": "cus-001",
 "merchant_id": "mer-001",
 "pin": "123456",
 "amount": 100.00,
 "message": "Payment for services"
}
```
#### Response Body:
```sh
{
"payment_id": "pay-001",
"customer_id": "cus-001",
"merchant_id": "mer-001",
"amount": 100.00,
"status": "Success",
"timestamp": "2024-11-26T03:55:57+07:00"
}
```
#
#
#### Logout: 
> Invalidate token, prevents it from being used
#### Endpoint: 
> POST localhost:8080/api/v1/logout
#### Header: 
> Authorization: Bearer <refreshToken>
#
#### History: Activity logging purposes
> No endpoint yet, but you can use this query below for checking table's content
```sh
SELECT * FROM history;
```

