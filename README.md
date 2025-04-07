# Sistema Na Mosca - Bolões de Futebol

Este é um sistema de bolões de futebol construído com Gin-Gonic, PostgreSQL e JWT, seguindo a arquitetura hexagonal.

## Requisitos

- Go 1.21 ou superior
- PostgreSQL
- Git

## Instalação

1. Clone o repositório:

```bash
git clone [url-do-repositorio]
cd na-mosca-server
```

2. Instale as dependências:

```bash
go mod download
```

3. Configure o banco de dados:

- Crie um banco de dados PostgreSQL chamado `na_mosca`
- Ajuste as credenciais no arquivo `.env` conforme necessário

4. Execute as migrações:

```bash
go run main.go
```

## Funcionalidades

- Criação e gerenciamento de bolões
- Participação em bolões existentes
- Palpites para jogos
- Ranking de participantes
- Notificações de resultados
- Sistema de pontuação automática

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

### Criar um bolão

```bash
curl -X POST http://localhost:8080/pools \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <seu-token-jwt>" \
  -d '{"name": "Brasileirão 2024", "description": "Bolão do campeonato brasileiro"}'
```

### Fazer um palpite

```bash
curl -X POST http://localhost:8080/guesses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <seu-token-jwt>" \
  -d '{"match_id": "123", "home_goals": 2, "away_goals": 1}'
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
