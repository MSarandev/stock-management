# stock-management
Basic stock management (CRUD/REST)

# Running the app
1. Network declaration - `docker network create stock_api_internal`
2. Env file generation - `make init`
3. Containers stater - `make start`
4. Migrations - `make migrate`
5. gRPC protos generation - `make pb-generate`

On UNIX systems simply run `./stocks-api`

# Endpoints
## [GET] localhost:9988
Returns all entries in the db. <br> 
Example:
```json
[
    {
        "id": "8dd6a556-dde0-4bc9-b61a-b1cfd6065dbc",
        "name": "test",
        "quantity": 1,
        "created_at": "2022-09-07T15:21:50.158495Z",
        "updated_at": "2022-09-07T15:21:50.158495Z"
    },
    {
        "id": "8dd6a556-dde0-4bc9-b61a-b1cfd6065db2",
        "name": "fromPost",
        "quantity": 1,
        "created_at": "2022-09-07T18:05:23.850864Z",
        "updated_at": "2022-09-07T18:05:23.850864Z"
    },
    {
        "id": "8dd6a556-dde0-4bc9-b61a-b1cfd6065db1",
        "name": "some text",
        "quantity": 10,
        "created_at": "2022-09-07T18:05:17.488877Z",
        "updated_at": "2022-09-07T18:50:58.695888Z"
    },
    {
        "id": "8dd6a556-dde0-4bc9-b61a-b1cfd6065db4",
        "name": "fromPost",
        "quantity": 10,
        "created_at": "2022-09-07T19:00:41.335912Z",
        "updated_at": "2022-09-07T19:00:41.348239Z"
    }
]
```

## [GET] localhost:9988/{id}
Returns the specific record, or returns a validation error <br>
Example:
```json
{
    "id": "8dd6a556-dde0-4bc9-b61a-b1cfd6065dbc",
    "name": "test",
    "quantity": 1,
    "created_at": "2022-09-07T15:21:50.158495Z",
    "updated_at": "2022-09-07T15:21:50.158495Z"
}
```

Validation error:
```text
record with id: 00000000-0000-0000-0000-000000000000 doesn't exist
```

## [POST] localhost:9988
Inserts a record, accepts JSON body <br>
Request body example:
```json
{
    "id": "8dd6a556-dde0-4bc9-b61a-b1cfd6065db4",
    "name": "fromPost",
    "quantity": 10
}
```
Response: 200 OK <br>
Validation error example:
```text
Key: 'InsertStock.ID' Error:Field validation for 'ID' failed on the 'uuid4' tag
```

## [POST] localhost:9988/{id}
Updates the given record, accepts JSON body
Request body example:
```json
{
    "name": "some update",
    "quantity": 50
}
```
Both fields are optional; quantity must be > 0

Response: 200 OK <br>
Validation error example:
```text
Input quantity is less than 0
```

## [PUT] localhost:9988/{id}
Deletes a record. <br>
Validation error example:
```text
record with id: 8dd6a556-dde0-4bc9-b61a-b1cfd6065d99 doesn't exist
```