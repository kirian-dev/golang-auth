# Go auth
This is a simple JWT (JSON Web Token) token package written in Go (Golang). It provides functionality to generate and parse JWT tokens, handle user registration and login, and store JWT tokens in cookies for authentication purposes.

### Features
User registration with password hashing and storage in-memory (non-persistent).
User login with password verification and token generation using JWT.
JWT token-based authentication via cookies for subsequent requests.
Lightweight, easy-to-understand code for learning purposes.

### Installation
1. Make sure you have Go (Golang) installed on your system.

2. Clone this repository to your local machine:
   `git clone https://github.com/yourusername/go-auth-test.git`

3. Navigate to the project directory:
   `cd golang-auth`

4. Install the necessary dependencies (if any) using:
   `go mod tidy`

5. To run the test project, use the following command:
   `go run main.go`
