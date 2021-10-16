# search-encrypted
Experimental repo to test searching on encrypted data


**_It is NOT production ready or secure, just an experiment!_**

So this script creates an in-memory SQLite db and inserts some data in an AES encrypted format and then tries to find some columns based on pre-generated indexes.
The indexes are also encrypted and there are different approaches how to interact with them.

There are very good articles about this topic, I used some ideas for this example from here: https://paragonie.com/blog/2017/05/building-searchable-encrypted-databases-with-php-and-sql

Example run:

```bash
$ go run main.go                                                                                                                                                                           130 ↵
2021/10/16 15:20:09 ID: 1 Data: {John Carter}
2021/10/16 15:20:09 ID: 1 Encrypted: ifhvqlwCKhvU83Upk+BMimJLhxgH+hyY4U/WwdrGzHTiB0BWPAqTzeDC2WRUCvCOejkXnxpi0DK1f5e55hQ4eeeSk8v9XQ== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:20:09 ID: 2 Data: {John Doe}
2021/10/16 15:20:09 ID: 2 Encrypted: 01Gpy8SHUHenPIudys7E6t0reWPSEu5Z72y1H75Il0Op/I14EprKfyNQC7LRZktXEGrgCb8f5mzp0WSEXb0EjfbKdg== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:20:09 ID: 3 Data: {John Wick}
2021/10/16 15:20:09 ID: 3 Encrypted: LBXPQV7Uud2jE0eZDTuUZaCRoYkddiI2Uvd9VuP0NB0xLDU4OISGSGKySe7hZlN4e3WjUh01RAna7UfmDtpGTAYzCoU= FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f
2021/10/16 15:20:09 ID: 4 Data: {Johnatan Somebody}
2021/10/16 15:20:09 ID: 4 Encrypted: e930X9BeFvz0TiSPo6GOs7ADQwDpErrbkfHUJ+JYU1UZbSq3cBpmEylruCnKDHXuXOT4/yB4YzGFFk/Ul+u8Jy/AIGNQtrhnLnBqbw== FirstName idx: f3415b01e81a1477a918458e072655d5381b623dd6361f3531ec0d84c58cd511 FirstName bloom: f3415b01e81a1477a918458e072655d5
2021/10/16 15:20:09 ID: 5 Data: {Somebody Else}
2021/10/16 15:20:09 ID: 5 Encrypted: 2na76YO1uBnZImo/Sk/ICV3a1QK3Bo2tNfLAUMHWI4MMspmkKtSYDMnSVhSM2pbbAeb8wiytQUcn9aYzxHThWPzMtnL1Vun+ FirstName idx: 8c9f971a8c5797dab1a47a6d9be056fc78e44415f9b0de7c54e6ea362afb7488 FirstName bloom: 8c9f971a8c5797dab1a47a6d9be056fc
2021/10/16 15:20:09 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:20:09 Fetch some data via index:
2021/10/16 15:20:09 ID: 1 Encrypted: ifhvqlwCKhvU83Upk+BMimJLhxgH+hyY4U/WwdrGzHTiB0BWPAqTzeDC2WRUCvCOejkXnxpi0DK1f5e55hQ4eeeSk8v9XQ==
2021/10/16 15:20:09 ID: 1 Data: {John Carter}
2021/10/16 15:20:09 ID: 2 Encrypted: 01Gpy8SHUHenPIudys7E6t0reWPSEu5Z72y1H75Il0Op/I14EprKfyNQC7LRZktXEGrgCb8f5mzp0WSEXb0EjfbKdg==
2021/10/16 15:20:09 ID: 2 Data: {John Doe}
2021/10/16 15:20:09 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:20:09 Do some blind index search:
2021/10/16 15:20:09 Looking for John
2021/10/16 15:20:09 Result, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:20:09 Result, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:20:09 Result, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:20:09 Looking for Somebody
2021/10/16 15:20:09 Result, ID: 5 Data: {"first_name":"Somebody","last_name":"Else"}
2021/10/16 15:20:09 Looking for Nobody
2021/10/16 15:20:09 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:20:09 Do some bloom filter search:
2021/10/16 15:20:09 Looking for John
2021/10/16 15:20:09 Result MATCHED, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:20:09 Result MATCHED, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:20:09 Result MATCHED, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:20:09 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/16 15:20:09 Do some bloom filter search with secondary filtering:
2021/10/16 15:20:09 Looking for John Doe
2021/10/16 15:20:09 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:20:09 Result matched, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:20:09 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
2021/10/16 15:20:09 Looking for John Carpenter
2021/10/16 15:20:09 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter"}
2021/10/16 15:20:09 Result NOT matched, ID: 2 Data: {"first_name":"John","last_name":"Doe"}
2021/10/16 15:20:09 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick"}
```

