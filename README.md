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