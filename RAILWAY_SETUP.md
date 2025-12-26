# 🚂 Railway Deployment Configuration

## ✅ Konfigurasi Otomatis

Repository ini sudah memiliki file konfigurasi Railway:
- `railway.json` - Konfigurasi Railway (build & start command)
- `Procfile` - Start command untuk Railway

**File-file ini memastikan hanya server yang dijalankan, BUKAN seed script.**

## ⚠️ Jika Masih Ada Masalah

Jika Railway masih menjalankan seed script (`/app/cmd/seed/main.go`), ikuti langkah berikut:

## 🔧 Cara Mengubah Konfigurasi di Railway

1. **Buka Railway Dashboard**
   - Login ke [railway.app](https://railway.app)
   - Pilih project Anda

2. **Cek Service Settings**
   - Klik pada service yang menjalankan aplikasi Go
   - Buka tab **Settings** atau **Variables**

3. **Periksa Build Command**
   - Cari **Build Command** atau **Build**
   - Pastikan tidak ada command yang menjalankan seed:
     ```
     ❌ SALAH: go run cmd/seed/main.go
     ❌ SALAH: go build ./cmd/seed && ./bin/seed
     ✅ BENAR: (kosong atau tidak ada)
     ```

4. **Periksa Start Command**
   - Cari **Start Command** atau **Start**
   - Pastikan hanya menjalankan server:
     ```
     ❌ SALAH: go run cmd/seed/main.go && go run cmd/server/main.go
     ❌ SALAH: ./bin/seed && ./bin/server
     ✅ BENAR: go run cmd/server/main.go
     ✅ BENAR: ./bin/server (jika sudah di-build)
     ```

5. **Periksa Deploy Command**
   - Beberapa Railway setup menggunakan **Deploy Command**
   - Pastikan tidak ada seed di sini juga

## 📝 Konfigurasi yang Benar

### Untuk Development/Testing:
- **Build Command**: (kosong atau `go build -o bin/server ./cmd/server`)
- **Start Command**: `go run cmd/server/main.go`

### Untuk Production:
- **Build Command**: `go build -o bin/server ./cmd/server`
- **Start Command**: `./bin/server`

## 🔍 Cara Cek Apakah Seed Masih Dijalankan

Jika Anda melihat log seperti ini di Railway:
```
/app/cmd/seed/main.go: ...
```

Berarti seed script masih dijalankan. Hapus command tersebut dari Railway settings.

## ✅ Setelah Diubah

1. **Redeploy** service di Railway
2. **Cek logs** untuk memastikan hanya server yang berjalan
3. **Pastikan** tidak ada error tentang seed script

## 🎯 Catatan

- Seed script sudah **dihapus** dari repository
- Data sudah **di-import** ke database
- Tidak perlu menjalankan seed lagi di production

