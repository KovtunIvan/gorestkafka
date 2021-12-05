# gorestkafka
Golang REST API to kafka (basic implemenation) on fasthttp and fasthttp/websocket


# Usage

<code> go run main.go </code>
<br> or <br>
<code> gin --appPort 8081 </code>
> The app will be listening 8081 (or 3000 for gin) port by default

# Handlers

1. Produce handler - /produce/{topicName} with `broker` header
2. Consume handler - /consume/{topicName} with headers below: <br>
    a. `broker` <br>
    b. `topic` <br>
    c. `consumer-group` <br>
    d. `timeout` <br>
    e. `commit-enable` (if `consumer-group` passed is `true` by default) <br>

> Consume handler has websocket implemenation
