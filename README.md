# Overview

This repository was created from the repository at https://github.com/SchwarzIT/go-code-review to make corrections and improvements

## Releases

The releases folder is used to store version-specific metadata for project updates, such as features, improvements, and bug fixes. Each release file is structured using the standard template located at ./templates/release.toml.

Workflow Automation
The repository includes a CI/CD pipeline to validate and manage release files:

1. Validation:
   When a pull request is opened or updated, the pipeline validates the release.toml file using the script scripts/validate_release.go to ensure it follows the required structure.
2. Merge Handling:
   When a pull request is merged, the pipeline automatically moves the validated release.toml file into the releases folder, renaming it with a timestamp for unique identification.

### How to Add a New Release:

1. Use the provided template at ./templates/release.toml to create a new release file.
2. Ensure the file is added to your pull request and follows the required structure.
3. Upon merging, the CI/CD pipeline will handle validation and placement in the releases folder.
