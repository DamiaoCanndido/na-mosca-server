# Sistema de Autenticação com Gin-Gonic

Este é um sistema de autenticação básico construído com Gin-Gonic, PostgreSQL e JWT, seguindo a arquitetura hexagonal.

## Requisitos

- Go 1.21 ou superior
- PostgreSQL
- Git

## Instalação

1. Clone o repositório:

```bash
git clone [url-do-repositorio]
cd bolao
```

2. Instale as dependências:

```bash
go mod download
```

3. Configure o banco de dados:

- Crie um banco de dados PostgreSQL chamado `bolao`
- Ajuste as credenciais no arquivo `.env` conforme necessário

4. Execute as migrações:

```bash
go run main.go
```

## Uso

### Criar um usuário

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email": "usuario@exemplo.com", "password": "senha123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "usuario@exemplo.com", "password": "senha123"}'
```

### Rotas Protegidas

Para acessar rotas protegidas, inclua o token JWT no header:

```
Authorization: Bearer <seu-token-jwt>
```

## Estrutura do Projeto

```
.
├── internal/
│   ├── domain/           # Entidades e interfaces do domínio
│   ├── ports/            # Casos de uso e serviços
│   └── adapters/
│       ├── driven/       # Adaptadores de saída (banco de dados)
│       └── drivers/      # Adaptadores de entrada (HTTP)
├── main.go              # Ponto de entrada da aplicação
├── go.mod               # Dependências do projeto
└── .env                 # Configurações do ambiente
```
