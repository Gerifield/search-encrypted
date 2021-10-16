# search-encrypted
Experimental repo to test searching on encrypted data


**_It is NOT production ready or secure, just an experiment!_**

So this script creates an in-memory SQLite db and inserts some data in an AES encrypted format and then tries to find some columns based on pre-generated indexes.
The indexes are also encrypted and there are different approaches how to interact with them.

There are very good articles about this topic, I used some ideas for this example from here: https://paragonie.com/blog/2017/05/building-searchable-encrypted-databases-with-php-and-sql

Example run:

```bash
$ go run main.go
2021/10/16 15:16:58 ID: 1 Data: {John Carter}
2021/10/16 15:16:58 ID: 1 Encrypted: bmYg+jda9pXN+MOMdSXvtknab3HgEr94Gweqzkdj9e1W4IN4YXtRzILwz0rRydC9NkUr0hoA9gQqPip5Wwbhpjz/RyWPXg== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:16:58 ID: 2 Data: {John Doe}
2021/10/16 15:16:58 ID: 2 Encrypted: OBdNtBTZRHf32XZ3wTaoCwLPZk7qC+Pp2jJXT4OkmFRC3xUGJd0Am+oTaE/dZ2Aejws58g20UonXZYbpIzj3fN/0Zw== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:16:58 ID: 3 Data: {John Wick}
2021/10/16 15:16:58 ID: 3 Encrypted: tv3RNV3/HBGPUx0GDEDRAh2cbRmrMW7D3/C1DqGOdUVzv8x4vUcs4yUY9Bf/1LZH4EwAiEmPmnRiGggjVkAwLCt/qHA= FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:16:58 ID: 4 Data: {Johnatan Somebody}
2021/10/16 15:16:58 ID: 4 Encrypted: /exsgBsaq2E1e1BnO2se5Itz1W+sJ6+TocOdEAt0sxMaExzfLn5ijW9FPWLKnR9Sw/l+3mT1tyKQcUpPqnVen40cYWqMThOTBesblA== FirstName idx: f3415b01e81a1477a918458e072655d5381b623dd6361f3531ec0d84c58cd511 FirstName bloom: f3415b01e81a1477a918458e072655d5
2021/10/16 15:16:58 ID: 5 Data: {Somebody Else}
2021/10/16 15:16:58 ID: 5 Encrypted: w+LQbZeD713dIU5JSp1+XK/J4e4O7b6krVQVBZZbvkEA1tjW9VKwPmLZL0LEgi/3lsZuz8fjDoPCyGiQVqnKJtcDM73s+rk9 FirstName idx: 8c9f971a8c5797dab1a47a6d9be056fc78e44415f9b0de7c54e6ea362afb7488 FirstName bloom: 8c9f971a8c5797dab1a47a6d9be056fc
2021/10/16 15:16:58 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:16:58 Test some data:
2021/10/16 15:16:58 ID: 1 Encrypted: bmYg+jda9pXN+MOMdSXvtknab3HgEr94Gweqzkdj9e1W4IN4YXtRzILwz0rRydC9NkUr0hoA9gQqPip5Wwbhpjz/RyWPXg==
2021/10/16 15:16:58 ID: 1 Data: {John Carter}
2021/10/16 15:16:58 ID: 2 Encrypted: OBdNtBTZRHf32XZ3wTaoCwLPZk7qC+Pp2jJXT4OkmFRC3xUGJd0Am+oTaE/dZ2Aejws58g20UonXZYbpIzj3fN/0Zw==
2021/10/16 15:16:58 ID: 2 Data: {John Doe}
2021/10/16 15:16:58 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:16:58 Do some index search:
2021/10/16 15:16:58 Looking for John
2021/10/16 15:16:58 Result, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:16:58 Result, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:16:58 Result, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:16:58 Looking for Somebody
2021/10/16 15:16:58 Result, ID: 5 Data: {"first_name":"Somebody","last_name":"Else"}
2021/10/16 15:16:58 Looking for Nobody
2021/10/16 15:16:58 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:16:58 Do some bloom filter search:
2021/10/16 15:16:58 Looking for John
2021/10/16 15:16:58 Result MATCHED, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:16:58 Result MATCHED, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:16:58 Result MATCHED, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:16:58 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:16:58 Do some bloom filter search with secondary filtering:
2021/10/16 15:16:58 Looking for John Doe
2021/10/16 15:16:58 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:16:58 Result matched, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:16:58 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:16:58 Looking for John Carpenter
2021/10/16 15:16:58 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:16:58 Result NOT matched, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:16:58 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
```

