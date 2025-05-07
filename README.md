# Desafio Client x Server

## Client

Comando para executar aplicação
```cmd
U:\GO\clientserver\client>go run ./cmd  
// Se for pelo executável
client
```

Comando para gerar a build do client
```cmd
go build -o ./dist/client.exe ./cmd
```

## Server

Comando para executar o servidor
```cmd
go run ./cmd/main.go start
// Se for pelo executável
server start
```
>Note: Caso tenha instaldo o MSYS2 será necessário baixar
> set PATH=%PATH%;C:\msys64\mingw64\bin
> set CGO_ENABLED=1

Comando para montar as tabelas do BD
```cmd
go run cmd/main.go migrate up
// Se for pelo executável
server migrate up
server 
```

Comando para apagar as tabelas do BD
```cmd
go run cmd/main.go migrate down
// Se for pelo executável
server migrate down
```

Comando para listar cotações do BD
```cmd
go run cmd/main.go list
// Se for pelo executável
server list
```

Comando para gerar a build do servidor
```cmd
go build -o ./dist/server.exe ./cmd
```
