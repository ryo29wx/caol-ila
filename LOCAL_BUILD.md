# Local build

Use macbook.

```
$ colima start
-- uam
$ cd uam
$ docker build -t caolila-auth .
$ docker run -dit -p 8080:8080 caolila-auth

-- mongo
$ docker run --name some-mongo -d mongo:8.0
$ docker exec -it 24dde02f3913 mongosh 

-- admin
$ cd admin
$ go run admin.go

-- front
$ cd front-service
$ npm run dev
```

