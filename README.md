# Coupon Service Repository

## 📝 Table of Contents

1. [📚 Overview](#-overview)
2. [🛠️ Features](#-features)
3. [🚀 Getting Started](#-getting-started)
   - [📦 Using Docker](#-using-docker)
   - [🛠️ Installation](#-installation)
4. [🖥️ API Documentation](#-api-documentation)
   - [📌 Base Path](#-base-path)
   - [📜 Endpoints](#-endpoints)
5. [🗂️ Data Models](#-data-models)
6. [🐳 Docker](#-docker)
7. [📜 Swagger Documentation](#-swagger-documentation)
8. [🛠️ Configuration](#-configuration)
   - [📋 Environment Variables](#-environment-variables)
   - [🗂️ `.env` File](#-env-file)
   - [🛠️ Precedence](#-precedence)
9. [🚀 Deployment](#-deployment)
10. [📦 Releases](#-releases)
    - [🔄 Workflow Automation](#-workflow-automation)
    - [📌 How to Add a New Release](#-how-to-add-a-new-release)
11. [📄 License](#-license)
12. [📝 Contributing](#-contributing)
13. [📞 Contact](#-contact)

---

## 📚 Overview

This repository is a fork of the [go-code-review](https://github.com/SchwarzIT/go-code-review) repository, enhanced with corrections and improvements to provide a robust coupon management service.

## 🛠️ Features

- **Coupon Creation:** Generate new coupons with validation.
- **Persistent Coupons:** Coupons will be saved in `/app/data/coupons.data.json`in container to ensure persistence.
- **Coupon Application:** Apply coupons to shopping baskets with validation.
- **Bulk Retrieval:** Retrieve multiple coupons by their codes.
- **Docker Support:** Easy deployment using Docker containers.
- **Swagger Documentation:** Interactive API documentation for development.
- **CI/CD Integration:** Automated release management and deployment.

## 🚀 Getting Started

### 📦 Using Docker

#### 1. **Pull the Docker Image**

```bash
docker pull matheuspolitano/couponservice:latest
```

#### 2. **Run the Docker Container**

```bash
docker run -d --name couponservice -p 80:80 -v ${PWD}/data:/app/data  matheuspolitano/couponservice:latest
```

- **Port Configuration:** Replace `-p 80:80` with the desired host and container ports.
- **Volume Configuration:** Replace `${PWD}/data:/app/data ` with the path that you want store the coupons.
- **Environment Variables & Volumes:** Add any necessary environment variables [🛠️ configuration](#-configuration) or volume mounts as needed.

### 🛠️ Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/matheuspolitano/couponservice.git
   cd couponservice
   ```

2. **Set Up Environment Variables**

Creating a .env file in the project root is optional. Instead of using a `.env` file, you can set the necessary environment variables directly in your system. If both the environment variables and the `.env` file are present, the environment variables will take precedence.

For detailed instructions on how to set environment variables, please refer to the [🛠️ configuration](#-configuration) section

```dotenv
API_PORT=9090
API_ENV=production
API_TIME_ALIVE=1y
API_SHUTDOWN_TIMEOUT=60s
API_ALLOW_ORIGINS=https://example.com,https://api.example.com
```

3. **Build and Run**

   ```bash
   docker build -t couponservice .
   docker run -d --name couponservice -p 80:80 couponservice
   ```

### Using Golang without deployment

```bash
go run ./cmd/coupon_service
```

**Notes:**

- **Example Coupons:** Two examples coupons have been added to the `data` folder to help you get started.
- **Golang Dependency:** Those examples only works when running the service with Golang.
- **Docker Image Creation:** The `data` folder is completely ignored during the Docker image creation process and will create another new one, so if you build the application using Docker, the contents of the data folder will not be included.

## 🖥️ API Documentation

### 📌 Base Path

`/api`

### 📜 Endpoints

#### 1. Create Coupon

- **Endpoint:** `POST /create`
- **Description:** Creates a new coupon with validation.
- **Request Body:**

  ```json
  {
    "code": "COUPON123",
    "discount": 10,
    "min_basket_value": 100
  }
  ```

- **Responses:**

  - **201 Created**

    ```json
    {
      "status": "success",
      "message": "Coupon created successfully",
      "data": {
        "id": "2e55fbb6-a856-4d25-9432-ce8f745bb118",
        "code": "COUPON123",
        "discount": 10,
        "min_basket_value": 100
      }
    }
    ```

  - **400 Bad Request**

    ```json
    {
      "status": "error",
      "message": "Error creating coupon",
      "error": "Coupon code already in use"
    }
    ```

- **Error Conditions:**
  - `discount` ≤ 0 or `min_basket_value` < 0.
  - `discount` > `min_basket_value`.
  - Duplicate coupon code.

#### 2. Apply Coupon

- **Endpoint:** `POST /apply`
- **Description:** Applies a coupon to a basket.
- **Request Body:**

  ```json
  {
    "basket": { "Value": 100 },
    "code": "COUPON123"
  }
  ```

- **Responses:**

  - **200 OK**

    ```json
    {
      "status": "success",
      "message": "Coupon applied successfully",
      "data": {
        "Value": 100,
        "applied_discount": 10,
        "application_successful": true
      }
    }
    ```

  - **400 Bad Request**

    ```json
    {
      "status": "error",
      "message": "Error applying coupon",
      "error": "Basket value 80 did not reach the minimum 100"
    }
    ```

- **Error Conditions:**
  - Coupon does not exist.
  - `basket.value` ≤ 0.
  - `basket.value` < `min_basket_value` of the coupon.

#### 3. Get Coupons by Code

- **Endpoint:** `GET /coupons`
- **Description:** Retrieves multiple coupons by their codes.
- **Query Parameters:**

  - `codes` (comma-separated, e.g., `COUPON12,COUPON123`)

- **Responses:**

  - **200 OK**

    ```json
    {
      "status": "success",
      "message": "Found all coupons",
      "data": [
        {
          "id": "2e55fbb6-a856-4d25-9432-ce8f745bb118",
          "code": "COUPON123",
          "discount": 10,
          "min_basket_value": 100
        }
      ]
    }
    ```

  - **404 Not Found**

    ```json
    {
      "status": "error",
      "message": "Error finding coupons",
      "error": "One or more coupons not found: COUPON000"
    }
    ```

- **Error Conditions:**
  - One or more coupon codes do not exist.
  - Missing `codes` query parameter.

## 🗂️ Data Models

- **`service.Basket`**

  ```json
  {
    "Value": 100
  }
  ```

- **`memdb.Coupon`**

  ```json
  {
    "id": "UUID",
    "code": "XYZ",
    "discount": 10,
    "min_basket_value": 100
  }
  ```

## 🐳 Docker

### Building the Docker Image

```bash
docker build -t matheuspolitano/couponservice:latest .
```

### Running the Container

```bash
docker run -d --name couponservice -p 80:80 matheuspolitano/couponservice:latest
```

- **Port Configuration:** Adjust `-p 80:80` as needed.
- **Environment Variables:** Pass environment variables using `-e` flags.
- **Volume Mounts:** Use `-v` flags to mount necessary volumes.

## 📜 Swagger Documentation

During **development**, access the Swagger UI at `http://localhost:<API_PORT>/swagger/index.html` to explore and test the API endpoints interactively.

## 🛠️ Configuration

Configure the application using environment variables. Supports a `.env` file and default values.

### 📋 Environment Variables

| Variable               | Description                                             | Required | Default                         |
| ---------------------- | ------------------------------------------------------- | -------- | ------------------------------- |
| `API_PORT`             | API port                                                | Yes      | -                               |
| `API_ENV`              | Environment type (`development` or `production`)        | No       | `development`                   |
| `API_TIME_ALIVE`       | Max server keep-alive duration (e.g., `10d`, `1y`)      | No       | Runs without duration           |
| `API_SHUTDOWN_TIMEOUT` | Max wait time for server to finish requests on shutdown | No       | `30s`                           |
| `API_ALLOW_ORIGINS`    | Allowed CORS origins (comma-separated URLs)             | No       | Allows all origins (production) |

### 🗂️ `.env` File

Create a `.env` file in the project root to set default values:

```dotenv
API_PORT=9090
API_ENV=production
API_TIME_ALIVE=1y
API_SHUTDOWN_TIMEOUT=60s
API_ALLOW_ORIGINS=https://example.com,https://api.example.com
```

### 🛠️ Precedence

1. **System Environment Variables** override `.env` values.
2. **`.env` File** provides defaults.
3. **Struct Defaults** apply if neither is set.

**Example:**

- `.env` has `API_PORT=9090`
- System sets `API_PORT=7070`
- **Effective `API_PORT`:** `7070`

## 🚀 Deployment

1. **Trigger & Skip:** Runs on main push or manual trigger; skips if the commit is by GitHub Actions.
2. **Release File:** Moves `release.toml` to the `releases` folder with a timestamped name, then commits.
3. **Version & Tag:** Finds the latest tag, determines the new version (major/minor/patch) from commit message or input, and creates the new tag.
4. **Docker:** Logs in, builds, and pushes Docker images to Docker Hub for the new tag and latest.
5. **GitHub Release:** Creates a GitHub Release and attaches the timestamped `release.toml` as an asset.

## 📦 Releases

The `releases` folder stores version-specific metadata for project updates, including features, improvements, and bug fixes. Each release file follows the standard template located at `./templates/release.toml`.

### 🔄 Workflow Automation

The repository includes a CI/CD pipeline to validate and manage release files:

1. **Validation:**
   - On pull request open or update, the pipeline validates the `release.toml` file using `scripts/validate_release.go` to ensure it follows the required structure.
2. **Merge Handling:**
   - On pull request merge, the pipeline automatically moves the validated `release.toml` file into the `releases` folder, renaming it with a timestamp for unique identification.

### 📌 How to Add a New Release

1. **Use the Template:**
   - Create a new release file using the template at `./templates/release.toml`.
2. **Add to Pull Request:**
   - Ensure the file is included in your pull request and follows the required structure.
3. **Merge Pull Request:**
   - Upon merging, the CI/CD pipeline will handle validation and placement in the `releases` folder.

## 📝 Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss changes.

## 📞 Contact

For any inquiries or support, please contact [matheuspolitano1@gmail.com](mailto:matheuspolitano1@gmail.com).
