# search-encrypted
Experimental repo to test searching on encrypted data


**_It is NOT production ready or secure, just an experiment!_**

So this script creates an in-memory SQLite db and inserts some data in an AES encrypted format and then tries to find some columns based on pre-generated indexes.
The indexes are also encrypted and there are different approaches how to interact with them.

There are very good articles about this topic, I used some ideas for this example from here: https://paragonie.com/blog/2017/05/building-searchable-encrypted-databases-with-php-and-sql

Example run:

```bash
$ go run main.go
2021/10/16 15:21:55 ID: 1 Data: {John Carter}
2021/10/16 15:21:55 ID: 1 Encrypted: L9s5hsjpEqA6ej3bEO9NX+ZxoQDCg6W8rQZBZsoivFzPFvhwDnuxm2TC2iGfH4sY9L11Eod81gkg0r2Uakj1RfI+rTIkHg== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:21:55 ID: 2 Data: {John Doe}
2021/10/16 15:21:55 ID: 2 Encrypted: rdpI1pM4iMMm/RG+VpSR5Ddrz7tWjWKuPuGIo0xxN9vx9JwBqCoJakJR8l1+VF1KufRu+TB8T76MIMSqIsO5z5m0yA== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:21:55 ID: 3 Data: {John Wick}
2021/10/16 15:21:55 ID: 3 Encrypted: 1m+/eEB3zZhGyg76tdngQBuN6pi8+pWPz3m5y5KdUzgGD2/98aCSsMr7rbsU+j8WjxopPOYNlFLVeigRPLHQJfxJPYc= FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:21:55 ID: 4 Data: {Johnatan Somebody}
2021/10/16 15:21:55 ID: 4 Encrypted: B+h71quIeCtJikxN9LUKeR1VkjpWnRSa2bvFziW+eQAoPcP7o0jg4SI5npnJGRLCuuZHndlgLnmoc666WXOpaUe8sSs1rLom2BrYKw== FirstName idx: f3415b01e81a1477a918458e072655d5381b623dd6361f3531ec0d84c58cd511 FirstName bloom: f3415b01e81a1477a918458e072655d5
2021/10/16 15:21:55 ID: 5 Data: {Somebody Else}
2021/10/16 15:21:55 ID: 5 Encrypted: 95hGllU4sJq7Wjt0A4NwWqS1Lp0JirW+zfKkRcbzh6raPBvqq5PMHVhKHzOwAhmDIli0Eqsj9HPXUHgF0pgEGb2pJEa/R422 FirstName idx: 8c9f971a8c5797dab1a47a6d9be056fc78e44415f9b0de7c54e6ea362afb7488 FirstName bloom: 8c9f971a8c5797dab1a47a6d9be056fc
2021/10/16 15:21:55 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:21:55 Fetch some data via index:
2021/10/16 15:21:55 ID: 1 Encrypted: L9s5hsjpEqA6ej3bEO9NX+ZxoQDCg6W8rQZBZsoivFzPFvhwDnuxm2TC2iGfH4sY9L11Eod81gkg0r2Uakj1RfI+rTIkHg==
2021/10/16 15:21:55 ID: 1 Data: {John Carter}
2021/10/16 15:21:55 ID: 2 Encrypted: rdpI1pM4iMMm/RG+VpSR5Ddrz7tWjWKuPuGIo0xxN9vx9JwBqCoJakJR8l1+VF1KufRu+TB8T76MIMSqIsO5z5m0yA==
2021/10/16 15:21:55 ID: 2 Data: {John Doe}
2021/10/16 15:21:55 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:21:55 Do some blind index search:
2021/10/16 15:21:55 Looking for John
2021/10/16 15:21:55 Result, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:21:55 Result, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:21:55 Result, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:21:55 Looking for Somebody
2021/10/16 15:21:55 Result, ID: 5 Data: {"first_name":"Somebody","last_name":"Else"}
2021/10/16 15:21:55 Looking for Nobody
2021/10/16 15:21:55 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:21:55 Do some bloom filter search:
2021/10/16 15:21:55 Looking for John
2021/10/16 15:21:55 Result MATCHED, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:21:55 Result MATCHED, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:21:55 Result MATCHED, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:21:55 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:21:55 Do some bloom filter search with secondary filtering:
2021/10/16 15:21:55 Looking for John Doe
2021/10/16 15:21:55 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:21:55 Result MATCHED, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:21:55 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:21:55 Looking for John Carpenter
2021/10/16 15:21:55 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:21:55 Result NOT matched, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:21:55 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
```

