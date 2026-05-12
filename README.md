🚀 App Backend Framework (Go)Este projeto é um framework backend modular, escalável e robusto escrito em Go (Golang). Ele foi projetado usando a arquitetura de "Package by Feature" (inspirada nos "Apps" do Django), permitindo que a aplicação cresça indefinidamente mantendo o código limpo e organizado.🛠️ Stack TecnológicoLinguagem: Go 1.22+Banco de Dados: PostgreSQL (via Docker)Roteador HTTP: Chi (go-chi/chi/v5) - Leve, 100% compatível com net/http.Database Tooling: SQLC (sqlc.dev) - Gera código Go typesafe a partir de SQL puro.Autenticação: JWT (golang-jwt) + Bcrypt (x/crypto/bcrypt)Validação de Dados: Validator (go-playground/validator/v10)Testes de Integração: Testcontainers (testcontainers-go) + Testify (stretchr/testify)Developer Experience (DX): Air (Hot Reload) + Makefile🏗️ Arquitetura e Padrões Aplicados1. Package by Feature (Módulos/Apps)Em vez de pastas clássicas como controllers/, models/ e routes/, o projeto é dividido em domínios de negócio (ex: users/, faq/). Cada módulo é auto-suficiente:schema.sql (Tabelas)queries.sql (Comandos DB)types.go (Structs de Request/Response)handler.go (Lógica HTTP)routes.go (Injeção no Roteador)2. O Padrão Flat REST & Validação Automática (internal/web)Criamos um utilitário central para padronizar respostas e automatizar o parse de dados:web.ReadJSON: Decodifica o JSON e roda o Validator. Se faltar um campo obrigatório (ex: validate:"required,email"), a API barra automaticamente retornando um erro 400 amigável.Tradução de Tags: O validador usa a tag label:"Nome" para retornar erros legíveis (ex: "O campo 'Nome' é obrigatório").Respostas Padronizadas: Uso de web.JSON e web.Error (retorna { "error": "mensagem" }), mantendo previsibilidade para o Frontend (React Native/Web).3. Injeção de Dependência (DI) & Configuração (internal/config)Variáveis de ambiente (.env) são carregadas uma única vez no pacote config. O main.go lê essas configurações e injeta as dependências (como a conexão do banco e o JWT Secret) nos construtores (NewApp). Os módulos não "caçam" variáveis globais.4. A Super Query do SQLC (O "Canivete Suíço")No módulo de faq, implementamos uma listagem avançada com SQLC sqlc.narg (Nullable Arguments).Isso permite que uma única query SQL lide com:Busca textual (ILIKE)Filtros booleanos opcionais (is_premium)Paginação (LIMIT e OFFSET)Ordenação dinâmica (ORDER BY CASE WHEN)Técnica aplicada no Go: Seleção dinâmica de campos (retornando apenas os fields que o Frontend pediu via map temporário).🧪 Estrutura de Testes (TDD & E2E)Aplicamos uma infraestrutura de testes nível sênior:Testcontainers (internal/testutil): Os testes não usam Mocks frágeis. O código sobe um container PostgreSQL real, aplica os schemas, roda o teste e apaga o container em ~2 segundos.Silêncio no Terminal: Configuramos o logger do Docker para ser descartado (io.Discard). O terminal foca apenas no que importa: PASS ou FAIL.Testes E2E Completos: Usamos httptest.NewRecorder e passamos a requisição pelo Roteador (r.ServeHTTP), garantindo que rotas, middlewares (como o AuthMiddleware) e o banco estão conectados.Air Test Watcher (.air.test.conf): TDD em tempo real. Salve o arquivo, os testes rodam.🚀 Como Rodar o ProjetoPré-requisitosGo instaladoDocker (ou Docker Desktop rodando no WSL2)Utilitário CLI: air (para hot reload) e sqlc (para gerar código de banco).1. Configurar o AmbienteCrie um arquivo .env na raiz do projeto:DB_USER=root
DB_PASSWORD=secret
DB_NAME=app_db
DB_HOST=localhost
DB_PORT=5432
SERVER_PORT=8080
JWT_SECRET=sua_chave_secreta_aqui
2. Subir o Banco de Dados (Desenvolvimento)make up
(Certifique-se de que o docker-compose.yml está configurado para o PostgreSQL)3. Rodar o Servidor (Com Hot-Reload)air
4. Modo TDD (Testes Contínuos)Em um terminal separado, inicie o observador de testes:# Roda todos os testes do projeto sempre que um arquivo for salvo
make test-watch

# Ou foque em um módulo específico para maior velocidade:
make test-watch pkg=faq
5. Trabalhando com Banco de Dados (SQLC)Sempre que você alterar um arquivo schema.sql ou queries.sql dentro de um app, regenere o código Go executando:sqlc generate
📞 Testando a APINa pasta de cada app (users/ e faq/), existe um arquivo request.http.Se você usa o VS Code com a extensão REST Client, basta abrir esses arquivos e clicar em "Send Request" para testar os endpoints diretamente do editor, simulando o registro, login e manipulação de entidades protegidas.



🚀 FAQ-Backend Framework GuideEste guia descreve os padrões arquiteturais, a estrutura de pastas e as melhores práticas para estender a nossa API Go modular.📂 1. Organização de Arquivos (Package by Feature)Cada novo domínio da aplicação deve ser tratado como um "App" isolado dentro de sua própria pasta.app/
├── internal/           # Código compartilhado e utilitários (Web, Testes, Config)
├── domain_app/         # Ex: users, faq, recipes
│   ├── db/             # Código gerado pelo SQLC (Package db)
│   ├── handler.go      # Lógica de negócio e Handlers HTTP
│   ├── routes.go       # Definição de rotas do módulo
│   ├── schema.sql      # DDL do banco de dados
│   ├── queries.sql     # SQL do SQLC
│   ├── db_test.go      # Helper para Testcontainers local
│   └── *_test.go       # Testes de Integração e E2E
└── main.go             # Orquestrador Principal
🏗️ 2. Criando um Novo App (Passo a Passo)Passo A: Banco de DadosCrie a pasta do app: mkdir -p recipes/db.Crie o schema.sql com as tabelas.Crie o queries.sql com os comandos.Adicione o novo módulo no sqlc.yaml e execute sqlc generate.Passo B: Handler e StructsDefina a struct App e as structs de entrada com tags de validação:type RegisterRequest struct {
    Title string `json:"title" validate:"required,min=3" label:"Título da Receita"`
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := web.ReadJSON(r, &req); err != nil {
        web.Error(w, http.StatusBadRequest, err.Error())
        return
    }
    // Lógica...
    web.JSON(w, http.StatusCreated, result)
}
🛣️ 3. Rotas e SegurançaAs rotas são definidas no arquivo routes.go de cada app. Utilizamos a injeção do Middleware de Autenticação para proteger endpoints específicos.Padrão de Registro de Rotas:func (a *App) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
    r.Route("/api/domain", func(router chi.Router) {
        // Rota Pública
        router.Get("/", a.List)

        // Grupo Protegido
        router.Group(func(protected chi.Router) {
            protected.Use(authMiddleware)
            protected.Post("/", a.Create) // Requer Token JWT
        })
    })
}
No main.go, basta injetar o middleware global:app.RegisterRoutes(r, users.AuthMiddleware)🛡️ 4. Validação de EntradaO framework utiliza o pacote validator/v10 integrado ao web.ReadJSON.Tags Suportadas:required: O campo não pode ser vazio.email: Deve ser um e-mail válido.min=N / max=N: Tamanho mínimo ou máximo.numeric: Apenas números.label:"Nome Amigável": Define o nome do campo na mensagem de erro.Exemplo de Erro Retornado:{"error": "O campo 'Título da Receita' deve ter no mínimo 3 caracteres"}🧪 5. Testes (TDD Flow)Estrutura de TestesUnitários: Ficam no *_test.go e testam lógica pura.Integração/E2E: Utilizam o testutil.SetupPostgresContainer para subir um banco Docker real.Executando TestesTudo: make testWatch Mode (TDD): make test-watchFocar em um App: make test-watch pkg=usersExemplo de Teste E2E:func Test_E2E(t *testing.T) {
    ctx := context.Background()
    dbConn, cleanup := SetupTestContainer(ctx, t)
    defer cleanup()

    r := chi.NewRouter()
    app := NewApp(dbConn)
    app.RegisterRoutes(r, users.AuthMiddleware)

    // Simular chamada HTTP Real
    req := httptest.NewRequest(http.MethodGet, "/api/domain", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}
✉️ 6. Padronização de Respostas (internal/web)Sempre utilize o utilitário web para responder ao cliente:web.JSON(w, status, data): Resposta de sucesso com corpo.web.Error(w, status, message): Resposta de erro formatada como {"error": "message"}.web.NoContent(w): Retorna Status 204 (usado em Deletes).web.ReadJSON(r, &dest): Decodifica e valida o corpo da requisição automaticamente.