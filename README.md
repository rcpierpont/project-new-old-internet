# project-new-old-internet

WORK IN PROGRESS

Using what I've learned from bootdev, I'm creating a full-stack social media app with basic post, comment, image, and video features. It will be deployable to docker and will require an aws account, but anyone with aws and docker accounts can start up their own environment and host it where they want, inviting only the people they want.

# Basic architecture
- Postgresql db for users, posts, comments, and file urls
- js front end (more design decisions tbd)
- Web App clients authenticate users and publish their posts so server can consume and add to db
- RabbitMQ post queues (users subscribe to other users, server subscribes to posts published by users so it can write to db)
- S3 bucket for storing image and video files

# Highest priority TODOs
- implement authenticated database CRUD operations for users, posts, and comments
- implement admin portal and self-service password changes
- implement ui for client-initiated CRUD operations (creating posts and comments) that leverages authn/authz where needed
- implement ability to attach image/video to posts by uploading to s3 storage
- set up rabbitMQ so created posts and comments are published to queue that server subscribes to
- set up rabbitMQ channel for server to publish posts, after being added to db, for users to subscribe to
- implement ability for users to follow other users and a ui page for posts of followed users

# Long-term TODOs
- implement ability to attach images/videos from other sites via links instead of requiring local upload
- implement likes and dislikes (maybe other reactions?)

# Tool dependencies
- goose and sqlc for golang-postgres layer
- docker
- rabbitmq
- s3

# Basic building approach

Start with basic db crud operations and then build a minimal front end to interact with them since it's needed to test image/video sharing. 

Implement authentication from frontend (user creation and login), then image/video sharing.

Finally, set up client and server publishing and subscribing, and then build a front end for client to view their subscription feed.

# Helpful commands

to refresh docker container with updated local code and start the container for testing:

```
cd <root of project or app directory>
docker compose up --build
``` 