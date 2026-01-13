# Meta (work in progress)

Meta is an application that extends the capabilities of Large Language Models (LLMs) by providing  them with a "meta"
layer of functionality. It includes a knowledge base for Retrieval-Augmented Generation (RAG), a recipe system for
guiding LLMs on specific tasks, and a secure way to manage these resources through a REST API, a minimal web interface
and a command-line
interface (CLI).

This application is not intended for enterprise use. It does not contain a user management system. The only barrier
between users is a long-lasting JWT token that should be considered an API key.

## Features

*   **Knowledge Base (RAG):** Ingest documents into a knowledge base and perform semantic searches to retrieve relevant
    information. This allows LLMs to access and reason about your private data.
*   **Recipes:** Create and manage "recipes," which are small, structured manuals that guide an LLM on how to perform
    specific tasks. This enables you to codify and reuse complex prompts and workflows.
*   **REST API:** A comprehensive REST API for managing the knowledge base and recipes.
*   **Model Context Protocol (MCP):** A protocol that allows LLMs to access the knowledge base and recipes.
*   **Multi-tenancy:** User and tenant isolation is enforced through JWT-based authentication.
*   **Web Interface:** A minimal web interface for managing the knowledge base and recipes.
*   **Command-Line Interface (CLI):** A CLI for the sole purpose of generating JWT tokens (API keys).

## Getting Started

### Prerequisites

* Docker (for containerized deployment)
* OpenSSL or any other tool of choice to generate the JWT RSA keys

### Setup

**Generate RSA keys for JWT:**

**IMPORTANT:** this repository contains a **test keypair. DO NOT USE FOR PRODUCTION!**
```bash
openssl genpkey -algorithm RSA -out etc/keys/private.pem -pkeyopt rsa_keygen_bits:4096
openssl rsa -pubout -in etc/keys/private.pem -out etc/keys/public.pem
```


### Running the Application

The recommended way to run the application is by using the provided Docker Compose file. This will set up and run all
the necessary services, including the main application, the database, and the LLM service (used exclusively to extract
text embeddings).

**Note:** beyond the Docker-compose, you will also need the `pg_init.sql` file.`

**Using Docker Compose:**

To start the application, run the following command in your terminal:

```bash
docker-compose up -d
```

This will start all the services in detached mode. The main application will be available at http://localhost:8080.

To stop the application, use the following command:
```bash
docker-compose down
```

### Exposed services

| Service  | Port  | Description       |
|----------|-------|-------------------|
| meta     | 8080  | The main application |
| postgres | 5432  | The database      |
| ollama   | 11435 | The Ollama service |


### Command-Line Interface (CLI)

The CLI's sole purpose is to generate tokens/api-keys. You can use the same docker container to run the CLI.

```shell
docker run -v ./etc:/usr/local/meta/etc ghcr.io/theirish81/meta:latest key --email foo@bar.com --permissions read
```

Options are:

*   `--email`: The email address of the user.
*   `--permissions`: The user's permissions (`read` or `write`).
*   `--subject` (optional): The subject of the JWT token. Consider this as a "user ID" and it's what identifies the
    ownership of the data. If you don't provide it, a random UUID will be generated. If you use it more than once to
    generate multiple tokens, it'll be as if you generated multiple tokens for the same user/organization.

The command will output a JWT token and the claims it contains.

### REST API

The REST API provides endpoints for managing the knowledge base and recipes. For a detailed description of the API,
please refer to the [OpenAPI specification](https://github.com/theirish81/meta/blob/main/spec/openapi.yaml).

**Authentication:**

All API requests must include a valid JWT token in the `Authorization` header. This token should be treated like an API
key.

```
Authorization: Bearer <your_jwt_token>
```

## Web Interface
You can access it at http://localhost:8080/web . It will require the same token mentioned above.

## Configuration

The application can be configured through a `.env` file or environment variables. Environment variables will always
supersede the values in the `.env` file.

* `DATABASE_URL`: the URL of the PostgreSQL database
* `KB_DISTANCE_THRESHOLD`: the vector distance beyond which a chunk of knowledge is considered irrelevant 
* `META_DISTANCE_THRESHOLD`: the vector distance beyond which a meta record is considered irrelevant
* `OLLAMA_BASE_URL`: the base URL of the Ollama service
* `EMBEDDING_MODEL`: the name of the text vectorization model to use. The model should produce vectors of exactly 2560
  dimensions (I recommend `qwen3-embedding:4b`)

## License

This project is licensed under the GNU Affero General Public License v3.0. See the `LICENSE` file for details.
