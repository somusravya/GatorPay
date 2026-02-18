# 🏃 Sprint 1 – Core Foundation & Authentication

Sprint Duration: Jan 6 – Jan 19, 2026  
Sprint Goal: Build authentication system, wallet management, and establish core frontend + backend architecture.

---

# 📌 1. Sprint Overview

Sprint 1 focused on building the foundational architecture of the GatorPay application.  
The goal was to create a secure, scalable, and production-ready system that supports authentication, wallet management, and a functional dashboard interface.

This sprint laid the groundwork for all future financial features such as bill payments, loans, peer-to-peer transfers, and analytics.

---

# 🎯 2. Sprint Objectives

The main objectives of Sprint 1 were:

- Implement secure Email/Password authentication
- Integrate Google OAuth login
- Implement JWT-based route protection
- Design database schema with auto-migration
- Build Wallet system with transactions
- Create responsive Angular frontend
- Establish clean backend–frontend communication
- Ensure atomic wallet operations using database transactions
- Deliver a functional dashboard with account overview

---

# 👥 3. User Stories

Below are the user stories planned for this sprint.

---

## 🔐 Authentication User Stories

1. As a new user, I want to register with email, password, and personal details so that I can create an account.

2. As a registered user, I want to login securely so that I can access my account dashboard.

3. As a user with a Google account, I want to authenticate using Google OAuth so I do not need to create a new password.

4. As a developer, I want to protect backend routes using JWT middleware so that only authenticated users can access secure APIs.

5. As an authenticated user, I want to retrieve my profile information so I can view my account details and wallet balance.

---

## 💰 Wallet User Stories

6. As a newly registered user, I want a wallet created automatically so that I can start using financial features immediately.

7. As a user, I want to add money to my wallet so that I can fund my account.

8. As a user, I want to withdraw money securely so that I can access my funds.

9. As a user, I want transactions to be recorded so that I can track my financial history.

---

## 🗄 Database & Backend Architecture Stories

10. As a developer, I want automatic database migration so that schema updates are applied during startup.

11. As a developer, I want well-defined models (User, Wallet, Transaction) so that the system is structured and scalable.

12. As a developer, I want soft delete support to prevent permanent data loss.

---

## 🎨 Frontend User Stories

13. As a user, I want a clean and responsive login page.

14. As a user, I want a registration page with validation.

15. As a developer, I want a centralized AuthService to manage authentication state.

16. As a developer, I want route guards to prevent unauthorized access.

17. As a developer, I want an HTTP interceptor to automatically attach JWT tokens.

18. As a user, I want a dashboard that shows my balance, KYC status, and recent transactions.

19. As a user, I want quick action buttons to easily access wallet features.

---

# 📋 4. Issues Planned in Sprint 1

Total Issues Planned: 19  
Total Story Points: 81

---

## Backend Issues

- User Registration API
- User Login API
- Google OAuth Integration
- JWT Authentication Middleware
- Get Current User API
- Database Setup & Auto-Migration
- User Model Implementation
- Wallet Model & Service
- Add Money API
- Withdraw Money API

---

## Frontend Issues

- Login Page UI
- Registration Page UI
- Authentication Service
- Route Guards
- HTTP Interceptor
- App Layout & Navigation
- Dashboard Component
- Quick Actions Grid
- Wallet Service

---

# ✅ 5. Successfully Completed Issues

All 19 planned issues were completed.

Completion Rate: 100%

---

# 🔐 6. Backend Implementation Details

### Authentication System
- Password hashing implemented using bcrypt
- JWT tokens generated with 7-day expiry
- Middleware validates token signature and expiration
- Protected routes require Bearer token
- Google OAuth integration implemented
- Auto wallet creation during registration

### Database
- SQLite used as primary database
- GORM ORM for model mapping
- Auto-migration on startup
- UUID primary keys
- Unique constraints on email, phone, username
- Soft delete support
- Seeded default billers and loan offers

### Wallet & Transactions
- Decimal precision using shopspring/decimal
- Atomic transactions using database transactions
- Add Money and Withdraw endpoints
- Sufficient balance validation
- Active wallet checks
- Transaction records created for every operation

---

# 🎨 7. Frontend Implementation Details

### Angular Architecture
- Standalone components
- Signal-based state management
- Centralized AuthService
- Persistent login via localStorage
- HTTP interceptor for automatic JWT injection

### UI/UX
- Premium dark theme
- Responsive layout (mobile + desktop)
- Sidebar navigation
- Header with user profile
- Logout functionality
- Animated quick action cards
- Dashboard with balance overview and transaction list

---

# 📊 8. Sprint Metrics

Total Issues: 19  
Total Story Points: 81  
Completed: 19/19  
Completion Rate: 100%

---

# 👨‍💻 9. Team Contribution

Backend Developer 1  
- Authentication APIs  
- Google OAuth  
- JWT Middleware  
- Get Current User API  
Story Points: 18  

Backend Developer 2  
- Database setup  
- Models  
- Wallet service  
- Add/Withdraw APIs  
Story Points: 23  

Frontend Developer 1  
- Login & Registration UI  
- AuthService  
- Route Guards  
- HTTP Interceptor  
Story Points: 21  

Frontend Developer 2  
- Layout & Navigation  
- Dashboard  
- Quick Actions  
- WalletService  
Story Points: 19  

---

# ❌ 10. Incomplete Issues

There were no incomplete issues in Sprint 1.

All planned features were delivered successfully within the sprint duration.

---

# 🚀 11. Sprint 1 Outcome

Sprint 1 successfully delivered:

- Secure authentication system (Email + Google OAuth)
- JWT-based authorization
- Wallet creation and transaction management
- Full backend–frontend integration
- Responsive Angular UI
- Dashboard with financial overview
- Production-ready foundation for future financial modules

This sprint established a strong architectural base for scaling the application in future sprints.

Future sprints will focus on:
- Bill payments
- Loan management
- Peer-to-peer transfers
- Advanced analytics
- Role-based access control
- Performance optimizations

---

# 🏁 Sprint 1 Status: Successfully Completed
