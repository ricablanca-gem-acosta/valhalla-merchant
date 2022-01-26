# valhalla-merchant

**Starting the App**
Run the command below from the project's root directory. Note that the database will be populated with dummy data on first run.
`go run main/main.go`

**Testing the App**
Run the command below from the project's root directory.
`go test ./...`

**URL Structures:**

Get the list of merchant codes

`GET localhost:3000/merchant`

Add a new merchant

`POST localhost:3000/merchant`

*JSON body param* `{"code": "{merchant code}"`

Delete a merchant

`DELETE localhost:3000/merchant/{merchant code}`

Add a merchant member

`POST localhost:3000/merchant/{merchant code}/addmember`

*JSON body param* `{"email": "{member's email}"`

Delete a merchant member

`DELETE localhost:3000/merchant/{merchant code}/{member's email}`

Get merchant members

`GET localhost:3000/merchant/{merchant code}/members?page={page}`

**Returns**

`page` - the current page number

`totalPages` - the total number of pages

`count` - total number of members on the page

`data` - array of member data

```
\\ sample
{
    "page": 1,
    "totalPages": 10,
    "count": 100
    "data": [
        {
            "email": "orvilleturner@beahan.org"
        },
        {
            "email": "olivercollins@murray.biz"
        }
    ]
}
```