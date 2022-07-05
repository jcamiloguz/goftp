# Go File Transfer Protocol (GOFTP)

## Actions

| Action      | Description                                         | Parameters |
| ----------- | -----------                                         | ---------  |
| REG         | Action to send/publish a file to a specific channel | Channel, Content(File)    |
| UNREG         | Action to send/publish a file to a specific channel | Channel, Content(File)    |
| PUB         | Action to send/publish a file to a specific channel | Channel, Content(File)    |
| SUB         | Action for subscribe to a specific channel          | Channel    |
| UNSUB         | Action for unsubscribe to a specific channel          |      |

## Status Codes

| Code        | Description |
| ----------- | -----------             |
| 10          | File Sended to server   |
| 12          | File Sendend to server and recived by (x) clients  |
| 15          | File Sended to server but no client received it        |
| 20          | File not sendend        |
| 22          | File not sendend client error        |
| 25         | File not sendend server error        |
