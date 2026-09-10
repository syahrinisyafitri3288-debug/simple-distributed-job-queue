# Submission Notes & Implementation Details

File ini berisi rangkuman fitur yang sudah diimplementasikan, cara menjalankan, serta catatan teknis pengujian project.

---

## 🚀 Cara Menjalankan

1. Generate schema GraphQL:
   ```bash
   go generate ./delivery/graphql/schema

```

2. Jalankan aplikasi:
```bash
go run main.go

```


3. Akses antarmuka:
* **GraphQL Playground**: `http://localhost:58579/graphiql`
* **HTMX Dashboard**: `http://localhost:58579/jobqueue/dashboard`



---

## Fitur yang Sudah Selesai

* **GraphQL Backend**: Query/Mutation `Enqueue`, `GetAllJobs`, `GetJobById`, dan `GetAllJobStatus` berjalan lancar. `SimultaneousCreateJob` dan `SimulateUnstableJob` dites via alias mutation.
* **Retry Mechanism**: Retry maksimal 3x dengan delay 500ms. Khusus `unstable-job` dibuat gagal 2x lalu sukses di percobaan ke-3.
* **Concurrency & Load Test**: Thread-safe menggunakan Mutex. Sudah dibuatkan unit test untuk memproses 100 job bersamaan di `service/job_load_test.go` (Status: **PASS**).
```bash
go test -v -run TestConcurrentLoad ./service/...

```


* **HTMX Dashboard**: Fitur form create job, tombol unstable job, serta status summary & tabel job auto-refresh tiap 2 detik. Tombol *View* berfungsi menampilkan detail job.

---

## Notes

* **Idempotency**: Saat ini dijamin lewat unik ID per job. Fitur *idempotency key* dari client belum diimplementasikan karena keterbatasan waktu (rencana ke depan: tambah param `idempotencyKey` di schema dan validasi di repository sebelum save).
* **Race Detector (`-race`)**: Tidak dijalankan secara otomatis di lokal karena keterbatasan environment Windows (butuh CGO/GCC). Namun, pengujian manual (spam click) dan load test 100 job terbukti aman tanpa crash atau corruption data.

```