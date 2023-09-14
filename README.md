# COVID-19 Contact Tracing Backend Server

This is the backend server for processing, storing, and tracking COVID-19 risk contacts. It is responsible for managing data, communicating with clients via gRPC, and coordinating with various services, including a PostgreSQL database, MQTT server, and Redis cache. Docker is used for containerization, allowing for easy deployment and scaling.

## Directory Structure

- **`/server`:** Contains the backend server code. In `/server/docker` you can find Docker Compose configurations for both local testing and cloud deployment.

- **`/data-analysis`:** In the directory, you can perform data analysis on the collected test data. Before proceeding, make sure you have Anaconda installed. 

## Technology Stack

The backend server is built using the following technologies:

-   Golang: The primary programming language.
-   gRPC: For efficient communication between clients and the server.
-   PostgreSQL: The database for storing COVID-19 contact data.
-   MQTT: A message broker for handling real-time communication with the [mobile app](https://github.com/mikaellafs/tcc-contact-tracing-mobile).
-   Redis: Used for caching and optimizing data retrieval.
-   Docker: For containerization and deployment.

## Running the Server

To run the server locally for testing, follow these steps:

1.  Navigate to the `/server` directory.
    
2.  Start the necessary containers using Docker Compose for local testing:
    ```
    docker-compose -f docker/docker-compose-local.yaml up -d
    ```
    
4.  Run the backend server:
    ```
    go run cmd/main.go
    ```
    

## Data Analysis

In the `/data-analysis` directory, you can perform data analysis on the collected test data. Before proceeding, make sure you have Anaconda installed. Here are the steps:

1.  Install [Anaconda](https://docs.anaconda.com/free/anaconda/install/index.html).
    
2.  Create a Python environment and activate it:
    ```bash
    conda create -n contact-tracing python=3.8
    conda activate contact-tracing
    ```
    
3.  Install the required dependencies:
    ```
    ./data-analysis/setEnv.sh
    ```
    
4.  Run data analysis:
    ```
    python data-analysis/main.py
    ```
    

## Contributing

If you want to contribute to this project, please follow these steps:

1.  Fork the repository.
    
2.  Create a new branch for your feature or bug fix.
    
3.  Make your changes and commit them.
    
4.  Submit a pull request.
    

## Issues

If you encounter any issues or have suggestions for improvements, please **open an issue**.

## Mobile App Repository

The mobile app responsible for sending user contacts to the server and receiving notifications can be found in the [mobile app repository](https://github.com/mikaellafs/tcc-contact-tracing-mobile).
