# Badger DB Schema

/_ cSpell:disable _/

| Key                                                       | Value passed to/from dal | Description                    |
| --------------------------------------------------------- | ------------------------ | ------------------------------ |
| user:{id}#id                                              | string                   | username used to log in        |
| user:{id}#email                                           | string                   | email address                  |
| user:{id}#phash                                           | string                   | password hash                  |
| user:{id}#version                                         | int                      | version of user record         |
| user:{id}#event:{date}#id:{id}#activity:{id}              | []byte                   | activity name                  |
| user:{id}#event:{id}#exercise:{id}#index:{index}#instance | []byte                   | exercise instance              |
| user:{id}#activity:{id}#name                              | string                   | activity name                  |
| user:{id}#activity:{id}#exercise:{id}                     | string                   | exercise type id               |
| user:{id}#exercise:{id}#type                              | []byte                   | exercise type value            |
| tokenkey:{keyID}                                          | []byte                   | token key                      |
| pepperkey:{keyID}                                         | []byte                   | pepper key                     |
| session:{sessionid}#userID:{userID}#expires               | int64                    | session expiration time        |
| user:{id}#activity:{id}#program:{id}                      | []byte                   | training program               |
| user:{id}#program:{id}#instance:{id}                      | []byte                   | a program instance             |
| user:{id}#activity:{id}#activeprogram:{id}                | string                   | {programID: programInstanceID} |
| user:{id}#bio:{date}                                      | []byte                   | health daily stats             |

/_ cSpell:enable _/

Event Keys:

- primary sort by date
- enable query by date range
- enable query by activity
- enable query by exercise type
