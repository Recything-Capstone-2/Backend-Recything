# Recything API Documentation

## Introduction
Recything API is a backend service designed to facilitate user interactions for waste management and reporting systems. This API supports features like user registration, reporting rubbish, and managing user points.

## Features

| Feature ID | Feature Name                     | Description                                                                                  | Endpoint                                    | Method | Auth Required |
|------------|----------------------------------|----------------------------------------------------------------------------------------------|--------------------------------------------|--------|---------------|
| 1          | Register                         | Register a new user with details such as name, email, and password.                          | `/api/v1/register`                         | POST   | No            |
| 2          | Login Admin                      | Admin login using credentials.                                                              | `/api/v1/login`                            | POST   | No            |
| 3          | Login User                       | User login using credentials.                                                               | `/api/v1/login`                            | POST   | No            |
| 4          | Logout                           | Logout current session for user or admin.                                                   | `/api/v1/logout`                           | GET    | Yes           |
| 5          | Update Photo                     | Update the profile photo for the user or admin.                                             | `/api/v1/users`                            | PUT    | Yes           |
| 6          | Update User Data                 | Update user details such as email, phone, and password.                                      | `/api/v1/user/data/:iduser`                | PUT    | Yes           |
| 7          | Get User Points                  | Retrieve points associated with a user.                                                     | `/api/v1/users/points`                     | GET    | Yes           |
| 8          | Admin: Get All User Points       | Fetch points for all users.                                                                 | `/api/v1/admin/users/points`               | GET    | Yes           |
| 9          | Admin: Deduct Points             | Reduce points for a user as part of a reward mechanism.                                     | `/api/v1/admin/users/points/deduct`        | POST   | Yes           |
| 10         | Admin: Get All Users             | Retrieve all users in the system.                                                           | `/api/v1/admin/users`                      | GET    | Yes           |
| 11         | Admin: Get User by ID            | Retrieve a specific user based on their ID.                                                 | `/api/v1/admin/users/:id`                  | GET    | Yes           |
| 12         | User: Add Rubbish Report         | Report rubbish by providing location, description, and a photo.                             | `/api/v1/report-rubbish`                   | POST   | Yes           |
| 13         | Admin: Get All Rubbish Reports   | Retrieve all rubbish reports with pagination options.                                       | `/api/v1/admin/report-rubbish`             | GET    | Yes           |
| 14         | Admin: Filter Rubbish Reports    | Filter rubbish reports by status or sorting.                                                | `/api/v1/admin/report-rubbish`             | GET    | Yes           |
| 15         | Admin: Get Report by ID          | Retrieve specific rubbish report details.                                                   | `/api/v1/admin/report-rubbish/:id`         | GET    | Yes           |
| 16         | Admin: Delete Report             | Delete a specific rubbish report by ID.                                                     | `/api/v1/admin/report-rubbish/:id`         | DELETE | Yes           |
| 17         | Admin: Get Latest Reports        | Retrieve the latest 10 rubbish reports.                                                     | `/api/v1/admin/latest-report`              | GET    | Yes           |
| 18         | Admin: Update Report Status      | Change the status of a report (e.g., approved, rejected, completed).                        | `/api/v1/report-rubbish/:idreport`         | PUT    | Yes           |
| 19         | Admin: Add Article               | Publish a new article with content, author, and multimedia links.                           | `/api/v1/admin/articles`                   | POST   | Yes           |
| 20         | Admin: Update Article            | Modify an existing article's details.                                                       | `/api/v1/admin/articles/:id`               | PUT    | Yes           |
| 21         | Admin: Delete Article            | Remove an article by its ID.                                                                | `/api/v1/admin/article/:id`                | DELETE | Yes           |
| 22         | User: Get All Articles           | Retrieve all articles available.                                                            | `/api/v1/articles`                         | GET    | Yes           |
| 23         | User: Get Article by ID          | Fetch a specific article by its ID.                                                         | `/api/v1/articles/:id`                     | GET    | Yes           |
| 24         | User: Get Rubbish Report History | Retrieve the history of rubbish reports made by the user.                                   | `/api/v1/report-rubbish/history`           | GET    | Yes           |
| 25         | Admin: Statistics                | View statistics related to rubbish reports.                                                 | `/api/v1/admin/reports/statistics`         | GET    | Yes           |
| 26         | Admin: Add Reward                | Add a reward to the user for specific achievements.                                         | `/api/v1/admin/users/reward`               | POST   | Yes           |

## Authentication
Certain endpoints require a Bearer token for authentication. Tokens are issued upon successful login and should be included in the `Authorization` header.

## Getting Started
1. Clone this repository.
2. Navigate to the project directory.
3. Run the following command to install dependencies:
   ```bash
   go mod tidy
   ```
4. Start the application using:
   ```bash
   go run main.go
   ```
5. Follow the API endpoints and authentication process to integrate.

## Additional Resources
- [Postman Collection](Recything-Capstone2.postman_collection.json) for testing.
