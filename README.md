# Clean Architecture: Listagem de Orders (REST, gRPC e GraphQL)

## Objetivo

Neste desafio, voce deve implementar a funcionalidade de Listagem de Orders em sua aplicacao Clean Architecture. O objetivo principal e provar o desacoplamento da arquitetura: voce criara um unico Use Case (`ListOrders`) e o expora atraves de tres interfaces de comunicacao diferentes simultaneamente.

## Tecnologias e Padroes

- Linguagem: Go (Golang)
- Arquitetura: Clean Architecture
- Comunicacao: REST, gRPC e GraphQL
- Infraestrutura: Docker e Docker Compose

## Requisitos Tecnicos

### Use Case

Crie o caso de uso de listagem de pedidos (`ListOrdersUseCase`).

### Interfaces de Entrada

Disponibilize o acesso a esse Use Case atraves de:

- REST: Endpoint `GET /order`.
- gRPC: Service `ListOrders`.
- GraphQL: Query `ListOrders`.

### Banco de Dados

- Crie as migracoes necessarias para criar as tabelas do banco de dados.
- O banco deve ser provisionado via Docker.

## Requisitos de Dockerizacao (Automacao Total)

O avaliador nao deve executar nenhum comando manual alem do Docker Compose up.

### Container da Aplicacao

Voce deve criar um Dockerfile para a sua aplicacao Go.

### Orquestracao

O `docker-compose.yaml` deve subir o banco de dados e o container da aplicacao.

### Execucao Automatica

Ao rodar o comando:

```bash
docker compose up
```

- O banco de dados deve subir.
- As migracoes devem ser aplicadas automaticamente.
- A aplicacao deve iniciar e ficar disponivel nas portas configuradas.

Atencao: garanta que a aplicacao aguarde o banco estar pronto antes de tentar rodar as migracoes ou iniciar (handling de race condition na inicializacao).

## Arquivos Auxiliares

Crie um arquivo `api.http` na raiz contendo as requisicoes prontas para:

- Criar uma Order (para popular o banco e testar).
- Listar as Orders (para validar o desafio).

## Entregavel

Link do Repositorio: o link para o seu repositorio no GitHub.

## README

O arquivo deve conter:

- O comando unico de execucao (`docker compose up`).
- As portas em que cada servico (Web, gRPC, GraphQL) esta rodando.

## Regras de Entrega

### Repositorio Exclusivo (Muito Importante)

Este repositorio deve conter APENAS o codigo deste desafio.

Nao entregue um repositorio "monorepo" contendo pastas de outros cursos ou desafios anteriores. Isso bloqueia o processo de correcao automatica.

### Branch Principal

Todo o codigo deve estar na branch `main`.
