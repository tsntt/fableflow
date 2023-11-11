# Fableflow Api

FableFlow is a fantasy accounts and transactions API that app can connect to extend its functionalities, add user balance for viewing content, for sharing, liking, task rewarding or with real currency relation translating funds as credit units for bets or supporting creators.

1. Make a single request with domain you want to use api in and your email
2. Activate you 'bank' by clicking on the activation link on the email
3. Make requests to the endpoints direct from your client

## Registering
```javascript
const res = await fetch("https://fableflow.api/register", {
    method: "POST",
    headers: {"Content-Type": "application/json" }
    body: JSON.stringfy({
        domain: "mydomain.com",
        email: "myemail@provider.com"
        
    })
})
```

Go to your **email** and **click** the **link** send

## Accounts
### Creating
```javascript
const res = await fetch("https://fableflow.api/account", {
    method: "POST",
    headers: {"Content-Type": "application/json" }
    body: JSON.stringfy({
        initial_deposit: 100.00
    })
})

// response content
{
    header: {
        "Fableflowaid": "account token"
    },
    body: {
        "id": "67997780-86cf-4924-9507-006bd0e17f0e"
    }
}
```
### Getting
```javascript
let id = "67997780-86cf-4924-9507-006bd0e17f0e"
const res = await fetch(`https://fableflow.api/account/${id}`)

// response content
{
    body: {
        "id": "67997780-86cf-4924-9507-006bd0e17f0e",
        "balance": 100.00,
        "created_at": "2023-11-11T10:20:12z"
    }
}
```
## Transactions
### Creating
```javascript
const res = await fetch("https://fableflow.api/account/transfer", {
    method: "POST",
    headers: {"Content-Type": "application/json" }
    body: JSON.stringfy({
        sender: "uuid", // can be omited for deposit
        receiver: "uuid", // can be omited for withdraw
        amount: 100.00,
        schedule: "2023-11-11T10:10:10Z" // omit for realtime
    })
})

// response content
{
    body: {
        "id": "67997780-86cf-4924-9507-006bd0e17f0e",
        "status": "created"
    }
}
```
### Transaction processing update
```javascript
const sse = new EventSource("https://fableflow.api/account/sse")

sse.onerror = event => {
    // do yout custom logic
    console.log(event.data)
}

sse.addEventListner("transactionupdate", event => {
    // do you custom logic
    console.log(event.data)
})
```

### Cancel
Only works with unprocessed transactions
```javascript
let id = "67997780-86cf-4924-9507-006bd0e17f0e"
const res = await fetch(`https://fableflow.api/account/tranfer/${id}`, {
    method: "PATCH"
})

// response content
{
    body: {
        "id": "67997780-86cf-4924-9507-006bd0e17f0e",
        "status": "canceled"
    }
}
```

### Get transaction by ID
```javascript
let id = "67997780-86cf-4924-9507-006bd0e17f0e"
const res = await fetch(`https://fableflow.api/account/tranfer/${id}`)

// response content
{
    body: {
        "id": "67997780-86cf-4924-9507-006bd0e17f0e",
        "sender": "uuid",
        "receiver": "uuid",
        "amount": "uuid",
        "status": "*current transfer status",
        "scheduled": "2023-11-11T10:10:10z",
        "created_at": "2023-11-11T10:10:10z"
    }
}
```

### Get last 30 days account transactions
```javascript
let id = "67997780-86cf-4924-9507-006bd0e17f0e"
const res = await fetch(`https://fableflow.api/account/tranfers`)

// response content
{
    body: [
        {
            "id": "67997780-86cf-4924-9507-006bd0e17f0e",
            "sender": "uuid",
            "receiver": "uuid",
            "amount": "uuid",
            "status": "*current transfer status",
            "scheduled": "2023-11-11T10:10:10z",
            "created_at": "2023-11-11T10:10:10z"
        },
        ...
    ]
}
```

### Get account transactions by period
```javascript
let id = "67997780-86cf-4924-9507-006bd0e17f0e"
const res = await fetch(`https://fableflow.api/account/tranfers/byperiod`, {
    method: "GET",
    body: JSON.stringfy({
        "start": "2023-10-01T18:30:00Z",
        "end": "2023-11-11T18:30:00Z"
    })
})

// response content
// only transfer that are scheduled on that period
{
    body: [
        {
            "id": "67997780-86cf-4924-9507-006bd0e17f0e",
            "sender": "uuid",
            "receiver": "uuid",
            "amount": "uuid",
            "status": "*current transfer status",
            "scheduled": "2023-11-11T10:10:10z",
            "created_at": "2023-11-11T10:10:10z"
        },
        ...
    ]
}
```

## Run local
 - Postgres available
    ```bash
    docker start postgres14
	docker exec -it postgres14 createdb --username=postgres --owner=postgres fableflow_api
    ```
 - Migrate tables
     ```bash
    migrate -path data/postgres/migration -database "postgresql://postgres:postgres@localhost:5432/fableflow_api?sslmode=disable" -verbose up
    ```
 - Run server
   ```bash
   make
   # or
   make dev
   ```
 - Have fun:
 https://www.postman.com/tsnt/workspace/fableflow/collection/31015366-7aec03cc-6cc1-4d39-9f5d-f06dc4953d6e

