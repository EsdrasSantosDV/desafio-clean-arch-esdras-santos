# Clean Architecture: Listagem de Orders (REST, gRPC e GraphQL)

Desafio de Clean Architecture em Go com um unico use case de listagem de pedidos (`ListOrdersUseCase`) exposto simultaneamente via REST, gRPC e GraphQL.

## Como executar

Execute apenas:

```bash
docker compose up
```

O Docker Compose sobe MySQL, RabbitMQ e a aplicacao Go. A aplicacao aguarda o banco ficar pronto, aplica automaticamente a migration `migrations/001_create_orders.sql` e inicia os servidores.

## Portas

- REST/Web: `http://localhost:8000`
- gRPC: `localhost:50051`
- GraphQL: `http://localhost:8080/query`
- GraphQL Playground: `http://localhost:8080`
- RabbitMQ Management: `http://localhost:15672` (`guest` / `guest`)
- MySQL: `localhost:3306`

## Endpoints

### REST

- Criar order: `POST /order`
- Listar orders: `GET /order`

### gRPC

Service: `OrderService`

- `CreateOrder`
- `ListOrders`

Proto: `internal/infra/grpc/protofiles/order.proto`

### GraphQL

- Mutation: `createOrder`
- Query: `ListOrders`

## Requests prontas

Use o arquivo `api.http` na raiz do projeto para:

- Criar uma order via REST.
- Listar orders via REST.
- Criar uma order via GraphQL.
- Listar orders via GraphQL.

## Repositorio

https://github.com/EsdrasSantosDV/desafio-clean-arch-esdras-santos
