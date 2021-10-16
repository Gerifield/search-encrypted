# search-encrypted
Experimental repo to test searching on encrypted data


**_It is NOT production ready or secure, just an experiment!_**

So this script creates an in-memory SQLite db and inserts some data in an AES encrypted format and then tries to find some columns based on pre-generated indexes.
The indexes are also encrypted and there are different approaches how to interact with them.

There are very good articles about this topic, I used some ideas for this example from here: https://paragonie.com/blog/2017/05/building-searchable-encrypted-databases-with-php-and-sql

Example run:

```bash
$ go run main.go
2021/10/17 00:16:30 ID: 1 Data: {John Carter 1982}
2021/10/17 00:16:30 ID: 1 Encrypted: ilkpMZXWM9HwIDM7RHQ/63pye2uuJisNV1v43zWaUPax7xe4F5w9dvoRI79KX+xkqT3ETRh4Iup4H6eeXk0vkZ8pkDkqIebxS17WJxuBrGBFdA== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f Born index bucket: 0e9e566acee2e4c98eddba664d1e503eb44a340c3ea0aeee92bbe7aa272b6b4b
2021/10/17 00:16:30 ID: 2 Data: {John Doe 1994}
2021/10/17 00:16:30 ID: 2 Encrypted: SWGvPbccjITgx10bkpzAN1Gdx746z2NOddEfDsPBDfeq7v+VhuBv1UWX1CV+G8DHzis2wtzpqJrgNxWwpjPt1S61i4NRgRSranV2WoghhQ== FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f Born index bucket: ab2d67f4f05e212c8de0b198078d6e86c4cf73febc9e6c792b37ab3d48975816
2021/10/17 00:16:30 ID: 3 Data: {John Wick 1999}
2021/10/17 00:16:30 ID: 3 Encrypted: gKQNdRRxkcDn07FPXpH8Kipt1RYPipdqILuq4NbGmSwPzUjO46f5vxh19wA98Gkf6ybqVb8Jw72t3WEE5Ih9uiojrapB2qb2VKaR2EHrmb8= FirstName idx: 48d4f543569a062455e6375ab9c0ff5f56494f6453d1b7573bbd38923d5a30a3 FirstName bloom: 48d4f543569a062455e6375ab9c0ff5f Born index bucket: ab2d67f4f05e212c8de0b198078d6e86c4cf73febc9e6c792b37ab3d48975816
2021/10/17 00:16:30 ID: 4 Data: {Johnatan Somebody 1993}
2021/10/17 00:16:30 ID: 4 Encrypted: SAEjUq5unb4tgFHZkmqn7bVdPIrneqssLSxizVtnhI45Z64lKPUm85dGf3E0pYcp63iO8XEWCgY6OvhKdC2dKnZZAE1DuTkp2VLj3NoVGhJlN1PIFB8WLQ== FirstName idx: f3415b01e81a1477a918458e072655d5381b623dd6361f3531ec0d84c58cd511 FirstName bloom: f3415b01e81a1477a918458e072655d5 Born index bucket: ab2d67f4f05e212c8de0b198078d6e86c4cf73febc9e6c792b37ab3d48975816
2021/10/17 00:16:30 ID: 5 Data: {Somebody Else 2001}
2021/10/17 00:16:30 ID: 5 Encrypted: CSyVgqPGETouLDs9yFMHWWUOxneh99qYYszZgWgpnUmW0Rpa/OzsPafroy1HX2TSCG8RUzkcm5yGU0QzCJqxEHuXY441vowIDR5HI18jCCfPxpuw FirstName idx: 8c9f971a8c5797dab1a47a6d9be056fc78e44415f9b0de7c54e6ea362afb7488 FirstName bloom: 8c9f971a8c5797dab1a47a6d9be056fc Born index bucket: 1aba233f5e7453c59ef7cc11bc9d6b599693dae9b9d7af6a2f709d0290d73bad
2021/10/17 00:16:30 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/17 00:16:30 Fetch some data via index:
2021/10/17 00:16:30 ID: 1 Encrypted: ilkpMZXWM9HwIDM7RHQ/63pye2uuJisNV1v43zWaUPax7xe4F5w9dvoRI79KX+xkqT3ETRh4Iup4H6eeXk0vkZ8pkDkqIebxS17WJxuBrGBFdA==
2021/10/17 00:16:30 ID: 1 Data: {John Carter 1982}
2021/10/17 00:16:30 ID: 2 Encrypted: SWGvPbccjITgx10bkpzAN1Gdx746z2NOddEfDsPBDfeq7v+VhuBv1UWX1CV+G8DHzis2wtzpqJrgNxWwpjPt1S61i4NRgRSranV2WoghhQ==
2021/10/17 00:16:30 ID: 2 Data: {John Doe 1994}
2021/10/17 00:16:30 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/17 00:16:30 Do some blind index search:
2021/10/17 00:16:30 Looking for John
2021/10/17 00:16:30 Result, ID: 1 Data: {"first_name":"John","last_name":"Carter","born":1982}
2021/10/17 00:16:30 Result, ID: 2 Data: {"first_name":"John","last_name":"Doe","born":1994}
2021/10/17 00:16:30 Result, ID: 3 Data: {"first_name":"John","last_name":"Wick","born":1999}
2021/10/17 00:16:30 Looking for Somebody
2021/10/17 00:16:30 Result, ID: 5 Data: {"first_name":"Somebody","last_name":"Else","born":2001}
2021/10/17 00:16:30 Looking for Nobody
2021/10/17 00:16:30 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/17 00:16:30 Do some bloom filter search:
2021/10/17 00:16:30 Looking for John
2021/10/17 00:16:30 Result MATCHED, ID: 1 Data: {"first_name":"John","last_name":"Carter","born":1982}
2021/10/17 00:16:30 Result MATCHED, ID: 2 Data: {"first_name":"John","last_name":"Doe","born":1994}
2021/10/17 00:16:30 Result MATCHED, ID: 3 Data: {"first_name":"John","last_name":"Wick","born":1999}
2021/10/17 00:16:30 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/17 00:16:30 Do some bloom filter search with secondary filtering:
2021/10/17 00:16:30 Looking for John Doe
2021/10/17 00:16:30 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter","born":1982}
2021/10/17 00:16:30 Result MATCHED, ID: 2 Data: {"first_name":"John","last_name":"Doe","born":1994}
2021/10/17 00:16:30 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick","born":1999}
2021/10/17 00:16:30 Looking for John Carpenter
2021/10/17 00:16:30 Result NOT matched, ID: 1 Data: {"first_name":"John","last_name":"Carter","born":1982}
2021/10/17 00:16:30 Result NOT matched, ID: 2 Data: {"first_name":"John","last_name":"Doe","born":1994}
2021/10/17 00:16:30 Result NOT matched, ID: 3 Data: {"first_name":"John","last_name":"Wick","born":1999}
2021/10/17 00:16:30 ---------------------------------------------------------------------------------------------------------------------------------------
2021/10/17 00:16:30 Do some bucket-like filtering for a between query:
2021/10/17 00:16:30 Looking for values between 1995 (1990), 2005 (2010)
2021/10/17 00:16:30 Added bucket for 1990 ab2d67f4f05e212c8de0b198078d6e86c4cf73febc9e6c792b37ab3d48975816
2021/10/17 00:16:30 Added bucket for 2000 1aba233f5e7453c59ef7cc11bc9d6b599693dae9b9d7af6a2f709d0290d73bad
2021/10/17 00:16:30 Added bucket for 2010 a9088e5effc0d226cce5fcefaa00c613e84e90206b1f628b5d288657e7a3f98a
2021/10/17 00:16:30 Result NOT matched, ID: 2 Data: {"first_name":"John","last_name":"Doe","born":1994}
2021/10/17 00:16:30 Result MATCHED, ID: 3 Data: {"first_name":"John","last_name":"Wick","born":1999}
2021/10/17 00:16:30 Result NOT matched, ID: 4 Data: {"first_name":"Johnatan","last_name":"Somebody","born":1993}
2021/10/17 00:16:30 Result MATCHED, ID: 5 Data: {"first_name":"Somebody","last_name":"Else","born":2001}
```

