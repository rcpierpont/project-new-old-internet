# project-new-old-internet
Using what I've learned from bootdev to recreate something deployable and full-stack that can be used by close friends and family.

# Basic architecture
- Postgresql db for users, posts, comments
- Web App clients authenticate users and publish their posts so server can consume and add to db
- RabbitMQ post queues (users subscribe to other users, server subscribes to posts published by users so it can write to db)
- S3 bucket for storing image and video files

# Detailed breakdown of TODOs
- configure database schema
- implement authenticated database CRUD operations for users, posts, and comments
- implement basic ui for client-initiated CRUD operations
- implement image/video posts and connect to s3 storage
- set up rabbitMQ so created posts and comments are published to queue that server subscribes to
- set up rabbitMQ channel for server to publish posts, after being added to db, for users to subscribe to

# Tool dependencies
- goose and sqlc for golang-postgres layer
- docker
- rabbitmq
- s3

# Basic building approach

Start with basic db crud operations and then build a minimal front end to interact with them since it's needed to test image/video sharing. 

Implement authentication from frontend (user creation and login), then image/video sharing.

Finally, set up client and server publishing and subscribing, and then build a front end for client to view their subscription feed.