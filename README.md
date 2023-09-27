# Music Sharing Microservices 🎵

I built this project in my own and it is a demonstration of my personal knowledge and skills in using **golang**, **nodejs** and **microservices architecture**. The goal of this project is to create a music sharing platform that allows users to upload, listen, like, and organize their favorite songs and playlists. 🎧

## Project Overview 📋

The project consists of three microservices that communicate with each other using the **Synchronous Messaging** pattern. No message brokers are involved in this project. Each microservice has its own database and API endpoints. The microservices are:

- **user-microservice**: This microservice is built in golang, using the gin-gonic framework. It handles user-related operations such as authentication and profile management. It uses MySQL as its database. 🔐
- **music-microservice**: This microservice is also built in golang, using the gin-gonic framework. It handles music-related operations such as uploading, updating, reading, and liking/unliking songs. It uses MongoDB as its database. 🎶
- **playlist-microservice**: This microservice is built in nodejs (typescript), using the fastify framework. It handles playlist-related operations such as creating, updating, and liking/unliking playlists. It also allows users to interact with other users' playlists. It uses PostgreSQL as its database. 📂

## How to Run 🚀

To run this project, you need to have docker and docker-compose installed on your machine. Then, follow these steps:

- Clone this repository to your local machine. 🖥️
- Navigate to the root directory of the project. 🗂️
- Run `docker-compose up -d` to start all the microservices databases. 🐳
- Wait for a few minutes until all the containers are up and running. ⏳
- You can access the API of each microservice at the following URLs:
  - user-microservice: [http://localhost:8080]
  - music-microservice: [http://localhost:8081]
  - playlist-microservice: [http://localhost:3000]

## Future Work 💡

This project is still a work in progress and there are many features and improvements that can be added. Some of them are:

- Implementing a front-end user interface using React or Angular. 💻
- Adding more tests and code coverage reports. 🧪
- Implementing a CI/CD pipeline using GitHub Actions or Jenkins. 🛠️

## Feedback 🙏

I hope you find this project useful and interesting. If you have any feedback or suggestions please feel free to contribute by adding pull requests. Thank you for your time and attention. 😊
