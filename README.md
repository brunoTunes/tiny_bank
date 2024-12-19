
> [!NOTE]  
> Please disregard the low commits, this is far from ideal, and it's now how it should be done. 

# Tiny bank

A Go service to simulate a tiny bank, it exposes an REST API allowing users and their accounts to be created, deposit/withdraw balance from accounts and manage transactions between accounts.

## Getting Started

### Prerequesites

```
- Go 1.23
```

### Running the application

```
make run
```

### Running unit tests

```
make test
```

## API Specification

| Method   | URL                          | Description                                                                                                                                   | Request schema                                                   | Response schema                                                                                                        |
|----------|------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| `GET`    | `/users`                     | Fetches users, has optional return-deleted query parameter that also returns deleted users if set as true                                     |                                                                  | [{'id':'string','name':'string', 'deleted_at':'string'}]                                                               |
| `POST`   | `/users`                     | Creates new user.                                                                                                                             | {'name':'string'}                                                | {'id':'string','name':'string', 'deleted_at':'string'}                                                                 |
| `GET`    | `/users/{id}/accounts`       | Return user with {id} accounts                                                                                                                |                                                                  | [{'id':'string', 'user_id':'string', 'balance':'int', 'deleted_at':'string'}]                                          |
| `DELETE` | `/users/{id}`                | Soft Deletes user with {id}                                                                                                                   |                                                                  |                                                                                                                        |
| `GET`    | `/account/{id}/transactions` | Returns transactions from account with {id}, has from-date and to-date params in format 2006-01-02 only returning transactions in those dates |                                                                  | [{'id':'string', 'from-account':'string', 'to-account':'string', 'amount':'int', 'created_at':'string}]                |
| `POST`   | `/transaction`               | Performs a transaction from an account to another account                                                                                     | {'from_account':'string', 'to_account':'string', 'amount':'int'} | {'id':'string', 'from-account':'string', 'to-account':'string', 'amount':'int', 'created_at':'string', 'type':'string} |
| `POST`   | `/account/{id}/deposit`      | Performs a deposit to account with {id}                                                                                                       | { 'amount':'int'}                                                | {'id':'string', 'to-account':'string', 'amount':'int', 'created_at':'string', 'type':'string'}                         |
| `POST`   | `/account/{id}/withdraw`     | Performs a withdraw to account with {id}                                                                                                      | { 'amount':'int'}                                                | {'id':'string', 'from-account':'string', 'amount':'int', 'created_at':'string', 'type':'string'}                       |

> [!NOTE]  
> Delete is a soft delete

### Curl Examples

```
Create 2 users

curl --request POST \
  --url http://localhost:8080/users \
  --header 'content-type: application/json' \
  --data '{"name": "u1"}'

curl --request POST \
  --url http://localhost:8080/users \
  --header 'content-type: application/json' \
  --data '{"name": "u2"}'
```

> [!IMPORTANT]  
> in subsequent requests {u1_id}, {u2_account_id}, {u1_account_id} most be matched against generated ids which are returned in previous requests
```
Fetch u1 accounts, alongisde their balance

curl --request GET \
  --url http://localhost:8080/users/{u1_id}/accounts
  

Deposit u1 account

curl --request POST \
  --url http://localhost:8080/account/{u1_account_id}/deposit \
  --header 'content-type: application/json' \
  --data '{"amount": 100}'

Withdraw u1 account

curl --request POST \
  --url http://localhost:8080/account/{u1_account_id}/withdraw \
  --header 'content-type: application/json' \
  --data '{"amount": 100}'
 
Perform a transaction

curl --request POST \
  --url http://localhost:8080/transaction \
  --header 'content-type: application/json' \
  --data '{
  "from_account": {u1_account_id},
  "to_account": {u2_account_id},
  "amount": 100
}'

Get u1 account transactions

curl --request GET \
  --url http://localhost:8080/account/{u1_account_id}/transactions

Delete user

curl --request DELETE \
  --url http://localhost:8080/user/{u1_id} 
```
