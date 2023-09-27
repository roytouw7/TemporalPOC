# Temporal POC

## Goals
Created a small POC for implementing Temporal.
My goals were:
- to keep low level abstractions(the way temporal stuff is executed) and higher level abstractions(business logic) separated as much as possible
- to allow for proper unit testing, with mocking certain parts
- to keep the Temporal surface as small as possible to prevent a vendor-lockin

This POC simulates the simplified version of upgrading an existing reservation, this upgrade consists of:
1. get the current room name for this specific reservation from a different microservice(mocked)
2. create an email message using this reservation and the fetched room name
3. calling the email microservice(mocked) to send this created e-mail to the user with the upgrade link

There are 2 non-deterministic pieces of logic involved with accompanied specific retry behavior, they are captures in Temporal activities:
1. getting the room name via http request
2. sending the e-mail via http request

## Running the POC

To run this project
1. open a terminal and navigate into the `TemporalServer` directory and run `docker compose up`
2. run the `worker/worker` *main* function
3. run the `start/main` *main* function