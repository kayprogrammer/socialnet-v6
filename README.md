# SOCIALNET V6
A realtime social networking API built with Fiber & GORM ORM

![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/fiber.png?raw=true)


#### FIBER DOCS: [Documentation](https://docs.gofiber.io/)
#### GORM DOCS: [Documentation](https://gorm.io/docs/)
#### PG ADMIN: [Documentation](https://pgadmin.org) 


## How to run locally

* Download this repo or run: 
```bash
    $ git clone git@github.com:kayprogrammer/socialnet-v6.git
```

#### In the root directory:
- Install all dependencies
```bash
    $ go install github.com/cosmtrek/air@latest 
    $ go mod download
```
- Create an `.env` file and copy the contents from the `.env.example` to the file and set the respective values. A postgres database can be created with PG ADMIN or psql

- Run Locally
```bash
    $ air
```

- Run With Docker
```bash
    $ docker-compose up --build -d --remove-orphans
```
OR
```bash
    $ make build
```

- Test Coverage
```bash
    $ go test ./tests -v -count=1
```
OR
```bash
    $ make test
```

![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp1.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp2.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp3.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp4.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp5.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp6.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp7.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp8.png?raw=true)
![alt text](https://github.com/kayprogrammer/socialnet-v6/blob/main/display/disp9.png?raw=true)