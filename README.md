# Advanced Programming

## Housing

> Description: Our website is designed for housing search

### Team Members:
- Alisher Tokesh | SE-2217
- Nurali Makhmutzhan | SE-2217

## Getting Started

To get started with this application, follow the steps below:

### Prerequisites

- Go (Golang) installed on your machine
- MongoDB 

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Tsarkashrk/housing.git
```

2. Install all necessary Go dependencies

3. Install MongoDB:
> Make sure you have MongoDB installed on your machine. If not, download and install it from the official MongoDB website.

4. Configure MongoDB connection:
> Update the database connection details in the main.go file:
```go
collection := client.Database("housing").Collection("users")
```

5. Run the application:
```bash
go run main.go
```

6. Open your web browser and go to http://localhost:8080 to access the application.

## Tools Used
> - Go (Golang)
> - Bcrypt
> - MongoDB (Database)

## API Endpoints
- GET /users: Get all users
- POST /register: Register a new user
- DELETE /users/delete/{id}: Delete user by ID
- UPDATE /users/update: Updating users list
