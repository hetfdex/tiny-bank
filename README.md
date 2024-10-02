The Makefile contains the following commands:
- start: Starts the application.
- start-docker: Starts the application within a Docker container.
- tests: Runs all the tests in the project (unit and integration).
- coverage: Runs the tests and generates a coverage report.

The API allows for:
- User creation and decactivation
- Account creation (multiple per user)
- Account deposit
- Account withdrawl
- Account transfer (includind between account of the same user)
- Account balance
- Account hisotry

Check the provided swagger file for more details on the API.

The internal directory contains the following subdirectories:
- handler: Defines the API endpoints and handles HTTP requests. It uses the Gin framework to route requests to the appropriate handlers.
- service: Contains the business logic of the application. It interacts with the repository layer to perform operations and return results.
- repository: Provides an abstraction for data storage. It defines interfaces and implementations for interacting with user, account, and transaction data.
- domain: Defines the core entities of the application, such as User, Account, and Transaction.

Assumptions:
- Built as a monolith service. User and account would be separate in a microservices approach.
- An assortement of tests to provide examples but lacking more.
- Transactions within Account model. Should likely be a different "table/repo".
- HTTP errors are incorrect for a lot of cases.
- Missing basic model props such as "currency" or "updated_at".
- Transaction model is not scalable.
- "Database" does not folow ACID principles.
- Missing API basics like "Get users".
- Validation of req models is basic.
- Some control of who can access account but also simplistic.
