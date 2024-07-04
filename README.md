# gometa 🚀

Full CRUD code Generator adapted for [gobase](https://github.com/wajox/gobase) template. 

### Installation 

```bash
go install github.com/k0marov/gometa@main
```

### Usage 

Create a schema file for a new entity, for example [blog_post.schema.json](examples/blog_post.json). 

Navigate to your project directory and type:

```bash
gometa blog_post.schema.json
```

After that, the CRUD code for this entity (BlogPost) will be generated for all 3 layers: 
- Presentation layer (`controllers` using [gin](https://github.com/gin-gonic/gin)) 
- Business logic layer (`services`) 
- Storage layer (`repository` using [gorm](https://gorm.io)) 

Error handling and logging will be automatically configured. 
New code will also be inserted into DI (`app` package), and Swagger annotations for swaggo will be added. 

You can just commit and use this new entity *without a single any code modification*, and, when needed, modify this CRUD to suit your business logic. 
It is possible because the generated code does not look like generated, it is fully readable and maintainable. 

### Features 

- Parsing schema files with support for many data types 
- Generating Go struct for the provided entity
- Support for either UUID or Integer Autoincrement for ID
- Support for timestamp fields 
- Generating handlers for 6 endpoints: create, update, delete, get one, get all
- "Get All" endpoint with automatically generated **filters and pagination**
- Generating service layer with logging and error handling 
- Generating repo layer with gorm bindings 
- Generating mappers for conversing between DTOs for different layers 
- Inserting into DI by editing .go files by manipulating the AST (Abstract Syntax Tree)
- Generating Swagger annotations 

### Docs 

Supported Data Types for Schema: 
- string 
- int
- float
- boolean
- Unix time (for specifying that field should be a Unix time in the `*.schema.json` file, use special value `1694801985`)

Every schema should have an "id" field. 
It can be of two variants: `uuid` and `integer autoincrement`. 

To specify that this entity's ID should be an integer autoincrement, use any integer, for example `{"id": 42}`

To specify that this entity's ID should be a UUID you can use any string, 
but it will be more clear if it's a UUID string, for example `{"id": "0814a807-077c-464b-8b82-8e41e8b4c68c"}`