# Mini PaaS - Golang Backend & React Frontend

Một nền tảng PaaS (Platform as a Service) mini cho phép triển khai và quản lý ứng dụng web thông qua giao diện web.

## 🚀 Tính năng

- **Quản lý ứng dụng**: Tạo, cập nhật, xóa ứng dụng
- **Triển khai tự động**: Tự động build và deploy ứng dụng
- **Monitoring**: Theo dõi trạng thái ứng dụng
- **Container Management**: Quản lý Docker containers
- **Kubernetes Integration**: Orchestration với K8s

## 🛠️ Tech Stack

### Backend
- **Golang** với Gin framework
- **PostgreSQL** database
- **Docker** containerization
- **Kubernetes** orchestration

### Frontend
- **React** với TypeScript
- **Tailwind CSS** styling
- **Axios** HTTP client

### DevOps
- **Docker Compose** cho development
- **GitLab CI/CD** (có thể mở rộng)
- **Minikube/Kind** cho local K8s

## 📁 Cấu trúc dự án

```
mini-paas-golang/
├── backend/                 # Golang backend
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   ├── pkg/                # Public libraries
│   ├── api/                # API definitions
│   ├── configs/            # Configuration files
│   └── Dockerfile
├── frontend/               # React frontend
│   ├── src/
│   ├── public/
│   └── Dockerfile
├── k8s/                    # Kubernetes manifests
├── docker-compose.yml      # Development environment
└── README.md
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL
- Minikube (optional)

### Development Setup
```bash
# Clone repository
git clone <repository-url>
cd mini-paas-golang

# Start development environment
docker-compose up -d

# Backend (in separate terminal)
cd backend
go mod download
go run cmd/server/main.go

# Frontend (in separate terminal)
cd frontend
npm install
npm start
```

## 📚 API Documentation

API documentation sẽ được tạo tự động khi chạy backend.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

MIT License