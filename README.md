# Web Forum

Welcome to the Web Forum project! This is a web application that allows users to communicate with each other through posts and comments. It also includes features such as category association, post liking and disliking, and post filtering. The project is built using Golang and SQLite.

## Objectives

The main objectives of this project are:

- Enable users to engage in discussions by creating posts and commenting on them.
- Categorize posts to help users find relevant content.
- Allow users to express their opinion by liking or disliking posts and comments.
- Implement a filtering system to easily find posts based on categories, created posts, and liked posts.
- Utilize SQLite as the database management system for storing the application data.
- Implement user authentication and session management using cookies.
- Adhere to best coding practices and ensure robust error handling.

## Installation

To run this project, you need to have Docker installed on your machine. Follow the instructions below to get started:

1. Clone this repository to your local machine.
2. Navigate to the project directory.
3. Build the Docker image by running the following command:

```bash
make build
```

4. Start the Docker container by running the following command:

```bash
make run
```

5. Open your web browser and visit http://localhost:4000 to access the Web Forum.

If you prefer to run the application without Docker, you can use the following command:

```bash
go run ./cmd/web/*
```
## Usage

To use the Web Forum application, follow these steps:

1. Register a new user by providing your email, username, and password.
2. Log in to the application using your registered email and password.
3. Create a post by providing a title, content, and one or more categories.

Enjoy communicating and exchanging ideas with other users through the Web Forum!