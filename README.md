# Schwarz IT Code Review Repository
This project provides a simple coupon management service that allows clients to retrieve valid coupons and handle invalid coupon codes. The service exposes a REST API where you can POST coupon codes and get the results, which include valid coupons and any invalid codes. It also has an in-memory database to store created coupons. The service, once raised, will be up for up to one year.

## Project Structure

```
coupon-service/
├── cmd/
│   └── coupon_service/
│       └── main.go              # Main entry point for the coupon service
├── internal/
│   └── api/
│       ├── handler.go           # HTTP handler and API logic
│       └── service.go           # Service logic for coupon retrieval
├── entity/
│   └── coupon.go                # Defines the Coupon and CouponRequest structs
├── Dockerfile                   # Dockerfile to build the service
├── go.mod                       # Go modules file
├── go.sum                       # Go sum file for dependency management
└── README.md                    # Project documentation
```

## API Endpoints

### POST `/coupons`

This endpoint accepts a JSON body containing an array of coupon codes and returns a response with the valid coupons and any codes that were not found.

#### Request Body

```json
{
  "codes": ["valid1", "valid2", "invalid1"]
}
```

#### Response Body

```json
{
  "data": {
    "coupons": [
      {
        "code": "valid1",
        "discount": 10,
        "min_basket_value": 30
      },
      {
        "code": "valid2",
        "discount": 15,
        "min_basket_value": 50
      }
    ],
  }
}
```

### Error Handling

If the request body is malformed or invalid, the API will respond with a `400 Bad Request` status.

## Running the Project

To run the project locally, ensure that Go is installed on your machine and the dependencies are resolved.

1. Clone the repository:
   ```bash
   git clone https://github.com/SchwarzIT/go-code-review
   cd go-code-review/review
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the service:
   ```bash
   go run cmd/coupon_service/main.go
   ```

The service will be available at `http://localhost:8080`.

---

## Using Docker to Set Up the Environment

This project includes a `Dockerfile` to simplify building and running the service in a containerized environment. Below are the steps to build and run the service using Docker.

### Steps

1. **Build the Docker image**

   First, build the Docker image from the `Dockerfile`. You can do this by running the following command in the root of the project:

   ```bash
   docker build -t coupon-service .
   ```

   This will use the `golang:1.18-alpine` image to build the Go application, and then copy the compiled binary to a minimal Alpine-based container.

2. **Run the Docker container**

   Once the Docker image is built, you can run the container with the following command:

   ```bash
   docker run -p 8080:8080 -e API_PORT=8080 -e API_HOST=0.0.0.0 -e SKIP_CPU_CHECK=1 coupon-service
   ```

   This will start the application inside a container and map port `8080` on your local machine to port `8080` inside the container, so you can access the API at `http://localhost:8080`.

3. **Environment variables**
    
   - **`API_PORT`**  
     Specifies the port the application will listen on. By default, the application listens on port `8080`.  
   
   - **`API_HOST`**  
     Defines the host or IP address where the application will be accessible. Set it to `0.0.0.0` to ensure it listens on all interfaces.
   
   - **`SKIP_CPU_CHECK`**  
     If set to `1`, this variable skips any CPU-related checks during the application startup. This can be useful in specific environments where CPU checks are not necessary.  

---

## Running Tests

To run tests locally, use the following command:

```bash
go test ./...
```

This will execute all the unit tests in the project.

