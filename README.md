# Iivineri API

API pentru aplicația Iivineri, implementat în Go cu Fiber framework.

## 🚀 Caracteristici

- **API RESTful** cu Fiber framework
- **Documentație Swagger/OpenAPI** integrată
- **Autentificare JWT** cu suport 2FA
- **PostgreSQL** pentru baza de date
- **Redis** pentru cache
- **Docker** pentru deployment
- **Prometheus** pentru monitoring
- **Tool de management** pentru Docker

## 📋 Cerințe

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL (local sau Docker)
- Redis (opțional, pentru cache)

## 🛠️ Instalare și Setup

### 1. Clonează proiectul
```bash
git clone <repository-url>
cd iivineri-api
```

### 2. Configurează environment-ul
```bash
# Copiază și editează fișierul de configurare
cp .env.example .env
# Editează .env cu valorile tale
```

### 3. Instalează dependențele
```bash
go mod download
```

### 4. Generează documentația Swagger
```bash
make swagger-gen
```

## 🐳 Docker Setup

### Tool de Management Docker

Proiectul include un tool bash pentru managementul Docker care:
- ✅ **Detectează conflicte de porturi** automat
- ✅ **Găsește porturi disponibile** dacă există conflicte
- ✅ **Configurează environment-ul** automat
- ✅ **Gestionează serviciile** (start/stop/restart)
- ✅ **Profiles pentru diferite scenarii** (dev, monitoring)

### Comenzi rapide

```bash
# Configurare inițială (verifică porturi și creează .env)
./tool setup

# Pornește serviciile de bază (API + PostgreSQL + Redis)
./tool start

# Pornește cu tools de development (+ Adminer + Redis Commander)
./tool dev

# Pornește cu monitoring (+ Prometheus)
./tool monitoring

# Vezi statusul serviciilor
./tool status

# Vezi log-urile
./tool logs api -f

# Oprește toate serviciile
./tool stop
```

### Comenzi disponibile

```bash
./tool help  # Pentru lista completă de comenzi
```

## 📖 Makefile Commands

```bash
# Development local
make build              # Compilează aplicația
make run               # Rulează local
make swagger-gen       # Generează documentația Swagger

# Docker
make docker-setup      # Configurează Docker environment
make docker-start      # Pornește serviciile
make docker-dev        # Pornește cu development tools
make docker-stop       # Oprește serviciile

# Testing
make test              # Rulează testele
make test-coverage     # Testele cu coverage
make fmt               # Formatează codul
make lint              # Verifică codul

make help              # Lista completă de comenzi
```

## 🌐 Accesare Servicii

După pornirea cu `./tool dev`, ai acces la:

- **🖥️ API**: http://localhost:8080
- **📚 Swagger UI**: http://localhost:8080/swagger/index.html
- **🗄️ Adminer (DB)**: http://localhost:8081
- **🔴 Redis Commander**: http://localhost:8082
- **📈 Prometheus**: http://localhost:9090 (cu profil monitoring)

*Porturile se ajustează automat în caz de conflicte.*

## 🔧 Configurare

### Environment Variables

Principalele variabile de configurare în `.env`:

```env
# API
ENV=development
PORT=8080
LOG_LEVEL=info

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=iivineri_dev
DB_USERNAME=iivineri_user
DB_PASSWORD=your_password

# Cache (Redis)
CACHE_HOST=localhost
CACHE_PORT=6379

# Docker ports (configurate automat de tool)
API_EXTERNAL_PORT=8080
DB_EXTERNAL_PORT=5432
CACHE_EXTERNAL_PORT=6379
```

### Docker Profiles

- **Fără profile**: API + PostgreSQL + Redis
- **dev**: + Adminer + Redis Commander
- **tools**: Doar tools (Adminer + Redis Commander)
- **monitoring**: + Prometheus

## 📚 API Documentation

### Endpoints disponibile

- `GET /api/v1/` - Informații despre API
- `GET /api/v1/health` - Health check
- `GET /swagger/*` - Documentația Swagger

### Autentificare

```bash
# Înregistrare
POST /auth/register
{
  "nickname": "user123",
  "email": "user@example.com", 
  "password": "password123",
  "date_of_birth": "1990-01-01T00:00:00Z"
}

# Login
POST /auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Folosește token-ul primit în header:
Authorization: Bearer <token>
```

## 🛠️ Development Workflow

### 1. Setup inițial
```bash
./tool setup          # Configurează environment-ul
make swagger-gen       # Generează documentația
```

### 2. Dezvoltare
```bash
./tool dev            # Pornește serviciile cu tools
make dev              # Hot reload (dacă ai air instalat)
# sau
make run              # Run simplu
```

### 3. Testing
```bash
make test             # Rulează testele
make test-coverage    # Cu coverage
```

### 4. Database operations
```bash
./tool backup-db              # Backup
./tool restore-db backup.sql  # Restore
./tool shell postgres         # Deschide shell în PostgreSQL
```

## 🐛 Troubleshooting

### Port conflicts
Tool-ul detectează și rezolvă automat conflictele de porturi. Dacă întâmpini probleme:

```bash
./tool setup          # Re-configurează porturile
./tool config         # Vezi configurația curentă
```

### Docker issues
```bash
./tool cleanup        # Șterge tot (containers, images, volumes)
./tool build          # Rebuilds images
```

### Logs și debugging
```bash
./tool logs api -f           # Follow API logs
./tool logs postgres         # Database logs  
./tool shell api            # Shell în container API
```

## 📊 Monitoring

Cu profilul monitoring ai acces la:
- **Prometheus**: http://localhost:9090
- **Metrici API**: http://localhost:8080/metrics

```bash
./tool monitoring     # Pornește cu Prometheus
```

## 🔒 Securitate

- JWT authentication cu suport 2FA
- Password hashing cu bcrypt
- Rate limiting (planificat)
- Validare input cu validator
- CORS configurat pentru production

## 📝 Contribuție

1. Fork proiectul
2. Creează branch pentru feature (`git checkout -b feature/amazing-feature`)
3. Commit schimbările (`git commit -m 'Add amazing feature'`)
4. Push la branch (`git push origin feature/amazing-feature`)
5. Deschide Pull Request

## 📄 Licență

Acest proiect este licențiat sub [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0.html).