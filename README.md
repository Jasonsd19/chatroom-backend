## Simple Chatroom Application in Go
This is a simple chatroom application written in Go, utilizing the default `net/http` package and websockets for real-time communication. The application allows users to join the chatroom, send messages, and engage in real-time conversations with other users. [Front-end code here](https://github.com/Jasonsd19/chatroom-frontend)

# Features
* Real-time chat functionality using websockets.
* Basic validation for username length and uniqueness to ensure a seamless user experience.

# What I've Learned and Skills Gained
During the development of this chatroom application, I learned various essential concepts and gained valuable skills in Go and websocket development:

* Websockets: I gained a deep understanding of websockets and how they enable real-time communication between clients and servers. Implementing websockets allowed me to build a dynamic and interactive chatroom that updates in real-time. As a consequence of working with websockets I also learned a great deal about setting up HTTP routes and handler functions in Go, and in general.

* Go Channels and Concurrency: Utilizing Go channels and concurrency, I managed to handle multiple client connections simultaneously. This improved the application's performance and responsiveness during chat sessions.

* Design Pattern: Implementing the observer design pattern for handling communication between clients and the chatroom provided a clean and decoupled approach. This design pattern allowed me to maintain a well-structured codebase and facilitate easy communication between components.

* Error Handling and Validation: I developed a strong understanding of handling errors and basic data validation in Go. Validating the username length and ensuring uniqueness helped to provide a smooth user experience and prevent potential issues.

* Real-time Web Application Development: Building this chatroom application provided me with hands-on experience in developing real-time web applications. This project expanded my knowledge of frontend-backend communication and the intricacies of real-time updates.
