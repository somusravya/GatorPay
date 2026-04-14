# GatorPay - Digital Payment & Trading Platform

## Team Members

**Frontend Development:**
- Somu Geetha Sravya
- Shivankita K

**Backend Development:**
- Tharun Kamsala
- Kaushik Ramesh

## Project Overview

GatorPay is a digital fintech application that combines everyday banking, payments, and stock trading into one platform. Users can manage their money, send payments to friends, pay bills, and invest in stocks—all from a single app.

## What We're Building

We're creating a web application with three main parts:
1. **Frontend** - A user-friendly web interface built with Angular
2. **Banking Backend** - Handles money transfers, bills, and accounts
3. **Trading Backend** - Manages stock market data and investments

## Key Features We'll Implement

### 1. User Accounts & Security
- Email and password registration
- Google sign-in option
- Secure login with authentication tokens
- Personal profile management

### 2. Digital Wallet
- Check account balance
- Add money from bank or card
- Withdraw money to bank account
- View transaction history

### 3. Send Money to Friends (P2P)
- Send money to other users by username, email, or phone
- Add notes to payments
- See recent contacts
- Search for users on the platform

### 4. Bill Payments
- Pay utility bills (electricity, water, gas, internet)
- Mobile and DTH recharges
- Save frequently used billers
- Track payment history

### 5. Loans
- Browse available loan offers
- Check credit score
- Apply for loans
- Pay monthly EMIs (installments)

### 6. Virtual Cards
- Create virtual debit/credit cards
- Set spending limits
- Freeze/unfreeze cards
- View card details securely

### 7. QR Code Payments
- Generate QR codes for receiving payments
- Scan QR codes to pay merchants
- Business dashboard for merchants
- Track sales and customers

### 8. Rewards Program
- Earn cashback on transactions
- Track reward points
- View reward history

### 9. Stock Trading
- Search for stocks
- View live stock prices and charts
- Buy and sell stocks
- Track your investment portfolio
- Create watchlists of favorite stocks
- Set price alerts

### 10. Download Statements
- Generate account statements
- Download in CSV or PDF format
- Select custom date ranges

## How It Works

### For Regular Users:
1. **Sign Up** - Create an account using email or Google
2. **Add Money** - Load your wallet from your bank account
3. **Make Payments** - Send money to friends or pay bills
4. **Earn Rewards** - Get cashback on transactions
5. **Invest** - Trade stocks with your available balance
6. **Track Everything** - View all transactions and investments in one place

### For Merchants:
1. **Register Business** - Sign up as a merchant
2. **Get QR Code** - Receive a unique QR code for your business
3. **Accept Payments** - Customers scan and pay instantly
4. **View Dashboard** - Track sales, revenue, and customer data

### For Investors:
1. **Open Trading Account** - Complete verification process
2. **Deposit Funds** - Transfer money from wallet to trading account
3. **Research Stocks** - Search stocks, view charts and news
4. **Trade** - Buy and sell stocks at market prices
5. **Monitor Portfolio** - Track profits, losses, and performance

## Technology Stack

### Frontend
- **Angular 19** - For building the web interface
- **TypeScript** - For writing clean, type-safe code
- **SCSS** - For styling the application

### Backend
- **Go (Golang)** - For building fast, reliable servers
- **Gin Framework** - For handling web requests
- **PostgreSQL Database** - For storing all data
- **JWT Tokens** - For secure authentication

### External Services
- **RapidAPI** - For real-time stock market data
- **Google OAuth** - For Google login

## Project Architecture

```
User's Browser
      ↓
Angular Frontend (Port 4200)
      ↓
      ├─→ Banking API (Port 8080)
      │   └─→ Database (gatorpay.db)
      │
      └─→ Trading API (Port 8081)
          └─→ Database (trading.db)
```

### Services Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        GatorPay Ecosystem                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────────┐    ┌─────────────────┐                    │
│  │ gatorpay-frontend│    │                 │                    │
│  │  (Angular 19)    │◄──►│   Web Browser   │                    │
│  │  Port: 4200      │    │                 │                    │
│  └────────┬─────────┘    └─────────────────┘                    │
│           │                                                     │
│           ▼                                                     │
│  ┌─────────────────────────────────────────────────────┐        │
│  │                   API Gateway Layer                 │        │
│  └─────────────────────────────────────────────────────┘        │
│           │                           │                         │
│           ▼                           ▼                         │
│  ┌─────────────────┐         ┌─────────────────┐                │
│  │ gatorpay-backend│         │ gatorpay-trading│                │
│  │ (Main Banking)  │         │ (Stock Trading) │                │
│  │ Port: 8080      │         │ Port: 8081      │                │
│  └────────┬────────┘         └────────┬────────┘                │
│           │                           │                         │
│           ▼                           ▼                         │
│  ┌─────────────────┐         ┌─────────────────┐                │
│  │  gatorpay.db    │         │   trading.db    │                │
│  │  (PostgreSQL)   │         │  (PostgreSQL)   │                │
│  └─────────────────┘         └─────────────────┘                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Database Structure

We'll store:
- **Users** - Account information, passwords (encrypted)
- **Wallets** - Balance and transaction records
- **Transactions** - All money movements
- **Bills** - Saved billers and payment history
- **Loans** - Loan applications and EMI records
- **Cards** - Virtual card details
- **Merchants** - Business accounts and QR codes
- **Trades** - Stock buy/sell records
- **Portfolios** - Investment holdings
- **Watchlists** - Saved stocks for tracking

## Security Features

- Passwords are encrypted (never stored in plain text)
- All sensitive actions require authentication
- Card numbers are masked (show only last 4 digits)
- OTP verification for viewing full card details
- Secure token-based login system
- Protection against unauthorized access

## Development Plan

### Phase 1: Basic Setup
- Set up Angular frontend
- Create Go backend servers
- Set up databases
- Implement user registration and login

### Phase 2: Core Banking
- Build wallet functionality
- Implement P2P transfers
- Add bill payment system
- Create transaction history

### Phase 3: Advanced Features
- Develop loan system
- Build virtual card management
- Implement merchant QR payments
- Add rewards program

### Phase 4: Trading Platform
- Integrate stock market API
- Build trading account system
- Implement buy/sell functionality
- Create portfolio tracking
- Add watchlists and alerts

### Phase 5: Polish & Testing
- Generate statements
- Add security features
- Test all functionality
- Fix bugs and improve UI

## Expected Outcomes

By the end of this project, we'll have:
- A fully functional digital payment platform
- Secure user authentication system
- Working wallet with add/withdraw money
- P2P money transfer capability
- Bill payment integration
- Loan management system
- Virtual card generation
- Merchant payment solution
- Real-time stock trading
- Portfolio management
- Transaction statements

## How to Run the Project

1. **Start Backend Services:**
   - Run the banking server on port 8080
   - Run the trading server on port 8081

2. **Start Frontend:**
   - Run the Angular app on port 4200

3. **Access the App:**
   - Open browser and go to http://localhost:4200
   - Register a new account or login
   - Start using all features!

## Target Users

- **Individuals** - People who want to manage money and make payments
- **Investors** - People interested in stock trading
- **Merchants** - Businesses that want to accept digital payments
- **Friends & Family** - Groups who regularly send money to each other

## Why This Project?

This project demonstrates:
- Full-stack development skills (Frontend + Backend)
- Database design and management
- API integration with external services
- Security best practices
- Real-world application development
- Microservices architecture
- Financial transaction handling
- Modern web development technologies

## Success Criteria

The project will be successful when:
- Users can register and login securely
- Money can be transferred between users
- Bills can be paid through the platform
- Stocks can be bought and sold
- All transactions are recorded accurately
- The application is secure and reliable
- The user interface is intuitive and easy to use

---

**Project Name:** GatorPay  
**Project Type:** Full-Stack Web Application  
**Domain:** Financial Technology (Fintech)  
**Duration:** Academic Project  
**Team Size:** 4 Members  
**Status:** In Development
