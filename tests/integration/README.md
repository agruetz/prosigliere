# Integration Tests for Blog API

This directory contains integration tests for the Blog API using Robot Framework.

## Prerequisites

Before running the tests, you need to:

1. Install Robot Framework and required libraries:
   ```
   make install-robot
   ```

2. Build the server:
   ```
   make build-server
   ```

3. Ensure the database is set up and migrations are applied:
   ```
   make db-migrate
   ```

## Running the Tests

There are several ways to run the integration tests:

### Automated Test Run

To build the server, start it, run all tests, and stop the server automatically:

```
make run-integration-tests
```

### Manual Test Run

If you want more control over the testing process:

1. Start the server in the background:
   ```
   make start-server-for-test
   ```

2. Run the tests:
   ```
   make integration-test-all
   ```

3. Stop the server:
   ```
   make stop-server-for-test
   ```

### Running Individual Test Files

To run specific test files:

```
make integration-test
```

Or run a specific test file directly:

```
robot tests/integration/robot/blog_crud_tests.robot
```

## Test Files

The integration tests are organized into the following files:

- `blog_crud_tests.robot`: Tests for creating, reading, updating, and deleting blog posts
- `blog_list_tests.robot`: Tests for listing blog posts and pagination
- `blog_comment_tests.robot`: Tests for adding and retrieving comments on blog posts
- `blog_error_tests.robot`: Tests for error handling and edge cases

## Common Resources

The `common.resource` file contains shared keywords and variables used by all test files.

## Test Results

After running the tests, Robot Framework generates HTML reports and logs in the current directory:

- `report.html`: Summary report of the test execution
- `log.html`: Detailed log of the test execution
- `output.xml`: XML output file containing the test results