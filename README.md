# Overview

This repository was created from the repository at https://github.com/SchwarzIT/go-code-review to make corrections and improvements

## 📦 Configuration

Configure the application using environment variables with support for a `.env` file and default values.

### Environment Variables

- **`API_PORT`**
  - **Explanation:** API port
  - **required**
- **`API_ENV`**
  - **Explanation:** Type of environment "development"|"production"
  - **Default:** `development`
- **`API_TIMEALIVE`**
  - **Explanation:** Time max to server keep alive
  - **Format:** Specify the duration using units like s (seconds), m (minutes), h (hours), d (days), w (weeks), or y (years). Example: 10d 10w 1y.
  - **Without value:** API will run without duration
- **`API_SHUTDOWNTIMEOUT`**
  - **Explanation:** Time max to server wait the request finishing when is shutdown
  - **Format:** Specify the duration using units like s (seconds), m (minutes), h (hours), d (days), w (weeks), or y (years). Example: 10d 10w 1y.
  - **Default:** `30s`

### `.env` File

Create a `.env` file in the project root to set default values:

```dotenv
API_PORT=9090
API_ENV=production
API_TIMEALIVE=1y
API_SHUTDOWNTIMEOUT=60s
```

### Precedence

1. **System Environment Variables** override `.env` values.
2. **`.env` File** provides defaults.
3. **Struct Defaults** apply if neither is set.

**Example:**

- `.env` has `API_PORT=9090`
- System sets `API_PORT=7070`
- **Effective `API_PORT`:** `7070`

## 📄 Releases

The releases folder is used to store version-specific metadata for project updates, such as features, improvements, and bug fixes. Each release file is structured using the standard template located at ./templates/release.toml.

Workflow Automation
The repository includes a CI/CD pipeline to validate and manage release files:

1. Validation:
   When a pull request is opened or updated, the pipeline validates the release.toml file using the script `scripts/validate_release.go` to ensure it follows the required structure.
2. Merge Handling:
   When a pull request is merged, the pipeline automatically moves the validated `release.toml` file into the releases folder, renaming it with a timestamp for unique identification.

### How to Add a New Release:

1. Use the provided template at `./templates/release.toml` to create a new release file.
2. Ensure the file is added to your pull request and follows the required structure.
3. Upon merging, the CI/CD pipeline will handle validation and placement in the releases folder.
