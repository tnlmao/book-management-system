# Book Management System

This project is a Book Management System built in Go using Gin, Kafka, PostgreSQL, and Redis. Swagger is used for API documentation. This README provides one comprehensive file with all instructions to set up your dependencies, build the app, generate Swagger docs, and run everything as systemd services on your EC2 instance.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation & Build](#installation--build)
3. [Generate Swagger Documentation](#generate-swagger-documentation)
4. [Systemd Service Setup](#systemd-service-setup)
    - [Zookeeper](#zookeeper)
    - [Kafka](#kafka)
    - [Application](#application)
5. [Running the Application](#running-the-application)
6. [Usage](#usage)

---

## Prerequisites

- **EC2 Instance (Ubuntu recommended)**
- **Go** (>= 1.23.2 recommended)
- **Java** (OpenJDK 11 or similar; ensure JAVA_HOME is set correctly)
- **Kafka & Zookeeper** (Download and extract Kafka to, e.g., `/home/ubuntu/kafka`)
- **PostgreSQL** (Installed via apt; comes with its own systemd service)
- **Redis** (Installed via apt; comes with its own systemd service)
- Basic familiarity with Linux command line, Git, and systemd.

---

## Installation & Build

1. **Clone the Repository**
   git clone https://github.com/yourusername/book-management-system.git
   cd book-management-system
2. **Install Go Dependencies**
   go mod tidy
3. **Build the Application**
   go build -o book-management-system

## Generate swagger documentation

1. **Install Swag CLI**
   go install github.com/swaggo/swag/cmd/swag@latest
   export PATH=$PATH:$(go env GOPATH)/bin
2. **Generate Swagger Docs**
   swag init
3. **View Swagger UI**
   http://13.201.102.102:8080/swagger/index.html

## Systemd Service Setup
### Zookeeper
1. **Create the Unit File**
    [Unit]
    Description=Apache Zookeeper Service
    After=network.target

    [Service]
    Type=simple
    User=ubuntu
    ExecStart=/bin/bash /home/ubuntu/kafka/bin/zookeeper-server-start.sh /home/ubuntu/kafka/config/zookeeper.properties
    ExecStop=/bin/bash /home/ubuntu/kafka/bin/zookeeper-server-stop.sh
    Restart=on-failure
    Environment=JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64

    [Install]
    WantedBy=multi-user.target
### kafka
2. **Create the Unit File**
    [Unit]
    Description=Apache Kafka Service
    After=zookeeper.service

    [Service]
    Type=simple
    User=ubuntu
    ExecStart=/bin/bash /home/ubuntu/kafka/bin/kafka-server-start.sh /home/ubuntu/kafka/config/server.properties
    ExecStop=/bin/bash /home/ubuntu/kafka/bin/kafka-server-stop.sh
    Restart=on-failure
    Environment=JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64

    [Install]
    WantedBy=multi-user.target

### Application (Bookservice)
3. **Create the Unit File**
    [Unit]
    Description=Book Management System Service
    After=network.target zookeeper.service kafka.service postgresql.service redis.service

    [Service]
    User=ubuntu
    WorkingDirectory=/home/ubuntu/book-management-system
    ExecStart=/home/ubuntu/book-management-system/book-management-system
    Restart=on-failure
    Environment="DB_HOST=localhost"
    Environment="DB_PORT=5432"
    Environment="DB_USER=postgres"
    Environment="DB_PASSWORD=**********"
    Environment="DB_NAME=booksdb"
    Environment="REDIS_HOST=localhost:6379"
    Environment="KAFKA_BROKER=localhost:9092"

    [Install]
    WantedBy=multi-user.target

4. **Reload and Start Services**

    sudo systemctl daemon-reload
    sudo systemctl enable zookeeper kafka bookservice
    sudo systemctl start zookeeper kafka bookservice

    sudo systemctl status zookeeper
    sudo systemctl status kafka
    sudo systemctl status bookservice

## Running the Application
    The Go application listens on port 8080.
    Access the API endpoints at http://13.201.102.102:8080/api/v1/books
    Swagger UI is available at http://13.201.102.102:8080/swagger/index.html

## Usage
### API Endpoints:

#### GET /api/v1/books – Retrieve all books.
    http://13.201.102.102:8080/api/v1/books

#### GET /api/v1/books/{id} – Retrieve a specific book.
    http://13.201.102.102:8080/api/v1/books/{id}
    
#### POST /api/v1/books – Create a new book.
    http://13.201.102.102:8080/api/v1/books/
    Request: 
    {
        "title": "The Go Programming Language",
        "author": "Alan A. A. Donovan",
        "year": 2015
    }

#### PUT /api/v1/books/{id} – Update an existing book.
    http://13.201.102.102:8080/api/v1/books/{id}
    Request: 
    {
        "title": "updated",
        "author": "updated",
        "year": 2015
    }
    
#### DELETE /api/v1/books/{id} – Delete a book.
    http://13.201.102.102:8080/api/v1/books/{id}
