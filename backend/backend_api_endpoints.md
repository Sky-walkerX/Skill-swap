# Skill Swap Backend API Endpoints

This document lists all REST API endpoints for the Skill Swap backend (Go + Gin). Each endpoint includes:
- **Method & Path**
- **Authentication & Headers**
- **Request Body Example**
- **Response Example**
- **Description**

---

## Authentication

### Register
- **POST** `/api/v1/auth/register`
- **Headers:** `Content-Type: application/json`
- **Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "yourpassword"
}
```
- **Response:**
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": { "user_id": "...", "name": "John Doe", "email": "john@example.com" }
}
```
- **Description:** Register a new user. Returns tokens and user info.

### Login
- **POST** `/api/v1/auth/login`
- **Headers:** `Content-Type: application/json`
- **Body:**
```json
{
  "email": "john@example.com",
  "password": "yourpassword"
}
```
- **Response:**
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": { "user_id": "...", "name": "John Doe", "email": "john@example.com" }
}
```
- **Description:** Authenticate and receive tokens.

### Refresh Token
- **POST** `/api/v1/auth/refresh`
- **Headers:** `Content-Type: application/json`
- **Body:**
```json
{
  "refresh_token": "..."
}
```
- **Response:**
```json
{
  "access_token": "...",
  "refresh_token": "..."
}
```
- **Description:** Get a new access token using a refresh token.

### Logout
- **POST** `/api/v1/auth/logout`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "message": "Logged out successfully" }
```
- **Description:** Logout user (client should discard tokens).

### Get Current User
- **GET** `/api/v1/auth/me`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "user_id": "...", "email": "john@example.com" }
```
- **Description:** Get info from JWT. For full profile, use `/users/profile`.

---

## Users

### Search Users
- **GET** `/api/v1/public/users/search?location=...&search_term=...&page=1&limit=10`
- **Response:**
```json
{
  "users": [ { "user_id": "...", "name": "...", ... } ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```
- **Description:** Search users by location or name.

### Get Profile
- **GET** `/api/v1/users/profile`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "user_id": "...", "name": "...", "email": "...", ... }
```
- **Description:** Get authenticated user's profile.

### Update Profile
- **PUT** `/api/v1/users/profile`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{ "name": "New Name", "location": "New City" }
```
- **Response:**
```json
{ "message": "Profile updated successfully" }
```
- **Description:** Update authenticated user's profile.

---

## Skills

### List All Skills
- **GET** `/api/v1/skills`
- **Response:**
```json
[
  { "skill_id": "...", "name": "Guitar", "created_at": "..." }
]
```
- **Description:** List all available skills.

### Get Skill by ID
- **GET** `/api/v1/skills/{id}`
- **Response:**
```json
{ "skill_id": "...", "name": "Guitar", "created_at": "..." }
```
- **Description:** Get details of a skill.

#### User's Offered/Wanted Skills (Protected)
- **GET** `/api/v1/users/skills/offered` — List offered
- **POST** `/api/v1/users/skills/offered` — Add offered
- **DELETE** `/api/v1/users/skills/offered/{id}` — Remove offered
- **GET** `/api/v1/users/skills/wanted` — List wanted
- **POST** `/api/v1/users/skills/wanted` — Add wanted
- **DELETE** `/api/v1/users/skills/wanted/{id}` — Remove wanted
- **Headers:** `Authorization: Bearer <access_token>`
- **Body for POST:**
```json
{ "skill_id": "..." }
```
- **Response:**
  - GET: Array of skills
  - POST/DELETE: 204 No Content

---

## Swaps

### Create Swap Request
- **POST** `/api/v1/swaps`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{
  "responder_id": "...",
  "offered_skill_id": "...",
  "wanted_skill_id": "..."
}
```
- **Response:**
```json
{ "swap_id": "...", "status": "pending", ... }
```
- **Description:** Create a new swap request.

### Get User's Swap Requests
- **GET** `/api/v1/swaps?status=pending&sent=true&received=true&limit=10&offset=0`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{
  "sent": [ { "swap_id": "...", ... } ],
  "received": [ { "swap_id": "...", ... } ]
}
```
- **Description:** List swap requests for the user.

### Get Swap by ID
- **GET** `/api/v1/swaps/{id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "swap_id": "...", ... }
```
- **Description:** Get details of a swap request.

### Update Swap Status
- **PUT** `/api/v1/swaps/{id}/status`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{ "status": "accepted" }
```
- **Response:**
```json
{ "swap_id": "...", "status": "accepted", ... }
```
- **Description:** Accept, reject, or cancel a swap request.

### Delete Swap Request
- **DELETE** `/api/v1/swaps/{id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:** 204 No Content
- **Description:** Delete a swap request (requester only).

### Get Potential Matches
- **GET** `/api/v1/swaps/matches`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
[
  { "user": { ... }, "offered_skill": { ... }, "wanted_skill": { ... }, "match_score": 90 }
]
```
- **Description:** Find potential swap matches for the user.

---

## Ratings

### Create Rating
- **POST** `/api/v1/ratings`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{
  "swap_id": "...",
  "score": 5,
  "comment": "Great experience!"
}
```
- **Response:**
```json
{ "rating_id": "...", "score": 5, ... }
```
- **Description:** Rate a completed swap.

### Get Rating by ID
- **GET** `/api/v1/ratings/{id}`
- **Response:**
```json
{ "rating_id": "...", "score": 5, ... }
```
- **Description:** Get a specific rating.

### Update Rating
- **PUT** `/api/v1/ratings/{id}`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{ "score": 4, "comment": "Updated comment" }
```
- **Response:**
```json
{ "rating_id": "...", "score": 4, ... }
```
- **Description:** Update a rating (only by the rater).

### Delete Rating
- **DELETE** `/api/v1/ratings/{id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:** 204 No Content
- **Description:** Delete a rating (only by the rater).

### Get Ratings for a Swap
- **GET** `/api/v1/ratings/swap/{swap_id}`
- **Response:**
```json
{ "ratings": [ { "rating_id": "...", ... } ] }
```
- **Description:** Get all ratings for a swap.

### Get User Ratings
- **GET** `/api/v1/users/{user_id}/ratings?as_rater=true&as_ratee=true&min_score=1&max_score=5&limit=10&offset=0`
- **Response:**
```json
{ "ratings": [ { "rating_id": "...", ... } ] }
```
- **Description:** Get ratings given/received by a user.

### Get User Rating Stats
- **GET** `/api/v1/users/{user_id}/ratings/stats`
- **Response:**
```json
{ "average_score": 4.8, "total_ratings": 10, ... }
```
- **Description:** Get rating statistics for a user.

---

## Availability

### Create Availability Slot
- **POST** `/api/v1/availability`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{
  "day": 1,
  "start_time": "09:00",
  "end_time": "11:00"
}
```
- **Response:**
```json
{ "slot_id": "...", "day": 1, "start_time": "09:00", "end_time": "11:00", ... }
```
- **Description:** Create a new availability slot.

### Get User's Availability Slots
- **GET** `/api/v1/availability`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "availability_slots": [ { "slot_id": "...", ... } ] }
```
- **Description:** List all availability slots for the user.

### Get Availability Slot by ID
- **GET** `/api/v1/availability/{id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "slot_id": "...", ... }
```
- **Description:** Get a specific availability slot.

### Update Availability Slot
- **PUT** `/api/v1/availability/{id}`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{ "day": 2, "start_time": "10:00", "end_time": "12:00" }
```
- **Response:**
```json
{ "slot_id": "...", "day": 2, ... }
```
- **Description:** Update an availability slot.

### Delete Availability Slot
- **DELETE** `/api/v1/availability/{id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:** 204 No Content
- **Description:** Delete an availability slot.

### Find Common Availability
- **GET** `/api/v1/availability/common/{user_id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "common_availability": [ { "slot_id": "...", ... } ] }
```
- **Description:** Find overlapping availability with another user.

### Search Availability by Day/Time
- **GET** `/api/v1/availability/search?day=1&start_time=09:00&end_time=11:00`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "availability_slots": [ { "slot_id": "...", ... } ] }
```
- **Description:** Find slots for a specific day/time range.

---

## Notifications

### Get Notifications
- **GET** `/api/v1/notifications?page=1&limit=10&unread_only=true`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{
  "notifications": [ { "notification_id": "...", ... } ],
  "pagination": { "page": 1, "limit": 10, "total": 1 }
}
```
- **Description:** List notifications for the user.

### Mark Notifications as Read
- **PUT** `/api/v1/notifications/mark-read`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{ "notification_ids": ["..."] }
```
- **Response:**
```json
{ "message": "Notifications marked as read" }
```
- **Description:** Mark specific notifications as read.

### Mark All as Read
- **PUT** `/api/v1/notifications/mark-all-read`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "message": "All notifications marked as read" }
```
- **Description:** Mark all notifications as read.

### Delete Notification
- **DELETE** `/api/v1/notifications/{id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "message": "Notification deleted successfully" }
```
- **Description:** Delete a notification.

### Get Notification Stats
- **GET** `/api/v1/notifications/stats`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "total": 10, "read": 8, "unread": 2 }
```
- **Description:** Get notification statistics.

### Create Notification (Admin)
- **POST** `/api/v1/notifications`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **Body:**
```json
{ "type": "info", "title": "...", "message": "...", "related_id": "..." }
```
- **Response:**
```json
{ "notification_id": "...", ... }
```
- **Description:** Create a notification (admin only).

### Get Notification by ID
- **GET** `/api/v1/notifications/{id}`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "notification_id": "...", ... }
```
- **Description:** Get a specific notification.

---

## Files

### Upload User Photo
- **POST** `/api/v1/files/users/photo`
- **Headers:** `Authorization: Bearer <access_token>`, `Content-Type: multipart/form-data`
- **Form Data:**
  - `file`: (file upload, jpg/png/gif/webp, max 5MB)
- **Response:**
```json
{ "url": "/api/v1/files/users/{user_id}/{filename}", ... }
```
- **Description:** Upload a profile photo.

### Delete User Photo
- **DELETE** `/api/v1/files/users/photo`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "message": "Photo deleted successfully", "success": true }
```
- **Description:** Delete the user's profile photo.

### Get User Photo
- **GET** `/api/v1/files/users/{user_id}/{filename}`
- **Response:** (image file)
- **Description:** Download a user's profile photo.

### Get User Photo Info
- **GET** `/api/v1/files/users/{user_id}/{filename}/info`
- **Headers:** `Authorization: Bearer <access_token>`
- **Response:**
```json
{ "filename": "...", "size": 12345, ... }
```
- **Description:** Get info about a user's photo.

---

## Search

### Global Search
- **GET** `/api/v1/search/global?q=term&types=users,skills,swaps&limit=5`
- **Response:**
```json
{ "users": [...], "skills": [...], "swaps": [...] }
```
- **Description:** Search across users, skills, and swaps.

### Search Suggestions
- **GET** `/api/v1/search/suggestions?q=term&type=skills`
- **Response:**
```json
{ "suggestions": ["Guitar", "Piano"], "query": "G", "type": "skills" }
```
- **Description:** Get autocomplete suggestions.

### Advanced User Search
- **GET** `/api/v1/search/users?...`
- **Headers:** `Authorization: Bearer <access_token>` (for advanced search)
- **Query Params:** `q`, `location`, `skills_offered`, `skills_wanted`, `min_rating`, `is_public`, `sort_by`, `sort_order`, `limit`, `offset`
- **Response:**
```json
{ "users": [...], "total": 1, "limit": 10, "offset": 0 }
```
- **Description:** Advanced user search with filters.

### Advanced Swap Search
- **GET** `/api/v1/search/swaps?...`
- **Headers:** `Authorization: Bearer <access_token>`
- **Query Params:** `q`, `status`, `offered_skill_id`, `wanted_skill_id`, `requester_id`, `responder_id`, `created_after`, `created_before`, `sort_by`, `sort_order`, `limit`, `offset`
- **Response:**
```json
{ "swaps": [...], "total": 1, "limit": 10, "offset": 0 }
```
- **Description:** Advanced swap search with filters.

### Advanced Skill Search
- **GET** `/api/v1/search/skills?...`
- **Headers:** `Authorization: Bearer <access_token>`
- **Query Params:** `q`, `category`, `sort_by`, `sort_order`, `limit`, `offset`
- **Response:**
```json
{ "skills": [...], "total": 1, "limit": 10, "offset": 0 }
```
- **Description:** Advanced skill search with filters.

---

## Admin Endpoints

All admin endpoints require `Authorization: Bearer <access_token>` and admin privileges.

### Manage Skills
- **POST** `/api/v1/admin/skills` — Create skill
- **PUT** `/api/v1/admin/skills/{id}` — Update skill
- **DELETE** `/api/v1/admin/skills/{id}` — Delete skill
- **Body:** `{ "name": "Skill Name" }`
- **Response:** Skill object or 204 No Content

### Manage Users
- **GET** `/api/v1/admin/users?...` — List users
- **PUT** `/api/v1/admin/users/{id}/ban` — Ban user
- **PUT** `/api/v1/admin/users/{id}/unban` — Unban user
- **DELETE** `/api/v1/admin/users/{id}` — Delete user
- **PUT** `/api/v1/admin/users/{id}/make-admin` — Make admin
- **PUT** `/api/v1/admin/users/{id}/remove-admin` — Remove admin

### Manage Swaps
- **GET** `/api/v1/admin/swaps?...` — List swaps
- **PUT** `/api/v1/admin/swaps/{id}/cancel` — Cancel swap (body: `{ "reason": "..." }`)

### Platform Stats & Reports
- **GET** `/api/v1/admin/stats` — Platform statistics
- **GET** `/api/v1/admin/reports` — Reported content

---

## Health Checks

- **GET** `/health` — Service health
- **GET** `/ready` — Readiness probe
- **GET** `/live` — Liveness probe

---

## Error Responses

Most endpoints return errors in the form:
```json
{ "error": "Error message" }
```

- **401 Unauthorized:** Missing/invalid token
- **403 Forbidden:** Insufficient privileges
- **404 Not Found:** Resource not found
- **400 Bad Request:** Invalid input
- **409 Conflict:** Duplicate or already exists
- **500 Internal Server Error:** Server error

---

## Authentication

- For protected endpoints, send `Authorization: Bearer <access_token>` in headers.
- For file uploads, use `multipart/form-data`.
- For JSON requests, use `Content-Type: application/json`.

---

## Notes
- All IDs are UUID strings.
- Pagination: `page`, `limit`, `offset` as query params.
- All times are ISO8601 strings.
- Some endpoints may require admin privileges. 