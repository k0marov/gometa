# gometa 🚀

Full CRUD code generator adapted for [gobase](https://github.com/wajox/gobase) template. 

### Installation 

```bash
go install github.com/k0marov/gometa/cmd/cli@main
```

### Usage 

Create a schema file for a new entity, for example [client.schema.json](). 

Navigate to your project directory and type

```bash
gometa client.schema.json
```

After that, the CRUD code for this entity (Client) will be generated for all 3 layers: 
- Presentation layer (`controllers`) 
- Business logic layer (`services`) 
- Storage layer (`repository` using gorm) 

Error handling and logging will be automatically configured. 
New code will also be inserted into DI (`app` package), and Swagger annotations for swaggo will be added. 

### Features 

- Parsing schema files with support for many data types 
- Generating Go struct for the provided entity
- Generating handlers for 6 endpoints: create, update, delete, get one, get all
- Get all endpoint with automatically generated filters and pagination
- Generating service layer with logging and error handling 
- Generating repo layer with gorm bindings 
- Generating mappers for conversing between DTOs for different layers 
- Inserting into DI by editing .go files by manipulating the AST (Abstract Syntax Tree)
- Generating Swagger annotations 