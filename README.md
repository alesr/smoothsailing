# SmoothSailing

SmoothSailing is a backend application that utilizes the Go programming language and the Encore Framework to provide robust REST APIs for managing users. The application is backed by PostgreSQL as its database, ensuring high performance and reliability.

The API offers various CRUD endpoints for user management, including:

- POST /users for user registration
- GET /users for listing all users
- GET /users/:id for fetching a user by its ID
- GET /me for fetching the current user
- DELETE /users/:id for deleting a user
- PATCH /users/:id for updating a user

For added security, requests to the backend API must contain a Bearer token for authorization. This token is kept as a secret within the backend and is referred to as the API_ACCESS_TOKEN. It is provided to clients upon request by the administrator.

Authentication of users is implemented using the industry-standard JSON Web Token (JWT) mechanism.

In the future, the application aims to offer an additional feature which allows users to fetch weather forecast based date and on geographic coordinates (latitude, longitude) or city name (to be defined).
