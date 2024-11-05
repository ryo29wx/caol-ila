# Run docker
docker run --name some-mongo -d mongo:8.0 -p 27017:27017
docker run -it --rm mongo:8.0 -p 27017:27017 mongosh --host some-mongo test
docker exec -it 24dde02f3913 mongosh


# Connect to mongodb by using MongoSH
/ mongosh -u bar -p bar
test> use admin

# Create admin user
admin> db.createUser({user: "admin", pwd:"password", roles: [{ role: "root", db:"admin" }]})

# Swich and create user_info db
admin> use user_info

# Create db manager
user_info> db.createUser({user: "user_info_owner", pwd: "password", roles: [{ role: "dbOwner", db: "user_info" }]})

# Create `users` collection 
user_info> db.createCollection('users')

# data input
user_info> db.users.insert( { userid:'550e8400-e29b-41d4-a716-446655440000', displayname:'ryo kiuchi', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'./uploads/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125});

user_info> db.users.insert( {userid:'550e8400-e29b-41d4-a716-446655440001', displayname:'sample02', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125});

user_info> db.users.insert( {userid:'550e8400-e29b-41d4-a716-446655440002', displayname:'sample03', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125});

user_info> db.users.insert( {userid:'550e8400-e29b-41d4-a716-446655440003', displayname:'sample04', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125});

user_info> db.users.insert( {userid:'550e8400-e29b-41d4-a716-446655440004', displayname:'sample05', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125});

user_info> db.users.insert( {userid:'550e8400-e29b-41d4-a716-446655440006', displayname:'sample06', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125}); 

user_info> db.users.insert({ userid:'550e8400-e29b-41d4-a716-446655440007', displayname:'sample07', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125}); 

user_info> db.users.insert( {userid:'550e8400-e29b-41d4-a716-446655440008', displayname:'sample08', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125}); 

user_info> db.users.insert({ userid:'550e8400-e29b-41d4-a716-446655440009', displayname:'sample09', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125});
```


<details><summary> exec log </summary>

```
$ kubectl exec -it dev-mongo-statefulset-0 bash
/ mongosh -u bar -p bar
------
   The server generated these startup warnings when booting
   2024-02-20T10:03:47.906+00:00: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine. See http://dochub.mongodb.org/core/prodnotes-filesystem
   2024-02-20T10:03:53.407+00:00: vm.max_map_count is too low
------
test> use admin
switched to db admin
admin> db.createUser({user: "admin", pwd:"password", roles: [{ role: "root", db:"admin" }]})
{ ok: 1 }

admin> use user_info
switched to db user_info

user_info> db.createUser({user: "user_info_owner", pwd: "password", roles: [{ role: "dbOwner", db: "user_info" }]})

user_info> db.createCollection('users')
{ ok: 1 }

user_info> show collections
users

user_info> db.users.insert( { userid:'550e8400-e29b-41d4-a716-446655440000', displayname:'ryo kiuchi', sex:1, age:29, title:'swe', company:'rk', likes:['550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440099'], dislikes:['550e8400-e29b-41d4-a716-446655440080','550e8400-e29b-41d4-a716-446655440081'], blocks:['550e8400-e29b-41d4-a716-446655440070'], mainimage:'http://sample.com/image01', imagepath:['http://sample.com/image02', 'http://sample.com/image03'], registday:1645243125, lastlogin:1645243125} );
DeprecationWarning: Collection.insert() is deprecated. Use insertOne, insertMany, or bulkWrite.
{
  acknowledged: true,
  insertedIds: { '0': ObjectId('65d48b21448a95ef428cfe41') }
}
```

</details>