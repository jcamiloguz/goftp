# Go File Transfer Protocol (GOFTP)

## Actions

| Action      | Description                                          | Parameters                | Sender |
| ----------- | -----------                                          | ---------                 | ------ |
| REG         | Action to Register in a GOFTP server                 |                           | Client |
| OUT         | Action to go out in a GOFTP server                   |                           | Client |
| PUB         | Action to send/publish a file to a specific channel  | Channel, Content(File)    | Client |
| INFO        | Action to especify the file info to a subscriber     | FileName, size            | Server |
| SUB         | Action for subscribe to a specific channel           | Channel                   | Client |
| UNSUB       | Action for unsubscribe to a specific channel         |                           | Client |
| OK          | Action for unsubscribe to a specific channel         |                           | Server |
| ERR         | Action for unsubscribe to a specific channel         |                           | Server |

### References
