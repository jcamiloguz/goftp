# Go File Transfer Protocol (GOFTP)

This project is a TCP server to transfer file using a custom Protocol

## Actions

| Action      | Description                                          | Parameters                | Sender        |
| ----------- | -----------                                          | ---------                 | ------        |
| REG         | Action to Register in a GOFTP server                 |                           | Client        |
| OUT         | Action to go out in a GOFTP server                   |                           | Client        |
| PUB         | Action to send/publish a file to a specific channel  | Channel, FileName, Size   | Client/Server |
| FILE        | Action to send file content                          | Content                   | Client/Server |
| SUB         | Action for subscribe to a specific channel           | Channel                   | Client        |
| UNSUB       | Action for unsubscribe to a specific channel         |                           | Client        |
| OK          | Success confirmation                                 |                           | Client/Server |
| ERR         | Error notification                                   |                           | Client/Server |

### References
