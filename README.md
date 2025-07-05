# Golang API REST

API REST em Go, seguindo Clean/Hexagonal Architecture, pronto para produção e escalabilidade.

## Visão Geral
- Estrutura modular (api, application, domain, infrastructure)
- CRUD completo de usuários, produtos, projetos e itens de projeto
- Autenticação JWT
- **Logging abrangente com Logrus**
- Observabilidade (Prometheus, OpenTelemetry)
- Pronto para Docker, Kubernetes e CI/CD

## Pré-requisitos
- Go 1.21+
- Docker e Docker Compose
- PostgreSQL

## Instalação
```sh
git clone https://github.com/edumes/golang-api-rest.git
cd golang-api-rest
go mod tidy
```

## Como rodar
```sh
# Usando make
make run

# Ou diretamente
go run cmd/api/main.go

# Com Docker
docker-compose up
```

## Configuração
- Copie `.env.example` para `.env` e ajuste as variáveis.

## Logging

### Visão Geral
O projeto implementa logging abrangente usando **Logrus** em todas as camadas da aplicação:

- **Handlers**: Log de requisições, respostas e erros
- **Services**: Log de operações de negócio e validações
- **Repository**: Log de operações de banco de dados
- **Middleware**: Log de autenticação e performance
- **Main**: Log de inicialização e shutdown

### Níveis de Log
- **DEBUG**: Informações detalhadas para desenvolvimento
- **INFO**: Eventos normais da aplicação
- **WARN**: Situações que merecem atenção
- **ERROR**: Erros que não impedem a execução
- **FATAL**: Erros críticos que param a aplicação

### Estrutura dos Logs
Todos os logs incluem campos estruturados:
```json
{
  "level": "info",
  "msg": "User created successfully",
  "time": "2024-01-15T10:30:00Z",
  "user_id": "uuid",
  "email": "user@example.com",
  "method": "POST",
  "path": "/api/v1/users",
  "ip": "192.168.1.1",
  "latency": "150ms"
}
```

### Logs por Camada

#### Handlers (API Layer)
- Log de entrada de requisições
- Log de validação de dados
- Log de respostas de sucesso/erro
- Log de autenticação e autorização

#### Services (Application Layer)
- Log de operações de negócio
- Log de validações
- Log de transformações de dados
- Log de chamadas para repositories

#### Repository (Infrastructure Layer)
- Log de conexões com banco
- Log de queries executadas
- Log de filtros aplicados
- Log de resultados de queries

#### Middleware
- Log de performance (latência)
- Log de autenticação JWT
- Log de headers de requisição
- Log de recovery de panics

### Monitoramento de Performance
O middleware de logging captura automaticamente:
- **Latência** de cada requisição
- **Status codes** de resposta
- **User Agent** do cliente
- **IP** do cliente
- **Trace ID** (se fornecido)

## Comandos úteis
- Build: `make build` ou `go build -o golang-api-rest cmd/api/main.go`
- Migrations: `make migrate-up`
- Seeds: `make seeds-all` ou `make seeds-users`
- Swagger: `make swag`
- Testes: `make test`

## Documentação
- Swagger: `/swagger/index.html`

## Autenticação
- `POST /v1/auth/login` para obter JWT
- Use o token no header: `Authorization: Bearer <token>`

## Seeds

O projeto inclui um sistema de seeds para popular o banco de dados com dados iniciais.

### Executando Seeds

```bash
# Executar todos os seeds
make seeds-all

# Executar apenas seeds de usuários
make seeds-users

# Via linha de comando
go run cmd/seeds/main.go -type=users
```

## Documentação
- Swagger: `/swagger/index.html`