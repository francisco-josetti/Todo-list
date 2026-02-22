# 📝 Todo App

Aplicação web de lista de tarefas desenvolvida em Go, com autenticação de usuários, sessões e banco de dados SQLite.

## 🚀 Funcionalidades

- ✅ Cadastro de usuário
- 🔐 Login e logout com sessão
- 📝 Criar tarefas
- ✔️ Marcar tarefa como concluída
- 🗑️ Remover tarefas concluídas
- 👤 Tarefas separadas por usuário
- 🔒 Senhas criptografadas com bcrypt
- 📦 Banco de dados SQLite

---

## 🛠️ Tecnologias Utilizadas

- Go (net/http)
- SQLite
- Gorilla Sessions
- bcrypt
- HTML Templates
- TailwindCSS (frontend)

---

## 📂 Estrutura do Projeto
.
├── main.go<br>
├── tasks.db<br>
├── templates/<br>
│ ├── base.html<br>
│ ├── index.html<br>
│ ├── login.html<br>
│ ├── auth.html<br>
│ └── create.html<br>
└── README.md


---

## ▶️ Como Executar

1. Clone o repositório:

```bash
git clone https://github.com/francisco-josetti/Todo-list.git
cd Todo-list
```
2. Instale as dependências
```bash
go mod tidy
```
4. Execute o projeto
```bash
go run main.go
```
6. Acesse o navegador
```browser
http://localhost:8080
```
## 🗄️ Banco de Dados

O banco SQLite é criado automaticamente ao iniciar a aplicação:

users → armazena usuários cadastrados

tasks → armazena tarefas vinculadas ao usuário

## 🔐 Segurança

Senhas armazenadas com hash usando bcrypt

Sessões protegidas com cookie store

Cada usuário só acessa suas próprias tarefas

## 👨‍💻 Autor

Desenvolvido por Francisco Josetti.

## imagem Ilustrativa

<img width="296" height="529" alt="swappy-20260222_043829" src="https://github.com/user-attachments/assets/4e61f3c6-79d2-4324-b051-a68cfaf98437" />




