# sphelper

Library to interact with https://www.spieleplanet.eu/.

## Features
* Login and logout to https://www.spieleplanet.eu/
* Send get and post requests

## Usage

### Add sphelper to Your Project
```bash
go get github.com/yawn77/sphelper
```

### Set Environment
```bash
export SP_USERNAME=<username>
export SP_PASSWORD=<password>
```

### Use sphelper in Your Project
```go
// read credentials from env variables
creds, _ := sphelper.GetCredentials()
// create http client
client, _ := sphelper.GetClient()
_ = client.Login(creds)
_ = client.Get("https://www.spieleplanet.eu/")
_ = client.Post("https://www.spieleplanet.eu/...", url.Values{...})
_ = client.Logout()
```

## Development
For testing you have to provide a valid login for https://www.spieleplanet.eu/. If you are using vscode, you can do this by providing a `.env.test` file of the format
```
SP_USERNAME="<username>"
SP_PASSWORD="<password>"
```
In any case you can provide a username and a password as command-line arguments for the `go test` command
```bash
go test -v --race ./... -args -username="<username>" -password="<password>"
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
