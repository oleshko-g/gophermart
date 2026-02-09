# Gophermart

## What it is

- A loyalty points balance tracking service.

## Why it was built

- It's a required milestone project on the "Advanced Go Programmer" course by Yandex
  Practicum.

## What it does

It's implemented by two services.

### User

A new user can:

- Register
- Authenticate

### Balance

An authenticated user can:

- Upload an "accrual" order to track its accrual of loyalty points.
- List their uploaded orders
- Check their balance
- Withdraw points from the balance
- List their withdrawals

## It's built with

- goa.design to handcraft the design of HTTP API and generate the HTTP server
- PostgreSQL as the storage
- goose to run SQL migrations
- sqlc to generate the code to serialize data to and from a PostgreSQL DB
