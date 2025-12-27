# Cool-Games API 🎮

A high-performance digital game marketplace backend built with **Go** and **PostgreSQL**. This project utilizes **Clean Architecture** (Hexagonal Architecture) to maintain a strict separation between business logic, delivery mechanisms, and data persistence.

## 🌟 Key Features

* **Atomic Transactions:** Buying a game triggers a multi-step database transaction (updates balance, reduces stock, records order, adds to library, and logs to ledger) with full rollback on failure.
* **Role-Based Access Control (RBAC):**
    * **Admin:** Full control over genres and global oversight.
    * **Publisher:** Can list games, restock inventory, and view detailed sales reports.
    * **Customer:** Can top up balances, purchase games, and view their digital library.
* **Inventory Tracking:** Automated stock level management with a historical log of every change (`game_quantity_history`).
* **Smart Filtering:** Search games by name (case-insensitive) and price range (min/max).
* **Financial Ledger:** A transparent record of all `credit` (top-ups) and `debit` (purchases) transactions.

## 🏗️ Project Structure

```text
├── config/             # Database connection logic
├── internal/
│   ├── auth/           # User management & login
│   ├── game/           # Game & inventory logic
│   ├── order/          # Transactions & Library
│   ├── genre/          # Category management
│   ├── domain/         # Shared interfaces & entities
│   └── middleware/     # JWT & Role-based security
├── main.go             # Entry point
└── .env                # Environment variables
```

## 🛠️ Tech Stack

* **Core:** Go (1.20+)
* **Framework:** Gin Gonic
* **Database:** PostgreSQL
* **Auth:** JWT (HS256)
* **Security:** Bcrypt (Password hashing)

## 🚦 API Endpoints

### Public / Auth

* `POST /register`: Create account (`admin`, `publisher`, or `customer`).
* `POST /login`: Receive JWT.

### Store (Protected)

* `GET /games`: Search & filter games.
* `GET /games/:id`: Get game details.
* `POST /games`: Create game (**Publisher**).
* `PATCH /games/:id/restock`: Update stock (**Publisher**).

### Orders & Finance (Protected)

* `POST /orders/topup`: Add balance (**Customer**).
* `POST /orders/buy`: Purchase game (**Customer**).
* `GET /orders/library`: View owned games (**Customer**).
* `GET /orders/sales-report`: View revenue analytics (**Publisher**).

## 🔧 Setup

1. **Configure `.env**`:
```env
DB_USER=your_user
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cool_games
JWT_SECRET=your_secret_key

```

2. **Migrate Database**: Run the provided SQL schema in your PostgreSQL instance.
3. **Run Server**:

```bash
go mod tidy
go run main.go
```