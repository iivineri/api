# Authentication API Endpoints

## Overview

Modulul de autentificare oferă funcționalități complete pentru gestionarea utilizatorilor, sesiunilor, autentificare 2FA și managementul parolelor.

**Base URL:** `/api/v1/auth`

## Public Endpoints (fără autentificare)

### 1. Register User

**Endpoint:** `POST /api/v1/auth/register`

**Descriere:** Înregistrează un utilizator nou în sistem.

**Request Body:**
```json
{
  "nickname": "johndoe",
  "email": "john@example.com", 
  "password": "securepassword123",
  "date_of_birth": "1990-01-15T00:00:00Z"
}
```

**Validări:**
- `nickname`: obligatoriu, 3-32 caractere, doar alfanumeric
- `email`: obligatoriu, format valid de email
- `password`: obligatoriu, minim 8 caractere
- `date_of_birth`: obligatoriu, format ISO date

**Response Success (200):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "nickname": "johndoe",
    "email": "john@example.com",
    "enabled_2fa": false,
    "date_of_birth": "1990-01-15T00:00:00Z",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "email already exists",
  "error": null
}
```

---

### 2. Login User

**Endpoint:** `POST /api/v1/auth/login`

**Descriere:** Autentifică un utilizator și returnează JWT token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123",
  "totp_code": "123456"
}
```

**Validări:**
- `email`: obligatoriu, format valid
- `password`: obligatoriu, minim 6 caractere
- `totp_code`: opțional, exact 6 cifre (obligatoriu dacă 2FA este activat)

**Response Success (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nickname": "johndoe",
      "email": "john@example.com",
      "enabled_2fa": true,
      "date_of_birth": "1990-01-15T00:00:00Z",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "expires_at": "2024-01-16T10:30:00Z",
    "requires_2fa": false
  }
}
```

**Response 2FA Required (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "requires_2fa": true
  }
}
```

**Response Error (401):**
```json
{
  "success": false,
  "message": "invalid credentials",
  "error": null
}
```

---

### 3. Refresh Token

**Endpoint:** `POST /api/v1/auth/refresh`

**Descriere:** Refresh JWT token folosind refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nickname": "johndoe",
      "email": "john@example.com",
      "enabled_2fa": true,
      "date_of_birth": "1990-01-15T00:00:00Z",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "expires_at": "2024-01-16T10:30:00Z"
  }
}
```

---

### 4. Request Password Reset

**Endpoint:** `POST /api/v1/auth/password/reset`

**Descriere:** Inițiază procesul de resetare parolă prin email.

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Password reset email sent if account exists",
  "data": null
}
```

---

### 5. Confirm Password Reset

**Endpoint:** `POST /api/v1/auth/password/reset/confirm`

**Descriere:** Confirmă resetarea parolei cu token-ul primit prin email.

**Request Body:**
```json
{
  "token": "uuid-token-from-email",
  "new_password": "newsecurepassword123"
}
```

**Validări:**
- `token`: obligatoriu, format UUID
- `new_password`: obligatoriu, minim 8 caractere

**Response Success (200):**
```json
{
  "success": true,
  "message": "Password reset successfully",
  "data": null
}
```

---

## Protected Endpoints (necesită JWT token)

**Header de autentificare:** `Authorization: Bearer <jwt_token>`

### 6. Get User Profile

**Endpoint:** `GET /api/v1/auth/profile`

**Descriere:** Returnează profilul utilizatorului autentificat.

**Response Success (200):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "nickname": "johndoe",
    "email": "john@example.com",
    "enabled_2fa": true,
    "date_of_birth": "1990-01-15T00:00:00Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

---

### 7. Logout

**Endpoint:** `POST /api/v1/auth/logout`

**Descriere:** Logout din sesiunea curentă.

**Response Success (200):**
```json
{
  "success": true,
  "message": "Logged out successfully",
  "data": null
}
```

---

### 8. Logout All Devices

**Endpoint:** `POST /api/v1/auth/logout/all`

**Descriere:** Logout din toate sesiunile utilizatorului.

**Response Success (200):**
```json
{
  "success": true,
  "message": "Logged out from all devices successfully",
  "data": null
}
```

---

### 9. Change Password

**Endpoint:** `POST /api/v1/auth/password/change`

**Descriere:** Schimbă parola utilizatorului autentificat.

**Request Body:**
```json
{
  "current_password": "currentpassword123",
  "new_password": "newpassword123",
  "totp_code": "123456"
}
```

**Validări:**
- `current_password`: obligatoriu
- `new_password`: obligatoriu, minim 8 caractere
- `totp_code`: opțional, exact 6 cifre (obligatoriu dacă 2FA este activat)

**Response Success (200):**
```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

---

## Two-Factor Authentication (2FA)

### 10. Enable 2FA

**Endpoint:** `POST /api/v1/auth/2fa/enable`

**Descriere:** Inițiază procesul de activare 2FA.

**Request Body:**
```json
{
  "password": "currentpassword123"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "2FA setup initiated",
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_code_url": "otpauth://totp/Iivineri:john%40example.com?secret=JBSWY3DPEHPK3PXP&issuer=Iivineri",
    "backup_codes": [
      "ABCD1234",
      "EFGH5678",
      "IJKL9012"
    ]
  }
}
```

---

### 11. Confirm 2FA

**Endpoint:** `POST /api/v1/auth/2fa/confirm`

**Descriere:** Confirmă și activează 2FA cu codul TOTP.

**Request Body:**
```json
{
  "secret": "JBSWY3DPEHPK3PXP",
  "totp_code": "123456"
}
```

**Validări:**
- `secret`: obligatoriu
- `totp_code`: obligatoriu, exact 6 cifre

**Response Success (200):**
```json
{
  "success": true,
  "message": "2FA enabled successfully",
  "data": null
}
```

---

### 12. Disable 2FA

**Endpoint:** `POST /api/v1/auth/2fa/disable`

**Descriere:** Dezactivează 2FA pentru utilizator.

**Request Body:**
```json
{
  "password": "currentpassword123",
  "totp_code": "123456"
}
```

**Validări:**
- `password`: obligatoriu
- `totp_code`: opțional, exact 6 cifre

**Response Success (200):**
```json
{
  "success": true,
  "message": "2FA disabled successfully",
  "data": null
}
```

---

## Error Responses

### Format Standard

Toate endpoint-urile pot returna următoarele tipuri de erori:

**400 Bad Request - Validation Error:**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": [
    {
      "field": "email",
      "message": "email must be a valid email"
    },
    {
      "field": "password",
      "message": "password must be at least 8 characters long"
    }
  ]
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": null
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "user is banned",
  "error": null
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Internal server error",
  "error": null
}
```

---

## Flow de Autentificare

### 1. Înregistrare Utilizator Nou
1. `POST /api/v1/auth/register` - creează cont nou
2. `POST /api/v1/auth/login` - autentificare

### 2. Login cu 2FA
1. `POST /api/v1/auth/login` (fără totp_code) → primești `requires_2fa: true`
2. `POST /api/v1/auth/login` (cu totp_code) → primești token

### 3. Activare 2FA
1. `POST /api/v1/auth/2fa/enable` → primești secret și QR code
2. Scanează QR code în Google Authenticator
3. `POST /api/v1/auth/2fa/confirm` → confirmă cu primul cod TOTP

### 4. Reset Parolă
1. `POST /api/v1/auth/password/reset` → trimite email cu token
2. `POST /api/v1/auth/password/reset/confirm` → confirmă cu token din email

---

## Securitate

- **JWT Tokens:** Expiră în 24 ore
- **Refresh Tokens:** Expiră în 7 zile
- **Reset Tokens:** Expiră în 24 ore
- **Rate Limiting:** Aplicat în producție (100 requests/minut per IP)
- **Password Hashing:** bcrypt cu cost 12
- **2FA:** TOTP compatibil cu Google Authenticator
- **Soft Delete:** Utilizatorii nu sunt șterși fizic din baza de date