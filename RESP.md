# REdis Serialization Protocol (RESP)

### Success Reply
- start with `+`，end with `\r\n`

> +OK\r\n



### Error Reply

- start with `-`, end with `\r\n`

> -Error message\r\n



### Integer Transfer

- start with `:`, end with `\r\n`

> wana sent 123456
>
> - :123456\r\n



### Multi line string Transfer

- start with $, followed by the actual `number of bytes sent`, followed by `\r\n`, and ends with `\r\n` 

> wanna sent "hedon.com"
>
> - $9\r\nhedon.com\r\n
>
> wanna send ""
>
> - $0\r\n\r\n



### Array Transfer

- start with `*`，followed by `the number of members`

> wanna send "SET key value"
>
> - *3\r\n$3\r\nSET\r\n\$3\r\nkey\r\n\$5\r\nvalue\r\n