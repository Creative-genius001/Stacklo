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

10. Infrastructure
    Cloud Provider: AWS, Google Cloud Platform (GCP), or Azure. (e.g., AWS for this example).
    Compute:
    Kubernetes (EKS on AWS): For container orchestration, scaling, and managing your Golang microservices.
    Docker: For containerizing each Golang service.
    Database:
    Amazon RDS for PostgreSQL: A fully managed relational database service.
    Message Queue:
    Apache Kafka (or Amazon MSK): For asynchronous communication between services (e.g., for transaction processing, notifications, or event sourcing). This ensures loose coupling and resilience.
    Caching:
    Redis (or Amazon ElastiCache for Redis): For caching frequently accessed data (e.g., user sessions, crypto prices) to reduce database load and improve response times.
    Load Balancer:
    AWS ALB (Application Load Balancer): To distribute incoming traffic across your microservices.
    Monitoring & Logging:
    Prometheus/Grafana: For metrics collection and visualization.
    ELK Stack (Elasticsearch, Logstash, Kibana) or CloudWatch Logs: For centralized logging.
    CI/CD:
    GitHub Actions, GitLab CI, Jenkins, AWS CodePipeline: For automated testing, building, and deployment.

                           +----------------+
                           |     Client     |
                           +----------------+
                                   |
                             API Gateway (e.g. Kong, NGINX)
                                   |
           -----------------------------------------------------
           |           |               |           |            |

    User Service Wallet Service Crypto Service Txn Service Payment Service
    | | | | |
    PostgreSQL PostgreSQL PostgreSQL PostgreSQL PostgreSQL
    | | | | |
    Redis Redis Redis Redis Redis (caching)
    \_**\_\_\_\_** All services communicate via NATS **\_\_\_\_**/
    |
    Monitoring & Logging
    (Prometheus + Grafana, Loki, Jaeger)
