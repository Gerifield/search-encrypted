# search-encrypted
Experimental repo to test searching on encrypted data


**_It is NOT production ready or secure, just an experiment!_**

So this script creates an in-memory SQLite db and inserts some data in an AES encrypted format and then tried to find some columns based on pre-generated indexes.
(The indexes are also encrypted!)

There are very good articles about this topic, I used some ideas for this example from here: https://paragonie.com/blog/2017/05/building-searchable-encrypted-databases-with-php-and-sql

Example run:

```bash
$ go run main.go
2021/10/16 14:58:45 ID: 1 Data: {John Carter}
2021/10/16 14:58:45 ID: 1 Encrypted: GwLD0UQfCx9BR3mgQ5m656sZlNNW4tvTN+g+8vaMoPh8JLjayJtW2KeSZMykJ9kxC43vduN4hitKgUch34u8hICBGweZdQ== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3
2021/10/16 14:58:45 ID: 2 Data: {John Doe}
2021/10/16 14:58:45 ID: 2 Encrypted: 835loJ7sVusBy5CrlflMyd4eePbGv8A3IOJ2zLXQaGZYbPoKnzn8c1rbUkVbRac5QxYu34FjBdweYSaY+kDqTvqu2g== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3
2021/10/16 14:58:45 ID: 3 Data: {John Wick}
2021/10/16 14:58:45 ID: 3 Encrypted: K1pt0wMQApaSoOKl4Z7DLmLt+o7GxRUPd94/AExz7U8+z/Tn9T7f8hU3122ebBQQMYbAw3y4u2Vw1kaiabsAgNFPbrY= FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3
2021/10/16 14:58:45 ID: 4 Data: {Johnatan Somebody}
2021/10/16 14:58:45 ID: 4 Encrypted: MC+HqL4FRIyTl+9pOMSAg58X2GIh/9vIT85fOJHPWhpLTrQtKGcQIil6kj3+zur1A37QRy+Pzvwtpg++0MUrtJ5odGHRMiG33bPzJQ== FirstName idx: f3415b01e81a1477a918458e072655d5381b623dd6361f3531ec0d84c58cd511
2021/10/16 14:58:45 ID: 5 Data: {Somebody Else}
2021/10/16 14:58:45 ID: 5 Encrypted: FBbJ56WDhXR6BfXseZ+eDPCK9lDCxZj6vkT0l4S+aniAGBLGVFizjpASZEV8RMGQPDj7yTd/UkhFOW0z2lf7DB9DszXKMKDN FirstName idx: 8c9f971a8c5797dab1a47a6d9be056fc78e44415f9b0de7c54e6ea362afb7488
2021/10/16 14:58:45 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 14:58:45 Test some data:
2021/10/16 14:58:45 ID: 1 Encrypted: GwLD0UQfCx9BR3mgQ5m656sZlNNW4tvTN+g+8vaMoPh8JLjayJtW2KeSZMykJ9kxC43vduN4hitKgUch34u8hICBGweZdQ==
2021/10/16 14:58:45 ID: 1 Data: {John Carter}
2021/10/16 14:58:45 ID: 2 Encrypted: 835loJ7sVusBy5CrlflMyd4eePbGv8A3IOJ2zLXQaGZYbPoKnzn8c1rbUkVbRac5QxYu34FjBdweYSaY+kDqTvqu2g==
2021/10/16 14:58:45 ID: 2 Data: {John Doe}
2021/10/16 14:58:45 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 14:58:45 Do some index search:
2021/10/16 14:58:45 Looking for John
2021/10/16 14:58:45 Result, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 14:58:45 Result, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 14:58:45 Result, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 14:58:45 Looking for Somebody
2021/10/16 14:58:45 Result, ID: 5 Data: {"first_name":"Somebody","last_name":"Else"}
2021/10/16 14:58:45 Looking for Nobody
```

