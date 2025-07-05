# Golang API REST

API REST robusta em Go implementando Clean Architecture com observabilidade, validação, tratamento de erros padronizado e configuração centralizada.

## 🚀 Funcionalidades

### ✅ Implementadas
- **Clean Architecture** - Separação clara entre domínio, aplicação e infraestrutura
- **Configuração Centralizada** - Viper para gerenciamento de configurações
- **Contexto Global e Cancelamento** - Shutdown gracioso com timeout
- **Middlewares Avançados**:
  - Request ID para rastreamento
  - Logging estruturado com contexto
  - Tratamento de erros padronizado
  - CORS dinâmico
  - Recovery de panics
- **Validação de Payload** - go-playground/validator com mensagens customizadas
- **Tratamento de Erros Padronizado** - AppError com códigos HTTP e detalhes
- **Observabilidade**:
  - Prometheus metrics em `/metrics`
  - Health checks em `/health/live`, `/health/ready`, `/health/detailed`
  - Logging estruturado com Logrus
- **Documentação Swagger** - Auto-gerada em `/docs/index.html`
- **Automação** - Makefile com comandos para desenvolvimento

### 🔄 Em Desenvolvimento
- Migrations versionadas com golang-migrate
- Tracing com OpenTelemetry e Jaeger
- Testes automatizados (unit, integration, e2e)

## 📋 Pré-requisitos

- Go 1.24.4+
- PostgreSQL 12+
- Docker (opcional)

## 🛠️ Instalação

1. **Clone o repositório**
```bash
git clone https://github.com/edumes/golang-api-rest.git
cd golang-api-rest
```

2. **Instale as dependências**
```bash
make deps
make tidy
```

3. **Configure o ambiente**
```bash
cp env.example .env
# Edite o arquivo .env com suas configurações
```

4. **Execute as migrations**
```bash
make migrate-up
```

5. **Execute o projeto**
```bash
# Desenvolvimento
make dev

# Produção
make run
```

## 📊 Observabilidade

### Health Checks
- `GET /health/live` - Verifica se a aplicação está viva
- `GET /health/ready` - Verifica se está pronta para receber requests
- `GET /health/detailed` - Informações detalhadas de saúde

### Métricas Prometheus
- `GET /metrics` - Métricas do Prometheus
- HTTP requests total, duração, em andamento
- Database connections e query duration
- Business operations

### Logging
- Logging estruturado com Logrus
- Request ID para rastreamento
- Contexto de usuário quando autenticado
- Níveis configuráveis (debug, info, warn, error)

## 🔐 Autenticação

A API usa JWT para autenticação. Inclua o token no header:

```bash
Authorization: Bearer <your-jwt-token>
```

## 📚 Documentação

### Swagger UI
Acesse a documentação interativa em:
```
http://localhost:8080/docs/index.html
```

### Endpoints Principais

#### Autenticação
- `POST /api/v1/auth/login` - Login de usuário
- `POST /api/v1/auth/register` - Registro de usuário

#### Usuários
- `GET /api/v1/users` - Listar usuários
- `POST /api/v1/users` - Criar usuário
- `GET /api/v1/users/:id` - Buscar usuário
- `PUT /api/v1/users/:id` - Atualizar usuário
- `DELETE /api/v1/users/:id` - Deletar usuário

#### Produtos
- `GET /api/v1/products` - Listar produtos
- `POST /api/v1/products` - Criar produto
- `GET /api/v1/products/:id` - Buscar produto
- `PUT /api/v1/products/:id` - Atualizar produto
- `DELETE /api/v1/products/:id` - Deletar produto

#### Projetos
- `GET /api/v1/projects` - Listar projetos
- `POST /api/v1/projects` - Criar projeto
- `GET /api/v1/projects/:id` - Buscar projeto
- `PUT /api/v1/projects/:id` - Atualizar projeto
- `DELETE /api/v1/projects/:id` - Deletar projeto

## 🧪 Testes

```bash
# Executar todos os testes
make test

# Executar com coverage
make coverage

# Executar linter
make lint

# Executar security check
make security-check

# Pipeline completo
make ci
```

## 🐳 Docker

```bash
# Build da imagem
make docker-build

# Executar container
make docker-run
```

## 🔄 Comandos Úteis

```bash
# Desenvolvimento
make dev              # Executar em modo debug
make fmt              # Formatar código
make vet              # Verificar código
make tidy             # Organizar dependências

# Database
make migrate-up       # Executar migrations
make migrate-down     # Reverter migrations
make seed             # Popular banco com dados

# Documentação
make swagger          # Gerar documentação Swagger

# Ferramentas
make install-tools    # Instalar ferramentas de desenvolvimento
make help             # Ver todos os comandos
```

## 🚀 Deploy

### Local
```bash
make build
./bin/golang-api-rest
```

### Docker
```bash
make docker-build
docker run -p 8080:8080 --env-file .env golang-api-rest:latest
```

### Docker Compose
```bash
docker-compose up -d
```

## 📈 Monitoramento

### Prometheus
Configure o Prometheus para coletar métricas de:
```
http://localhost:8080/metrics
```

### Grafana
Importe dashboards para visualizar:
- HTTP requests
- Database performance
- Business metrics