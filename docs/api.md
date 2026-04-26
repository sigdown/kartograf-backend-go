# API

## Auth

- Auth type: JWT Bearer
- Header: `Authorization: Bearer <token>`
- `POST /auth/register` and `POST /auth/login` return `token` and `user`
- Admin routes require `role = ADMIN`

## Notes

- Map `slug` is unique and immutable after creation
- Active PMTiles object is stored by stable key: `<slug>.pmtiles`
- Archive upload goes through presigned upload URL
- Archive download returns presigned download URL

## Public Endpoints

### `GET /health`

- Returns service health

### `POST /auth/register`

- Registers user
- Body:

```json
{
  "username": "user",
  "display_name": "User",
  "email": "user@example.com",
  "password": "password123"
}
```

### `POST /auth/login`

- Logs user in
- Body:

```json
{
  "login": "user@example.com",
  "password": "password123"
}
```

### `GET /maps`

- Returns public catalog of maps

### `GET /maps/:slug`

- Returns map details by slug

## Authorized Endpoints

### `GET /auth/me`

- Returns current authenticated user

### `GET /maps/by-id/:id/download`

- Returns presigned archive download URL

### `GET /points`

- Returns current user remote points

### `POST /points`

- Creates remote point
- Body:

```json
{
  "name": "Point 1",
  "description": "demo",
  "lat": 55.75,
  "lon": 37.61
}
```

### `PATCH /points/:id`

- Updates remote point

### `DELETE /points/:id`

- Deletes remote point

### `PATCH /account`

- Updates current user account
- Body can include:
  - `username`
  - `display_name`
  - `email`
  - `password`

### `DELETE /account`

- Deletes current user account

## Admin Endpoints

### `POST /admin/maps/upload-url`

- Starts map archive upload
- Returns presigned `PUT` URL and storage metadata
- Body:

```json
{
  "slug": "old-map",
  "title": "Old Map",
  "description": "Demo map",
  "year": 1943,
  "archive_name": "old-map.pmtiles",
  "archive_mime_type": "application/pmtiles"
}
```

### `POST /admin/maps`

- Finalizes map creation after file upload to storage
- Body:

```json
{
  "map_id": "uuid",
  "archive_id": "uuid",
  "storage_key": "old-map.pmtiles",
  "slug": "old-map",
  "title": "Old Map",
  "description": "Demo map",
  "year": 1943
}
```

### `PATCH /admin/maps/:id`

- Updates map metadata
- `slug` is not updated here
- Body:

```json
{
  "title": "Old Map Updated",
  "description": "Updated description",
  "year": 1944
}
```

### `POST /admin/maps/:id/archive/upload-url`

- Starts active archive replacement
- Returns presigned `PUT` URL and storage metadata
- Body:

```json
{
  "archive_name": "old-map.pmtiles",
  "archive_mime_type": "application/pmtiles"
}
```

### `PUT /admin/maps/:id/archive`

- Finalizes archive replacement after upload
- Body:

```json
{
  "archive_id": "uuid",
  "storage_key": "old-map.pmtiles"
}
```

### `DELETE /admin/maps/:id`

- Deletes map
