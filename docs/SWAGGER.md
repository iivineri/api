# Swagger Documentation Guide

Această documentare explică cum funcționează implementarea Swagger/OpenAPI în proiectul Iivineri API și cum să adaugi documentația pentru endpoint-uri noi.

## 📋 Cuprins

1. [Structura actuală](#structura-actuală)
2. [Cum funcționează](#cum-funcționează)
3. [Adăugarea documentației pentru endpoint-uri noi](#adăugarea-documentației-pentru-endpoint-uri-noi)
4. [Tipuri de adnotări](#tipuri-de-adnotări)
5. [Exemple practice](#exemple-practice)
6. [Generarea documentației](#generarea-documentației)
7. [Debugging și rezolvarea problemelor](#debugging-și-rezolvarea-problemelor)

## 📁 Structura actuală

```
swag init --generalInfo main.go --dir ./ --output ./pkg/swagger --parseDependency --parseInternal
```

```
iivineri-api/
├── main.go                      # Configurația generală Swagger
├── docs/
│   ├── docs.go                  # Generat automat - NU EDITA
│   ├── swagger.json             # Generat automat - NU EDITA
│   ├── swagger.yaml             # Generat automat - NU EDITA
│   └── SWAGGER.md              # Această documentare
├── cmd/serve/serve.go           # Configurare endpoint Swagger UI
└── internal/fiber/modules/auth/handler/handler.go  # Exemple de endpoint-uri documentate
```

## ⚙️ Cum funcționează

### 1. Configurația globală (main.go)

```go
// @title Iivineri API
// @version 1.0
// @description API pentru aplicația Iivineri
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

**Explicația adnotărilor:**
- `@title` - Numele API-ului din documentația Swagger
- `@version` - Versiunea API-ului
- `@description` - Descrierea generală a API-ului
- `@host` - Hostname-ul unde rulează API-ul
- `@BasePath` - Path-ul de bază pentru toate endpoint-urile
- `@schemes` - Protocoalele suportate (http, https)
- `@securityDefinitions.apikey` - Definește autentificarea Bearer JWT

### 2. Configurarea endpoint-ului Swagger UI (cmd/serve/serve.go)

```go
import fiberSwagger "github.com/swaggo/fiber-swagger"

// În funcția de configurare a aplicației:
app.Get("/swagger/*", fiberSwagger.WrapHandler)
```

Acest cod face ca documentația Swagger să fie disponibilă la: `http://localhost:8080/swagger/index.html`

### 3. Importul documentației generate

```go
import _ "iivineri/docs"  // Import pentru side-effects
```

## 🚀 Adăugarea documentației pentru endpoint-uri noi

### Pasul 1: Adaugă adnotările Swagger în handler

Pentru fiecare endpoint nou, adaugă comentarii speciale **direct deasupra funcției handler**:

```go
// @Summary Scurtă descriere a endpoint-ului
// @Description Descriere detaliată a funcționalității
// @Tags categoria-endpoint-ului
// @Accept json
// @Produce json
// @Param nume-parametru tip-locație tip-date required "Descriere parametru"
// @Security BearerAuth
// @Success 200 {object} TipulRaspunsului "Descriere succes"
// @Failure 400 {object} map[string]interface{} "Descriere eroare"
// @Router /calea/endpoint-ului [metoda-http]
func NumeleFunctiei(c *fiber.Ctx) error {
    // Implementarea endpoint-ului
}
```

### Pasul 2: Definește modelele de date

Pentru request/response bodies, creează structuri în `internal/fiber/modules/{modul}/models/`:

```go
package models

// RegisterRequest reprezintă datele necesare pentru înregistrarea unui utilizator
type RegisterRequest struct {
    Nickname    string    `json:"nickname" validate:"required,min=3,max=32"`
    Email       string    `json:"email" validate:"required,email"`
    Password    string    `json:"password" validate:"required,min=8"`
    DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}

// RegisterResponse reprezintă răspunsul pentru înregistrarea cu succes
type RegisterResponse struct {
    Message string `json:"message"`
    UserID  uint   `json:"user_id"`
}
```

### Pasul 3: Regenerează documentația

```bash
# Instalează swag CLI dacă nu e deja instalat
go install github.com/swaggo/swag/cmd/swag@latest

# Generează documentația
swag init
```

## 📝 Tipuri de adnotări

### Adnotări de bază

| Adnotare | Descriere | Exemplu |
|----------|-----------|---------|
| `@Summary` | Titlul scurt al endpoint-ului | `@Summary Create new user` |
| `@Description` | Descrierea detaliată | `@Description Creates a new user account with email verification` |
| `@Tags` | Gruparea endpoint-urilor | `@Tags users, auth` |
| `@Accept` | Tipul de conținut acceptat | `@Accept json` |
| `@Produce` | Tipul de conținut returnat | `@Produce json` |
| `@Router` | Calea și metoda HTTP | `@Router /users [post]` |

### Parametri

| Tip parametru | Sintaxă | Exemplu |
|--------------|---------|---------|
| Path parameter | `@Param nume path tip required "descriere"` | `@Param id path int true "User ID"` |
| Query parameter | `@Param nume query tip required "descriere"` | `@Param page query int false "Page number"` |
| Body parameter | `@Param nume body tip required "descriere"` | `@Param user body models.CreateUser true "User data"` |
| Header parameter | `@Param nume header tip required "descriere"` | `@Param Authorization header string true "Bearer token"` |

### Răspunsuri

```go
// @Success cod-status {tip} tip-obiect "descriere"
@Success 200 {object} models.User "User retrieved successfully"
@Success 201 {object} models.CreateUserResponse "User created"
@Success 204 "No content"

// @Failure cod-status {tip} tip-obiect "descriere"
@Failure 400 {object} map[string]interface{} "Bad request"
@Failure 401 {object} map[string]interface{} "Unauthorized"
@Failure 404 {object} map[string]interface{} "User not found"
@Failure 500 {object} map[string]interface{} "Internal server error"
```

### Securitate

```go
// Pentru endpoint-uri care necesită autentificare
@Security BearerAuth

// Pentru endpoint-uri publice - nu specifica nimic
```

## 💡 Exemple practice

### Exemplu 1: GET endpoint simplu

```go
// @Summary Get user by ID
// @Description Retrieve a specific user by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} models.User "User found"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
    // Implementarea
}
```

### Exemplu 2: POST endpoint cu body

```go
// @Summary Create new product
// @Description Create a new product in the system
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Product data"
// @Security BearerAuth
// @Success 201 {object} models.Product "Product created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid product data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "Product already exists"
// @Router /products [post]
func CreateProduct(c *fiber.Ctx) error {
    // Implementarea
}
```

### Exemplu 3: GET endpoint cu query parameters

```go
// @Summary List all users
// @Description Get a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Security BearerAuth
// @Success 200 {object} models.UserListResponse "Users retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users [get]
func ListUsers(c *fiber.Ctx) error {
    // Implementarea
}
```

### Exemplu 4: PUT endpoint pentru update

```go
// @Summary Update user profile
// @Description Update the profile information of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "Updated user data"
// @Security BearerAuth
// @Success 200 {object} models.User "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - cannot update other users"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
    // Implementarea
}
```

### Exemplu 5: DELETE endpoint

```go
// @Summary Delete user
// @Description Permanently delete a user account
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 204 "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - insufficient permissions"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
    // Implementarea
}
```

## ⚙️ Generarea documentației

### Comenzi esențiale

```bash
# Installează tool-ul swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generează documentația (din directorul root al proiectului)
swag init

# Verifică versiunea
swag --version

# Help pentru mai multe opțiuni
swag --help
```

### Comanda cu opțiuni avansate

```bash
# Generează cu opțiuni specifice
swag init --generalInfo main.go --dir ./ --output ./docs --parseDependency --parseInternal
```

**Opțiuni explicate:**
- `--generalInfo main.go` - Fișierul cu informațiile generale
- `--dir ./` - Directorul sursă de scanat
- `--output ./docs` - Directorul de output
- `--parseDependency` - Parsează dependențele externe
- `--parseInternal` - Parsează pachetele interne

### Integrare în workflow

Pentru a automatiza generarea, poți crea un script sau adăuga în Makefile:

```bash
# În ./.tool sau similar
swagger-gen() {
    echo "🔄 Generating Swagger documentation..."
    swag init
    if [ $? -eq 0 ]; then
        echo "✅ Swagger documentation generated successfully"
        echo "📚 Available at: http://localhost:8080/swagger/index.html"
    else
        echo "❌ Failed to generate Swagger documentation"
        exit 1
    fi
}
```

## 🐛 Debugging și rezolvarea problemelor

### Probleme comune

#### 1. **Swagger UI nu se încarcă**

**Soluție:**
```bash
# Verifică dacă documentația e generată
ls -la docs/

# Regenerează documentația
swag init

# Verifică importul în main.go
grep "docs" main.go
```

#### 2. **Endpoint-urile nu apar în documentație**

**Cauze posibile:**
- Comentariile Swagger nu sunt direct deasupra funcției
- Sintaxa comentariilor e incorectă
- Funcția nu e exportată (nu începe cu majusculă)
- Pachetul nu e importat

**Soluție:**
```bash
# Regenerează cu verbose pentru debugging
swag init --verbose

# Verifică sintaxa comentariilor
# Exemple corecte:
// @Summary Descriere scurtă
// @Router /path [method]
```

#### 3. **Modele/structuri nu apar**

**Soluție:**
```go
// Asigură-te că structurile sunt exportate
type User struct {  // ✅ Corect - începe cu majusculă
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

type user struct {  // ❌ Greșit - nu e exportată
    ID   uint   `json:"id"`
    Name string `json:"name"`
}
```

#### 4. **Informațiile generale nu se actualizează**

**Soluție:**
```bash
# Șterge fișierele generate și regenerează
rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml
swag init
```

### Verificarea documentației generate

```bash
# Verifică fișierul generat
cat docs/swagger.json | jq .

# Verifică dacă toate endpoint-urile sunt incluse
grep -r "@Router" internal/ | wc -l
grep "paths" docs/swagger.json
```

### Testarea în browser

1. **Pornește aplicația:**
   ```bash
   go run main.go serve
   ```

2. **Accesează Swagger UI:**
   - URL: http://localhost:8080/swagger/index.html
   - Verifică că toate endpoint-urile apar
   - Testează funcționalitatea "Try it out"

3. **Verifică JSON-ul raw:**
   - URL: http://localhost:8080/swagger/doc.json

## 📚 Referințe și resurse utile

- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [swaggo/swag Documentation](https://github.com/swaggo/swag)
- [Fiber Swagger Middleware](https://github.com/swaggo/fiber-swagger)
- [Go Struct Tags pentru JSON](https://golang.org/pkg/encoding/json/#Marshal)

## 🔄 Workflow complet pentru un endpoint nou

1. **Creează handler-ul cu documentația Swagger**
2. **Definește modelele de request/response**
3. **Adaugă validările necesare**
4. **Regenerează documentația: `swag init`**
5. **Testează endpoint-ul în Swagger UI**
6. **Verifică că documentația e corectă și completă**

## ⚡ Tips și best practices

- **Folosește tag-uri consistente** pentru gruparea endpoint-urilor
- **Documentează toate răspunsurile posibile** (success + error cases)
- **Folosește modele pentru request/response** în loc de `map[string]interface{}`
- **Specifică validation rules în struct tags**
- **Regenerează documentația după fiecare modificare**
- **Testează endpoint-urile în Swagger UI** înainte de commit
- **Menține descrierile clare și concise**
