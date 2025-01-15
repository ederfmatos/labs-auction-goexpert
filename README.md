# Auction API

Este projeto é uma **API REST para leilões**, desenvolvida em **Go** utilizando o framework **Gin** e o banco de dados **MongoDB**. A aplicação permite a criação de leilões, a realização de lances e a consulta de informações sobre leilões, incluindo o ganhador. Além disso, uma funcionalidade importante foi implementada: **o fechamento automático de leilões após um período de tempo definido**.

Foi adicionado um **teste de integração** que valida o fechamento automático de leilões. Este teste pode ser executado de forma simples, utilizando **Testcontainers** para subir o ambiente de teste e simular as requisições à API.

---

## 📚 **Funcionalidades Implementadas**

1. **Criação de leilões**: Endpoint para criar novos leilões com item e valor inicial.
2. **Listagem de leilões**: Endpoint para listar todos os leilões criados.
3. **Fechamento automático de leilões**: Após um tempo configurado, os leilões são automaticamente fechados.
4. **Teste de integração**: Teste automatizado que valida o fluxo de criação e fechamento de leilões.
5. **Arquivo `request.http`**: Contém duas requisições de exemplo para testar a API manualmente.

---

## 🧪 **Teste de Integração**

Foi adicionado um **teste de integração no arquivo `cmd/auction/main_test.go`** utilizando o pacote **Testcontainers** para validar o fechamento automático de leilões.

O teste:

1. Sobe um container com o banco de dados **MongoDB**.
2. Inicia o servidor Gin.
3. Faz uma requisição para **criar um leilão**.
4. Faz uma requisição para **listar os leilões** e obter o ID do leilão criado.
5. Espera um tempo definido e verifica se o leilão foi **automaticamente fechado** conforme a funcionalidade solicitada.

---

## 🛠️ **Como Executar o Projeto**

O projeto utiliza **Docker Compose** para facilitar a execução dos serviços.

### **Passo 1: Subir os containers com Docker Compose**

Execute o comando:

```bash
docker-compose up -d
```

Isso iniciará o container do **MongoDB** e outros serviços necessários.

---

### **Passo 2: Fazer as Requisições de Teste**

No arquivo `request.http`, existem duas requisições de exemplo:

1. **POST /auction** – Cria um novo leilão.
2. **GET /auction** – Lista todos os leilões.

Você pode usar um plugin no VS Code ou uma ferramenta como o **Postman** para executar essas requisições.

---

### **Passo 3: Executar os Testes Automatizados**

Para rodar o teste de integração diretamente pelo terminal, utilize o comando:

```bash
go test ./cmd/auction -v
```

Esse comando executará o teste automatizado que valida o fechamento automático de leilões. Durante o teste:

1. Um leilão será criado.
2. Após um tempo de espera, o leilão será marcado como **"Completed"** automaticamente.
3. O teste verificará se o comportamento está correto.

---

## 📄 **Exemplo de Requisição HTTP**

### **Criar um Leilão**

```http
POST http://localhost:8080/auction
Content-Type: application/json

{
  "product_name": "Product Test",
  "category": "Category Test",
  "description": "Description Test",
  "condition": 1
}
```

### **Listar Leilões**

```http
GET http://localhost:8080/auction
```

---
