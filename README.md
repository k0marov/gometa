# gometa
Gometa: Generate and not degenerate 😉

## Motivation 

Have you ever noticed that 90% of backend code is just cruds? 

You have entities that users can Create, Read, Update and then Delete. And that's all. Maybe you have some business logic in the middle, but it is just hidden behind this massive bloat of json parsing, db calls, error handling, logging and so on and so on. And the same crud logic is everywhere, for every model that you have in your application. And the only difference is the properties that these entities have and maybe some domain-specific business logic and calling external services.

This situtation is very sad. Writing so much degenerate repeating code is horrible not only for programmer's time, but also for the code quality. Because even when you conform to all modern standards of clean architecture, good separation of concerns, and DRY (Don't Repeat Yourself), you ironically have massive code repetition at the core of your project. 

Actually, it's code repetition in a more high-level essence. Normally you see repeating code where it is either fully identical, or has different variables or literals, and you DRY it up into a function or a class and everything becomes clean and good. But this CRUD gigapattern cannot be beaten by this method, because it is entities that are different and not variables, it also spans multiple layers and needs too much customization. Functions are too low-level for this scenario, so probably the thing that can help in this case needs to be like one level above your program.. And here metaprogramming comes. 

Gometa wants to overcome this problem by using a form of metaprogramming - code generation. In the past people referred to code generation as a bad practice, but nowadays it is not viewed as such: we are all now accustomed to such things as mock generation and protoc for gRPC.  

It will generate all of the CRUD code for a given entity and let you fill in the gaps for business logic layer. 

## Goals 

Gometa will aim to be highly customizable, supporting changes to templates to accompany specific project needs, selection of databases and Go web frameworks, lots of places where you may put custom domain logic. 

It will be as modular as possible, with separate generators for different application layers. 

The above should be achieved by having as simple codebase as possible. 

It is not a goal of Gometa to generate code that would be commited to a VCS or read by a human, because then it will save developer time, but would not solve the problem of overbloated systems. Ideally code should be regenerated every time at the precompile stage.   

Gometa will not parse go code (like a struct with tags for generation), since it is troublesome, not concise and has too many edgecases. Instead, gometa will use a well-known config format, such as YAML or JSON. 

Gometa should have a separated generator core so that it can later be used for other languages, frontend, etc. 

## Current State 

Currently, Gometa is just an idea. It looks like a hard task, and I will try to develop it incrementally. 

## Todo 

- [x] Project initialization
- [x] Parsing the entity scheme
- [x] Generating repository layer
- [x] Generating entity struct
- [x] Generating service layer 
- [x] Generating http handlers layer
- [x] Split layers into separate directories
- [x] DI file 
- [x] API Endpoints file
- [x] Generating base main file
- [ ] Logs
- [ ] Logging in templates
- [ ] Error handling 
- [ ] Acceptance test for generator
- [ ] Customization of templates
- [ ] CLI interface similar to protoc
- [ ] Swagger generation
