# first and foremost do this
- git clone this repo

# must haves to run the program
- PostgreSQL installed.
- Go installed.

# to have the gator cli
- after cloning
```
go build -o gator/gator.exe
```
- and then
``` 
go install
```
- this will avail the gator as a cli app in go bin directory

# setup configs and running the program
- make a .gatorconfig.json file containing the fields db_url and currentUserName in the home directory of ur system
- once that is done, init the psql and run the migrations and sqlc generate
- and use go run . <cmd name> <args> / gator <cmd name> <args> to run the program

# commands u can run
- register
- reset
- login
- users
- addFeed
- follow
- following
- browse
- agg