# messenger

Simple messaging service. Created to try out basic examples of building microservices using golang.<br>
Consists of four services communicating through gRPC:
- frontend - accepting client WebSocket connections and serves as API Gateway for other services;
- user - managing user infos;
- auth - managing authorization tokens;
- chat - managing chats and messages.

CLI app is included to test it out.

#### Stack: Go, MongoDB, ScyllaDB, gRPC, Docker 

---
## How to run

Repository have `Taskfile.yml` included to run and build bins.<br> 
Commands are written for windows only, here are the list:

- gen - generate protobuf part;
- build - build Docker image for every service;
- run - start necessary docker images via docker-compose, including mongo and scylla database clusters;
- test - run server tests;
- build-cli - build cli application bin;
- run-cli - run cli application bin.

---
## CLI Commands
List of available cli commands

- reg user_name password - register new user;
- login user_name password - log in user;
- logout - log out;
- init - required command after logging in, printing existing user chats;
- cur - get current user info;
- user user_id - get user info by id;
- users - get user info of all users;
- chats - get chats for currently logged in user;
- cchat user_id - create chat with specified user;
- msg chat_id content - send message to chat;
- msgs chat_id - get messages from chat.
