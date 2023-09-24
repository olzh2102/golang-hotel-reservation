# Hotel reservation backend

## Project outline
- users -> book room from an hotel
- admins -> going to check reservation/bookings
- authentication and authorization -> jwt tokens
- hotels -> crud api -> json
- rooms -> crud api -> json
- scripts -> db management -> seeding, migration

## Resources
### MongoDB driver
Documentation
```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber
Documentation
```
https://gofiber.io
```

Installing gofiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Installing mongodb as a Docker container
```
docker run --name mongodb -d mongo:latesst -p 27017:27017
```