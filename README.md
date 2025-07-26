Stacklo

This is a financial payment system that allows users buy crypto currecnies online using fiat currency and store safely in a wallet internally.

A microservices architecture was used for scalability, resilience, and maintainability.

1.  Overall Architecture: Microservices

2.  Core Services (Golang)
    Each service will be a separate Golang application, communicating via internal APIs (e.g., gRPC for high performance or REST).

3.  Authentication & User Service:
    Purpose: Handles user registration, login, session management, user profiles (e.g., KYC status), and password management.
    Features: User signup/login, password hashing, JWT token generation/validation.

4.  Wallet Service:
    Purpose: Manages user virtual wallets (NGN) and all fiat currency transactions.
    Features: Wallet creation on signup, NGN deposit/withdrawal, balance tracking, transaction history.

5.  Crypto Service:
    Purpose: Manages user's cryptocurrency holdings (BTC, ETH) and handles crypto buy/sell operations.
    Features: Crypto wallet generation (if self-custodial, otherwise just balance tracking), buy/sell orders, crypto balance updates.

6.  Transaction Service:
    Purpose: Orchestrates and records all financial transactions (fiat and crypto). This service ensures atomicity and consistency across wallet and crypto balances.
    Features: Transaction logging, status updates, idempotency checks.

7.  Notification Service:
    Purpose: Sends real-time notifications to users for transactions, login alerts, etc.
    Features: Email, SMS, push notifications.

8.  Payment Service: This focuses only on interacting with external payment providers (Paystack, crypto exchanges). This includes handling their specific API contracts, authentication, webhooks, and potential rate limits.

9.  API Design
    External APIs (for Frontend): RESTful APIs (JSON over HTTP) are standard for client-server communication.
    Internal APIs (between services): gRPC for efficient, high-performance inter-service communication using Protocol Buffers.

                           +----------------+
                           |     Client     |
                           +----------------+
                                   |
                             API Gateway
                                   |
           -----------------------------------------------------
           |           |                 |         |

    User Service Wallet Service Txn Service Payment Service
    | | | | |
    PostgreSQL PostgreSQL PostgreSQL PostgreSQL PostgreSQL
    | | | | |
    Redis Redis Redis Redis Redis (caching)
    \_**\_\_\_\_** All services communicate via NATS **\_\_\_\_**/
    |
    Monitoring & Logging
    (Prometheus + Grafana, Loki, Jaeger)
