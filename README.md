# presenceHoursLog

this is a warmup project for go rest api development and uses basic authentication for user actions

one can create users and monitor their working hours.
users can record their presence hours by starting or stopping their sitting sessions

## Setup requirements
You'll need to run `go build` to create your own .exe file. prior to that you'll need the following
## .env file
The directory requires a .env file 
- specifying the mongodb connection url as "DB_URI"
- specifying the server connection port as "PORT"
---

- the "/user" endpoints require level 1 access
- the "/admin" endpoints require level 2 access

supported handlers:

- GET request to server root
> returns a string to show that server is up and running

- GET request to user/admin/all
> returns a json of all users 

> the __admins__ collection from the __presenceLog__ database contains all valid username password combinations

- POST request to user/admin/users
> do include a json body as the example suggests:

> `{
  "name" :"Hesam",
  "username":"DoubleH",
  "password": "Hesam12345"
  "Access": 3
  "age": 24
}`

- GET request to /admin/user/myinfo
> returns the complete user info

- POST request to user/start/
> starts a session for the user

- POST request to user/stop/
> stops the session for the user


