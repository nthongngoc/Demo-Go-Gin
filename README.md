# API project written by Golang

## Database diagram

<https://dbdiagram.io/d/5dc4e612edf08a25543dabec>

## Technology stack

- Web framework: Gin (<https://github.com/gin-gonic/gin>)
- Live reload: Air (<https://github.com/cosmtrek/air>)
- Swagger: Gin-swagger (<https://github.com/swaggo/gin-swagger>)
- Auth0: auth0-go (<https://github.com/auth0-community/auth0-go>)
- Environment parameter: GoDotEnv (<https://github.com/joho/godotenv>)
- Linter: Golint (<https://github.com/golang/lint>), GolangCI-lint (<https://github.com/golangci/golangci-lint>)
- JWT: jwt-go (<https://github.com/dgrijalva/jwt-go>), go-jwt-middleware (<https://github.com/auth0/go-jwt-middleware>)
- ODM framework: (<https://github.com/go-mgo/mgo/tree/v2>)
- Cloud Storage: (Firebase Cloud Storage: <https://firebase.google.com/docs/storage/admin/start>)

## Start Project

To use live-reloading in local environment:

```zsh
make run-local
```

To use live-reloading in development environment:

```zsh
make run-dev
```

To explicitly compile the code before you run the server with production environment:

```zsh
make run
```

To Compiling for every OS and Platform:

```zsh
make compile
```
