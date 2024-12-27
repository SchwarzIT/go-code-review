# Releases Folder

This folder contains all the release files for the project. Each release file documents the details of a specific version, including metadata, features, improvements, or bug fixes.

## Folder Purpose

The `releases` folder is intended to store version-specific metadata for project releases. It serves as a structured repository of information for tracking updates and ensuring that all release notes are consistently formatted.

## Release File Structure

Each release file follows a standard template to ensure clarity and uniformity. Below is the structure of a typical release file:

### [Metadata]

- **Author**: The name of the person creating the release.
- **Release File Version**: The version of the release metadata format.

### [Description]

- **Notes**: General notes about the release. Include details or summaries as necessary.

### [[Digest.Types]]

A release file can only contain one type of digest at a time. Choose from the following:

#### Features Digest

- **Name**: The name of the new feature.
- **Issue**: The related issue number (if applicable).
- **Description**: A detailed description of the feature.

#### Improvements Digest

- **Name**: The name of the improvement.
- **Issue**: The related issue number (if applicable).
- **Description**: A detailed description of the improvement.

#### Bugs Digest

- **Name**: The name of the bug fixed.
- **Issue**: The related issue number (if applicable).
- **Description**: A detailed description of the bug and how it was resolved.

> **Note**: Each release file must focus on a single type of digest (Features, Improvements, or Bugs). If you need to document multiple types, create separate release files for each.

## Template for Release Files

To create a new release file, use the template provided at `./templates/release.toml`. This template ensures all release files follow the required structure and format.

### How to Use the Template:

1. Copy the `release.toml` template to the `releases` folder:
   ```bash
   cp ./templates/release.toml ./releases/release-v<version>.toml
   ```
