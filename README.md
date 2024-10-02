The Makefile contains the following commands:
- start: Builds and starts the application using Docker Compose.
- tests: Runs all the tests in the project.
- coverage: Runs the tests and generates a coverage report.

The internal directory contains the following subdirectories:
- handler: Defines the API endpoints and handles HTTP requests. It uses the Gin framework to route requests to the appropriate handlers.
- service: Contains the business logic of the application. It interacts with the repository layer to perform operations and return results.
- repository: Provides an abstraction for data storage. It defines interfaces and implementations for interacting with user, account, and transaction data.
- domain: Defines the core entities of the application, such as User, Account, and Transaction.
