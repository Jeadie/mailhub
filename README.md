# mailhub


### Testing POST
```shell
curl --header "Content-Type: application/json" \
     --request POST \
      --data '{
        "Phone": "61412345678",
        "Content": "Hello World",
        "Date": "1970-01-01 12:00:01"
        }' \
  http://0.0.0.0:8080/sms/phone_number
```