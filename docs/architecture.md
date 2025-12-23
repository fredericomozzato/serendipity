# Serendipity - Architecture Documentation

**Version:** 1.0
**Date:** 2025-12-20
**Author:** Frederico Mozzato
**Status:** Draft for Review

---

## 1. Introduction and Goals

### What is Serendipity?
> [!quote]
> The ability to find valuable or agreeable things not sought for.
> [Merriam-Webster's definition](https://www.merriam-webster.com/dictionary/serendipity)

Serendipity is a random music discovery app built without any recommendation
algorithm. The song selection is purely random and the users may expect
anything to pop up. It is an app for people that like discovering new artists
and songs and that enjoy listening to different and unexpected things.

The target user wants to explore songs out of the normal spectrum and is not
interested in the valuable but predictable recommendations used by regular
music streaming apps, that use listening patterns to find relevant songs. Our
users like to be surprised.

Serendipity allows any user to browse a player that generates a new random
song. The player will allow the users to filter based on genre, year and
country of origin, but that's about it. Everything else is random and subjected
to chance.

Users can access their listening history, like songs to add them to their
collection and create playlists to group songs together. The user will also be
able to download songs directly to their computers for sampling purposes.

The playlists can be shared through links so other users can hear
mixes/collections, adding a social aspect to the application.

---

## 2. Constraints

### Language
Serendipity will be implemented in Golang since this is a great language to
handle web services and allows the app to be fast with its asynchronous
capabilities. This is also a study project to help me, as a developer, learn
the language and build a portfolio with it.

### Architecture
The application will be divided in three main services:

1. The web server
2. The download service
3. The notification service

This choice is done to help separate the concerns, increase isolation and help
me learn the micro service paradigm with a real project.

The app will be backed by the [Discogs
API](https://www.discogs.com/developers). The music data will be all fetched
from this open and free database with millions of artists and songs.

One important consideration to make is that we'll be doing Server Side
Rendering in the web server so we can focus as much as possible in building the
backend (which is the main focus of this project) while not having to maintain
a second Javascript application.

This also means that the front end design and implementation will have 
massive help from AI tools, adding another layer for learning on prototyping 
and using AI coding agents to get things done. The initial idea is to use 
[Lovable](https://lovable.dev/) and [Claude Code](https://claude.com/product/claude-code) for this task.

## 3. System Context

### 3.1 Users
Serendipity will allow both anonymous and authenticated users to access the
service. No one will *need* to make an account to use the application but the
scope will be different for each type of user.

#### Anonymous users
Can access the application through the web. Will be able to use the player,
filter results and play tracks.

#### Authenticated users
Will access the app the same way as anonymous users and have all the same
capabilites plus being able to view and manage their listening history, liking
songs to add them to their collection, creating and sharing mixes and
downloading songs.

### 3.2 External Entities:
#### Discogs API
As mentioned, the main dependency is the Discogs API. This is the music
catalogue that Serendipity will drink from, fetching release data from it and
displaying it in the player.

The Discogs API returns results that have embedded Youtube video links. We'll
use these links to render the Youtube player so the users can listen to the
songs.

The API has a rate limit of 60 requests per minute for non-authorized
applications and 100 requests per minute for authorized applications. To
increase our rate limit we'll authorize our application by using the regular
auth flow to identify our application to the API server through a User Agent
and an authorization token. We'll have to keep our requests under 100 rpm to
avoid being throttled.

#### AWS
The application will be deployed to AWS in a EC2 instance. We'll also use more
services from the cloud provider.

##### S3
Will be used to store extracted audio files for user's download. The files will
be available for a period of 48 hours before automatic cleanup using S3
lifecycle policies.

##### SES
The download links will be sent to the users using the AWS SES. Users will
receive the email as soon as the audio file is uploaded to the S3 service. The
email will contain a link and a note about the expiration date. Users will be
able to directly download the mp3 from the link.


## 4. Solution Strategy

### 4.1 Technology Decisions
#### Golang
The main programming language chosen is Go. This has two reasons:

1. This is a project to help learn and create portfolio in Golang.
2. Go is a great language to build web servers for its concurrent capabilities.

#### Microservices architecture
The choice to split the application in services has to do with my need to learn
and implement this microservice patterns, since I'm a Ruby on Rails developer
and I'm very familiar with monolithic application.

#### PostgreSQL
As the industry standard database I chose to use it for it's JSON support
(since we'll be working extensively with external API data in this format).

#### RabbitMQ
The download and email services will work as background job workers. RabbitMQ
seems to be the industry standard for job processing with robust retry logic,
allowing us to ensure the downloads will be completed and delivered to users
even if there are errors in the way.

#### Server-Side Rendering
The server will do SSR so I don't have to maintain a second Javascript
application (as this is out of this project's scope).

## 5. Building Block View
### 5.1 Level 1: System Overview (Container Level)
![Containers diagram](./diagrams/serendipity_containers.jpg)

#### Web server
Used to handle HTTP requests and contains the main song discovery logic. We'll
do Server Side Rendering (SSR) to avoid creating and maintaining a second
application in Javascript, as this is out of scope for this project.

The web server will have read/write access to the main application's database
for storing and retrieving user data like listening history, the collection and
mixes. It'll also be responsible for querying the Discogs API to fetch releases
while ensuring we keep under the rate limits.

The server will implement a Token Bucket pattern to ensure we keep under the
100 RPM limit imposed by the Discogs API.

#### Downloads service
This service will be responsible for downloading and extracting audio from
videos on demand. Users can click on a *Download* button in the interface and
this will send a request to this service.

The service will use the `yt-dlp` package to download and convert the audio to
mp3. This file will then be uploaded to Amazon S3 so the user can dowload it.

Once the upload to S3 is complete this service will send messages to the web
server and the notifications service to indicate the success.

#### Notifications service
This service is responsible for sending emails to the users with their download
links. The email will contain the link to the S3 file, an expiration date and
the name of the song.

Users will be able to click and download the file they requested. This flow was
designed like this for the async nature of the download flow and to allow users
quick access to their links even without access to the app's web interface.

#### Database (PostgreSQL)
The database will be PostgreSQL 18. This was chosen for the good JSON support
(since we'll be working with API payloads) and for better text search support.

#### Message broker (RabbitMQ)
We'll use AMQP to send messages between our three services. The chosen broker
was RabbitMQ for its wide industry use and its robust capabilities for retries.
Since we'll use messages to orchestrate the download and notification services
we can't afford to lose information or jobs.

One important consideration about the broker design is that we want the
services to be database independent, so we'll pass all the data required for
completing a work unit within the messages. This way only the web server will
need access to the main database, while the downloads and notifications service
can do everything with the data received by the message.

#### Object storage system (Amazon S3)
S3 will be used to store the audio files extracted from the Youtube videos.
Users will receive the link to their downloads via email and the links will
point to the files in this storage.

> [!note]
> The S3 service will use lifecycle policies to clean up old download files.
> Users will have 7 days to download their files before the link is expired.

### 7.2 Deployment Strategy
**Needs more research**

The application will run in AWS in an EC2 instance.

- CI/CD pipeline: Github actions pipeline

### 7.3 Monitoring and Operations

- Logging: [What tool? What gets logged?]
- Metrics: [What do you track?]
- Alerts: [What triggers alerts?]
- Health checks: [What endpoints?]

---

## 8. Cross-Cutting Concepts
### 8.1 Security
**Authentication:**
The authentication will be done exclusively with oAuth2. This will ensure that
I learn how to implement this type of authentication while reducing the
security requirements for the app, since I won't need to store passwords. This
will also ensure users have a valid and correct email where we can send them
download links.
