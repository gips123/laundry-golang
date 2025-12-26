# 📮 Postman Collection Setup Guide

Panduan lengkap untuk menggunakan Postman Collection untuk testing API Laundry Marketplace.

## 📥 Import Collection

### 1. Import Collection
1. Buka Postman
2. Klik **Import** button (kiri atas)
3. Pilih file `Laundry_Marketplace_API.postman_collection.json`
4. Klik **Import**

### 2. Import Environment (Optional tapi Recommended)
1. Klik **Import** button lagi
2. Pilih file `Laundry_Marketplace_API.postman_environment.json`
3. Klik **Import**
4. Pilih environment **"Laundry Marketplace - Local"** di dropdown environment (kanan atas)

## 🔧 Setup Environment Variables

Setelah import environment, pastikan variables berikut sudah di-set:

- `base_url` - `http://localhost:8080` (default)
- `auth_token` - Akan di-set otomatis setelah login/register
- `user_id` - Akan di-set otomatis setelah login/register
- `laundry_id` - Set manual setelah mendapatkan laundry ID
- `service_id` - Set manual setelah mendapatkan service ID
- `order_id` - Set manual setelah membuat order

## 🚀 Quick Start

### Step 1: Health Check
1. Pilih **Health Check** → **Health Check**
2. Klik **Send**
3. Harus return `{"status": "ok"}`

### Step 2: Register atau Login
**Option A: Register (Create Account)**
1. Pilih **Authentication** → **Register User**
2. Edit request body sesuai kebutuhan
3. Klik **Send**
4. Token akan otomatis tersimpan di environment variable

**Option B: Login (Gunakan test credentials)**
1. Pilih **Authentication** → **Login**
2. Request body sudah diisi dengan test credentials:
   ```json
   {
     "email": "user@laundryhub.com",
     "password": "password123"
   }
   ```
3. Klik **Send**
4. Token akan otomatis tersimpan di environment variable

### Step 3: Test Protected Endpoints
Setelah login/register, semua protected endpoints akan otomatis menggunakan token dari environment.

## 📋 Collection Structure

### 1. Authentication (4 endpoints)
- ✅ **Register User** - Create new account (auto-save token)
- ✅ **Login** - Login user (auto-save token)
- ✅ **Get Current User** - Get authenticated user info
- ✅ **Update User Location** - Update user lat/lng

### 2. Laundries (2 endpoints)
- ✅ **Get All Laundries** - List dengan search, filter, pagination
- ✅ **Get Laundry Detail** - Get detail dengan services

### 3. Orders (5 endpoints - semua Protected)
- ✅ **Create Order** - Create new order
- ✅ **Get User Orders** - List orders dengan filter
- ✅ **Get Order Detail** - Get order detail
- ✅ **Cancel Order** - Cancel order
- ✅ **Update Order Status** - Update status (Owner only)

### 4. Health Check
- ✅ **Health Check** - Check server status

## 🔑 Auto Token Management

Collection ini sudah dikonfigurasi untuk:
- ✅ Auto-save token setelah login/register
- ✅ Auto-use token untuk semua protected endpoints
- ✅ Auto-save user_id setelah login/register

## 📝 Test Credentials

Setelah menjalankan seed script, gunakan credentials berikut:

**Customer:**
- Email: `user@laundryhub.com`
- Password: `password123`

**Laundry Owner:**
- Email: `admin@laundryhub.com`
- Password: `admin123`

**Test User:**
- Email: `test@test.com`
- Password: `test123`

## 🧪 Testing Flow

### Complete Testing Flow:

1. **Health Check** → Verify server is running
2. **Register User** → Create new account (atau **Login** dengan test credentials)
3. **Get Current User** → Verify authentication works
4. **Get All Laundries** → Browse laundries
5. **Get Laundry Detail** → View laundry with services
   - Copy `laundry_id` dan `service_id` ke environment variables
6. **Create Order** → Create new order
   - Copy `order_id` dari response ke environment variable
7. **Get User Orders** → List all orders
8. **Get Order Detail** → View specific order
9. **Cancel Order** → Cancel order (jika status pending/confirmed)
10. **Update Order Status** → Update status (hanya untuk laundry owner)

## 💡 Tips

1. **Environment Variables:**
   - Setelah login, token otomatis tersimpan
   - Set `laundry_id` dan `service_id` manual setelah get laundry detail
   - Set `order_id` manual setelah create order

2. **Query Parameters:**
   - Edit query params langsung di URL bar
   - Atau gunakan Params tab di Postman

3. **Request Body:**
   - Semua request body sudah diisi dengan contoh data
   - Edit sesuai kebutuhan sebelum send

4. **Error Handling:**
   - Check response status code
   - Error response akan ada di field `error` atau `message`

## 🔄 Update Base URL

Jika server berjalan di URL berbeda:

1. Pilih environment **"Laundry Marketplace - Local"**
2. Edit variable `base_url`
3. Contoh: `http://localhost:3000` atau `https://api.example.com`

## 📦 Files Included

- `Laundry_Marketplace_API.postman_collection.json` - Main collection file
- `Laundry_Marketplace_API.postman_environment.json` - Environment variables

---

**Happy Testing! 🚀**

