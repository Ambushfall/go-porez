# Go Porez

This is a re-write of the qr porez app in GOLANG with more features (testing functionality and portability of Golang)

## Usage

- serve as HTTP Server ready to receive a POST Request
```bash
./backend-* -s
```

```bash
curl --request POST \
--url http://localhost/ \
--header 'authorization: Basic dkfhsdlepwmdseA==' \
--header 'cache-control: no-cache' \
--header 'content-type: application/json' \
--data '{
    "email":"email@email.com",
"password":"password"}'
```

- use as CLI with scheduler

```sh
./backend-* -e email@email.com -p password
```

- double click and input params manually