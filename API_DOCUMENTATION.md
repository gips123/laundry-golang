# 📚 API Documentation - Laundry Marketplace Backend

Dokumentasi lengkap untuk semua API endpoints yang tersedia di backend Go.

## 🌐 Base URL

```
http://localhost:8080/api/v1
```

## 🔐 Authentication

Semua endpoint yang **Protected** memerlukan header:
```
Authorization: Bearer <token>
```

Token JWT akan expire dalam 24 jam (atau sesuai `JWT_EXPIRY` di `.env`).

---

## 📋 API Endpoints

### 🔑 Authentication Endpoints

#### 1. Register User (Create Account)
**POST** `/api/v1/auth/register`

**Description:** Endpoint untuk membuat akun user baru. Setelah register berhasil, user akan langsung mendapatkan JWT token untuk autentikasi.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "081234567890",
  "address": "Jl. Sudirman No. 123, Jakarta",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "role": "customer"
}
```

**Field Validation:**
- `name` - **Required**, string
- `email` - **Required**, valid email format, must be unique
- `password` - **Required**, minimum 8 characters
- `phone` - **Required**, string
- `address` - **Required**, string
- `latitude` - **Optional**, decimal number
- `longitude` - **Optional**, decimal number
- `role` - **Optional**, default: `"customer"`, valid values: `"customer"` or `"laundry_owner"`

**Note:** 
- `latitude` dan `longitude` adalah **optional**. Bisa diisi saat register atau diupdate kemudian via `/auth/update-location`
- Jika `role` tidak disediakan, default akan menjadi `"customer"`

**Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "081234567890",
      "address": "Jl. Sudirman No. 123, Jakarta",
      "latitude": -6.2088,
      "longitude": 106.8456
    },
    "token": "jwt-token-here"
  }
}
```

---

#### 2. Login User
**POST** `/api/v1/auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "081234567890",
      "address": "Jl. Sudirman No. 123, Jakarta",
      "latitude": -6.2088,
      "longitude": 106.8456
    },
    "token": "jwt-token-here"
  }
}
```

---

#### 3. Get Current User (Protected)
**GET** `/api/v1/auth/me`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "081234567890",
    "address": "Jl. Sudirman No. 123, Jakarta",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "role": "customer"
  }
}
```

---

#### 4. Update User Location (Protected)
**PATCH** `/api/v1/auth/update-location`

**Description:** Update user's current location. **Bisa dipanggil berkali-kali** untuk mengupdate lokasi sesuai tempat user berada (rumah, kosan, kampus, dll). Location akan otomatis digunakan untuk distance calculation di endpoint laundries.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "latitude": -6.2088,
  "longitude": 106.8456
}
```

**Validation:**
- `latitude` - Required, must be between -90 and 90
- `longitude` - Required, must be between -180 and 180

**Response (200):**
```json
{
  "success": true,
  "message": "Location updated successfully",
  "data": {
    "id": "uuid",
    "latitude": -6.2088,
    "longitude": 106.8456
  }
}
```

**Use Cases:**
- User di rumah → Update location ke koordinat rumah
- User pindah ke kosan → Update location ke koordinat kosan
- User di kampus → Update location ke koordinat kampus
- Setelah update, semua request ke `/api/v1/laundries` akan otomatis menggunakan location terbaru (jika tidak ada query params lat/lng)

---

### 🧺 Laundry Endpoints

#### 5. Get All Laundries
**GET** `/api/v1/laundries`

**Query Parameters:**
- `search` (optional): Search by name or address
- `is_open` (optional): Filter by open status (`true`/`false`)
- `page` (optional): Page number (default: `1`)
- `limit` (optional): Items per page (default: `10`)
- `lat` (optional): User latitude for distance calculation
- `lng` (optional): User longitude for distance calculation
- `sort_by` (optional): Sort by `distance` or `rating` (default: `distance` if lat/lng provided, else `rating`)

**Note:** 
- Jika `lat` dan `lng` disediakan via query params, akan digunakan untuk distance calculation
- Jika tidak ada query params tapi user sudah login, akan menggunakan lat/lng dari user profile
- Laundries akan diurutkan berdasarkan **distance terdekat** jika lat/lng tersedia, else sort by **rating tertinggi**

**Example:**
```
GET /api/v1/laundries?search=express&is_open=true&page=1&limit=10&lat=-6.2088&lng=106.8456
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "laundries": [
      {
        "id": "uuid",
        "name": "Laundry Express Jakarta",
        "description": "Laundry cepat dan berkualitas",
        "address": "Jl. Sudirman No. 123, Jakarta Pusat",
        "rating": 4.8,
        "review_count": 234,
        "image": "https://...",
        "price_range": "Rp 8.000 - Rp 25.000",
        "distance": 2.5,
        "is_open": true,
        "operating_hours": {
          "open": "08:00",
          "close": "20:00"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    },
    "user_location": {
      "latitude": -6.2088,
      "longitude": 106.8456
    }
  }
}
```

---

#### 6. Get Laundry Detail
**GET** `/api/v1/laundries/:id`

**Query Parameters:**
- `lat` (optional): User latitude for distance calculation
- `lng` (optional): User longitude for distance calculation

**Example:**
```
GET /api/v1/laundries/uuid-here?lat=-6.2088&lng=106.8456
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Laundry Express Jakarta",
    "description": "Laundry cepat dan berkualitas",
    "address": "Jl. Sudirman No. 123, Jakarta Pusat",
    "rating": 4.8,
    "review_count": 234,
    "image": "https://...",
    "price_range": "Rp 8.000 - Rp 25.000",
    "distance": 2.5,
    "is_open": true,
    "operating_hours": {
      "open": "08:00",
      "close": "20:00"
    },
    "services": [
      {
        "id": "uuid",
        "name": "Cuci Reguler",
        "description": "Cuci dan setrika pakaian biasa",
        "price": 8000,
        "unit": "kg",
        "estimated_time": 24,
        "category": "regular"
      }
    ]
  }
}
```

---

### 📦 Order Endpoints (Protected)

#### 7. Create Order
**POST** `/api/v1/orders`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "laundry_id": "uuid",
  "services": [
    {
      "service_id": "uuid",
      "quantity": 3
    },
    {
      "service_id": "uuid",
      "quantity": 2
    }
  ],
  "delivery_address": "Jl. Contoh No. 123, Jakarta",
  "notes": "Mohon hati-hati dengan pakaian putih",
  "estimated_pickup_at": "2024-01-15T14:00:00Z"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "uuid",
    "laundry_id": "uuid",
    "laundry_name": "Laundry Express Jakarta",
    "services": [
      {
        "service_id": "uuid",
        "service_name": "Cuci Reguler",
        "quantity": 3,
        "price": 8000,
        "unit": "kg"
      }
    ],
    "total_price": 34000,
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "estimated_pickup": "2024-01-15T14:00:00Z",
    "estimated_delivery": "2024-01-16T18:00:00Z",
    "address": "Jl. Contoh No. 123, Jakarta",
    "notes": "Mohon hati-hati dengan pakaian putih"
  }
}
```

---

#### 8. Get User Orders
**GET** `/api/v1/orders`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `status` (optional): Filter by status (`pending`, `confirmed`, `washing`, `ready`, `delivered`, `completed`, `cancelled`)
- `page` (optional): Page number (default: `1`)
- `limit` (optional): Items per page (default: `10`)

**Example:**
```
GET /api/v1/orders?status=pending&page=1&limit=10
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "orders": [
      {
        "id": "uuid",
        "laundry_id": "uuid",
        "laundry_name": "Laundry Express Jakarta",
        "services": [
          {
            "service_id": "uuid",
            "service_name": "Cuci Reguler",
            "quantity": 3,
            "price": 8000,
            "unit": "kg"
          }
        ],
        "total_price": 34000,
        "status": "washing",
        "created_at": "2024-01-15T10:30:00Z",
        "estimated_pickup": "2024-01-15T14:00:00Z",
        "estimated_delivery": "2024-01-16T18:00:00Z",
        "address": "Jl. Contoh No. 123, Jakarta"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

#### 9. Get Order Detail
**GET** `/api/v1/orders/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "laundry_id": "uuid",
    "laundry_name": "Laundry Express Jakarta",
    "services": [
      {
        "service_id": "uuid",
        "service_name": "Cuci Reguler",
        "quantity": 3,
        "price": 8000,
        "unit": "kg"
      }
    ],
    "total_price": 34000,
    "status": "washing",
    "created_at": "2024-01-15T10:30:00Z",
    "estimated_pickup": "2024-01-15T14:00:00Z",
    "estimated_delivery": "2024-01-16T18:00:00Z",
    "address": "Jl. Contoh No. 123, Jakarta",
    "notes": "Mohon hati-hati dengan pakaian putih"
  }
}
```

---

#### 10. Cancel Order
**PATCH** `/api/v1/orders/:id/cancel`

**Headers:**
```
Authorization: Bearer <token>
```

**Note:** Hanya bisa cancel order dengan status `pending` atau `confirmed`.

**Response (200):**
```json
{
  "success": true,
  "message": "Order cancelled successfully",
  "data": {
    "id": "uuid",
    "status": "cancelled",
    ...
  }
}
```

---

#### 11. Update Order Status (Laundry Owner Only)
**PATCH** `/api/v1/orders/:id/status`

**Headers:**
```
Authorization: Bearer <token>
```

**Note:** Hanya bisa diakses oleh user dengan role `laundry_owner`.

**Request Body:**
```json
{
  "status": "confirmed"
}
```

**Valid Status Values:**
- `pending`
- `confirmed`
- `picked-up`
- `washing`
- `drying`
- `ironing`
- `ready`
- `delivered`
- `completed`
- `cancelled`

**Response (200):**
```json
{
  "success": true,
  "message": "Order status updated",
  "data": {
    "id": "uuid",
    "status": "confirmed",
    ...
  }
}
```

---

## 🏥 Health Check

#### Health Check
**GET** `/health`

**Response (200):**
```json
{
  "status": "ok"
}
```

---

## 📝 Response Format

Semua response mengikuti format standar:

**Success Response:**
```json
{
  "success": true,
  "message": "Optional message",
  "data": {...}
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error message"
}
```

---

## 🚨 Error Codes

- **400** - Bad Request (validation errors)
- **401** - Unauthorized (invalid/missing token)
- **403** - Forbidden (insufficient permissions)
- **404** - Not Found
- **500** - Internal Server Error

---

## 🔑 Test Credentials

Setelah menjalankan seed script, Anda bisa menggunakan credentials berikut:

- **Customer:**
  - Email: `user@laundryhub.com`
  - Password: `password123`

- **Laundry Owner:**
  - Email: `admin@laundryhub.com`
  - Password: `admin123`

- **Test User:**
  - Email: `test@test.com`
  - Password: `test123`

---

## 📌 Important Notes

1. **Distance Calculation:**
   - Prioritas penggunaan lat/lng: Query params > User profile > Null
   - Distance dihitung menggunakan Haversine formula
   - Return dalam kilometers (rounded to 2 decimal places)

2. **Sorting:**
   - Jika lat/lng tersedia: Sort by distance (nearest first)
   - Jika tidak ada lat/lng: Sort by rating (highest first)

3. **Pagination:**
   - Default: `page=1`, `limit=10`
   - Response include `pagination` object dengan `total` dan `total_pages`

4. **Order Status Flow:**
   - `pending` → `confirmed` → `picked-up` → `washing` → `drying` → `ironing` → `ready` → `delivered` → `completed`
   - Bisa `cancelled` dari status `pending` atau `confirmed`

5. **Price Range:**
   - Dihitung dari min/max price dari services yang aktif
   - Format: "Rp {min} - Rp {max}"

---

## 🧪 Example Usage

### Complete Flow: Register → Update Location → Get Laundries

```javascript
// 1. Register (Create Account)
const registerResponse = await fetch('http://localhost:8080/api/v1/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    name: 'John Doe',
    email: 'john@example.com',
    password: 'password123',
    phone: '081234567890',
    address: 'Jl. Sudirman No. 123, Jakarta',
    latitude: -6.2088,  // Optional - bisa diisi saat register
    longitude: 106.8456, // Optional
    role: 'customer'
  })
});

const registerData = await registerResponse.json();
if (registerData.success) {
  const token = registerData.data.token;
  localStorage.setItem('token', token);
}

// 2. Request Location Permission (Frontend)
navigator.geolocation.getCurrentPosition(
  async (position) => {
    const lat = position.coords.latitude;
    const lng = position.coords.longitude;
    
    // 3. Update Location to Database (Bisa dipanggil berkali-kali)
    // Contoh: User di rumah
    await updateUserLocation(token, lat, lng);
    
    // 4. Get Laundries (Otomatis menggunakan location dari profile)
    const laundriesResponse = await fetch(
      'http://localhost:8080/api/v1/laundries', // Tidak perlu query params!
      {
        headers: { 'Authorization': `Bearer ${token}` }
      }
    );
    // Laundries akan otomatis diurutkan berdasarkan distance dari location terbaru
  },
  (error) => {
    console.error('Location access denied:', error);
  }
);

// Function untuk update location (bisa dipanggil kapan saja)
async function updateUserLocation(token, lat, lng) {
  const response = await fetch('http://localhost:8080/api/v1/auth/update-location', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ latitude: lat, longitude: lng })
  });
  return response.json();
}

// Contoh: User pindah ke kosan, update location lagi
// updateUserLocation(token, newLat, newLng);

// Contoh: User di kampus, update location lagi
// updateUserLocation(token, campusLat, campusLng);
```

### Login Flow

```javascript
// 1. Login
const loginResponse = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'john@example.com',
    password: 'password123'
  })
});

const { data } = await loginResponse.json();
const token = data.token;
localStorage.setItem('token', token);

// 2. Get Current Location (Frontend)
navigator.geolocation.getCurrentPosition(
  async (position) => {
    const lat = position.coords.latitude;
    const lng = position.coords.longitude;
    
    // 3. Update Location
    await updateUserLocation(token, lat, lng);
    
    // 4. Get Laundries (Auto-use location dari profile)
    const laundries = await fetch('http://localhost:8080/api/v1/laundries', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
  }
);
```

---

**Last Updated:** December 2024

