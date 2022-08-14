# Plow

Tools for the Ginkgo testing framework that simplifies microservices integration testing

### Inspired by Gonkey

Gonkey (github.com/lamoda/gonkey) is a testing toolkit that work with testing specifications written in yaml format.
Gonkey has many good ideas and paradigms, for example:
- response/request testing
- mocks for external services
- database fixtures
- database query checks

But the yaml format is not a very flexible format for writing extension.

Plow has adopted proven techiques from Gonkey codebase to the Ginkgo testing framework

### Codebase sharing

- some code migrated from github.com/lamoda/gonkey testing framework

### Features

- Request/response checks
- Data storages fixtures (PostgreSQL, MySQL, Aerospike, Redis)
- Database query checks
