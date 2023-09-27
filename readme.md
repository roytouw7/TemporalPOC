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

## Problem
Inside the `EmailWorkflow` lives the `UpgradeEmailWorkflowV3` function, this workflow function gets executed by
the *Temporal Worker* in a separate thread. This workflow has to be registered first, so it gets placed on the execution stack
for the workers, this happens in the `executeUpgradeEmailWorkflow` function.

Inside the `createChainOfCommand` function we create an orchestration of handlers composing the pieces of business logic wrapped in *Temporal Activites*
and normal functions, following the **Chain of Command design pattern**. The goal was to separate the orchestration and each piece of activity or business rule.

Passing in the factory used for this orchestration, or even just the collection of handlers, would allow for dependency injection and hereby to inject a mock factory or partial mocked collection of handlers.
This however is not possible due to all input of a workflow function has to be serializable to be placed
on the execution queue. Currently, the factory is directly used inside the workflow function preventing this. In the current solution every step can be unit tested properly, except for the orchestration.