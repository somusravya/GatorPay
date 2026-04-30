# GatorPay Sprint 4

GatorPay is a full-stack digital wallet and financial services application. Sprint 4 expands the product from core payments into a broader financial operating system with insights, budgeting, subscriptions, fraud review, notifications, social payments, merchant tools, admin analytics, and a trading microservice.

## Project Structure

```text
.
|-- backend/            # Go + Gin API for GatorPay wallet, payments, Sprint 4 modules
|-- frontend/           # Angular 19 web application
|-- flowpay-trading/    # Go + Gin trading microservice with SQLite persistence
|-- Sprint4.md          # Sprint 4 summary notes
|-- commands.txt        # Quick command reference
`-- README.md           # Main project guide
```

## Tech Stack

| Layer | Technology |
| --- | --- |
| Frontend | Angular 19, TypeScript, SCSS, RxJS, Cypress, Jasmine/Karma |
| Backend | Go, Gin, GORM, PostgreSQL, JWT authentication |
| Trading service | Go, Gin, GORM, SQLite |
| Testing | Go test, Jasmine/Karma, Cypress E2E |

## Sprint 4 Feature Summary

### Backend Features

- AI financial insights summary endpoint for spending analysis and financial health reporting.
- Budget goals and autosave rules for goal-based saving and roundup automation.
- Subscription tracking for recurring payments and autopay preferences.
- Fraud alert APIs for suspicious activity review.
- Notification APIs for in-app messages, read state, and user notification preferences.
- Social payment APIs for feed items, reactions, friends, and friend requests.
- Merchant invoicing and payment link APIs.
- Admin analytics APIs for platform metrics, users, and fraud review.
- Rate limiting middleware added globally.
- Database migrations for all Sprint 4 models.

### Frontend Features

- Dashboard with wallet overview, quick actions, and recent transactions.
- Wallet, add money, withdrawal, transaction history, rewards, and offers.
- Transfer flow with recipient search, contacts, and balance display.
- Bills workflow with categories, billers, saved billers, and bill payment.
- Loans page for offers, eligibility, applications, repayments, and cancellation.
- Cards page for virtual card creation, card details, OTP, and freeze controls.
- Trading page connected to the trading microservice for market data and portfolio actions.
- QR scan and payment flow.
- Sprint 4 screens for insights, budget planner, subscriptions, social feed, notifications, merchant portal, and admin console.
- Auth-protected routing with guest-only login/register screens.
- HTTP auth interceptor for JWT-based API requests.

### Trading Microservice Features

- Independent trading API running on port `8081`.
- SQLite database stored in `flowpay-trading/trading.db`.
- Market search, quote, details, chart, and market summary endpoints.
- Trading account verification, portfolio, order history, and trade execution endpoints.

## Prerequisites

- Go installed.
- Node.js and npm installed.
- PostgreSQL running for the main backend.
- Angular CLI is provided through the frontend npm dependencies and can be run through `npm start` or `npx ng`.

## Environment Variables

Create `backend/.env` or export these variables before running the backend:

```bash
PORT=8080
DB_DSN=host=localhost user=postgres password=postgres dbname=gatorpay port=5432 sslmode=disable
JWT_SECRET=gatorpay-super-secret-key-change-in-production
CORS_ORIGINS=http://localhost:4200
FRONTEND_URL=http://localhost:4200
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASS=
SMTP_FROM=noreply@gatorpay.app
```

`DB_DSN` is required. The backend exits at startup if it is missing.

The frontend API URL is configured in `frontend/src/environments/environment.ts`:

```ts
apiUrl: 'http://localhost:8080/api/v1'
```

## Install Dependencies

### Backend

```bash
cd backend
go mod download
```

### Trading Microservice

```bash
cd flowpay-trading
go mod download
```

### Frontend

```bash
cd frontend
npm install
```

## Quick Commands

Use these commands during the demo or while testing locally.

### Run Commands

```bash
cd backend
go run main.go

cd frontend
npm start

cd flowpay-trading
go run cmd/api/main.go
```

### Backend Tests

```bash
cd backend
go test ./...

cd flowpay-trading
go test ./...
```

### Frontend Tests

```bash
cd frontend
npm run test

cd frontend
npm run cypress:open

cd frontend
npm run cypress:run
```

## Run the Application

Open three terminals and run each service separately.

### 1. Start the Main Backend

```bash
cd backend
go run main.go
```

Default backend URL:

```text
http://localhost:8080/api/v1
```

### 2. Start the Trading Microservice

```bash
cd flowpay-trading
go run cmd/api/main.go
```

Default trading URL:

```text
http://localhost:8081/api/v1
```

Health check:

```bash
curl http://localhost:8081/api/v1/health
```

### 3. Start the Frontend

```bash
cd frontend
npm start
```

Default frontend URL:

```text
http://localhost:4200
```

## Test Commands

### Backend Tests

```bash
cd backend
go test ./...
```

Verbose mode:

```bash
cd backend
go test ./... -v
```

### Trading Microservice Tests

```bash
cd flowpay-trading
go test ./...
```

Verbose mode:

```bash
cd flowpay-trading
go test ./... -v
```

### Frontend Unit Tests

```bash
cd frontend
npm run test
```

Single-run mode for CI-style execution:

```bash
cd frontend
npm run test -- --watch=false --browsers=ChromeHeadless
```

### Frontend Cypress E2E Tests

Start the frontend first:

```bash
cd frontend
npm start
```

Then run Cypress in another terminal:

```bash
cd frontend
npm run cypress:run
```

Interactive Cypress runner:

```bash
cd frontend
npm run cypress:open
```

## API Documentation

All protected backend routes require:

```text
Authorization: Bearer <jwt_token>
```

### Auth

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login and start OTP flow |
| POST | `/api/v1/auth/verify-otp` | Verify OTP and receive auth token |
| POST | `/api/v1/auth/resend-otp` | Resend OTP |
| POST | `/api/v1/auth/google` | Google auth request |
| GET | `/api/v1/auth/me` | Get current authenticated user |

### Wallet

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/wallet/add` | Add money to wallet |
| POST | `/api/v1/wallet/withdraw` | Withdraw money |
| GET | `/api/v1/wallet/transactions` | List wallet transactions |
| GET | `/api/v1/wallet/statement` | Get wallet statement |

### Transfers

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/transfer/send` | Send money to another user |
| GET | `/api/v1/transfer/contacts` | Get recent transfer contacts |
| GET | `/api/v1/transfer/search` | Search users for transfer |

### Bills

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/bills/categories` | List bill categories |
| GET | `/api/v1/bills/billers` | List billers, optionally by category |
| POST | `/api/v1/bills/pay` | Pay a bill |
| GET | `/api/v1/bills/saved` | List saved billers |
| DELETE | `/api/v1/bills/saved/:id` | Remove a saved biller |

### Rewards

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/rewards` | Rewards summary |
| GET | `/api/v1/rewards/history` | Rewards history |
| GET | `/api/v1/rewards/offers` | Active reward offers |

### Loans

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/loans/offers` | List loan offers |
| GET | `/api/v1/loans/eligibility` | Check loan eligibility |
| POST | `/api/v1/loans/apply` | Apply for a loan |
| GET | `/api/v1/loans` | List user loans |
| GET | `/api/v1/loans/:id` | Get loan details |
| POST | `/api/v1/loans/:id/pay` | Pay loan EMI |
| POST | `/api/v1/loans/:id/cancel` | Cancel a loan |

### Cards

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/cards` | Create virtual card |
| GET | `/api/v1/cards` | List virtual cards |
| GET | `/api/v1/cards/:id` | Get card details |
| POST | `/api/v1/cards/:id/otp` | Request card OTP |
| POST | `/api/v1/cards/:id/details` | Reveal card details after OTP flow |
| POST | `/api/v1/cards/:id/freeze` | Freeze or unfreeze a card |

### QR and Merchant

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/merchant/register` | Register merchant profile |
| POST | `/api/v1/qr/generate` | Generate merchant QR |
| POST | `/api/v1/qr/lookup` | Look up QR payment data |
| POST | `/api/v1/qr/pay` | Pay through QR |

### Sprint 4: Insights

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/insights/summary` | Get financial insight summary |

### Sprint 4: Budget and Autosave

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/budget/goals` | Create budget or savings goal |
| GET | `/api/v1/budget/goals` | List budget goals |
| POST | `/api/v1/autosave/rules` | Create autosave rule |
| POST | `/api/v1/autosave/roundup/execute` | Execute roundup autosave |

### Sprint 4: Subscriptions

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/subscriptions` | List subscriptions |
| POST | `/api/v1/subscriptions/track` | Track a subscription |
| POST | `/api/v1/subscriptions/autopay` | Set subscription autopay |

### Sprint 4: Fraud

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/fraud/alerts` | List fraud alerts |
| POST | `/api/v1/fraud/review` | Review fraud alert |

### Sprint 4: Notifications

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/notifications` | List notifications |
| GET | `/api/v1/notifications/preferences` | Get notification preferences |
| PUT | `/api/v1/notifications/:id/read` | Mark notification as read |
| PUT | `/api/v1/notifications/preferences` | Update notification preferences |

### Sprint 4: Social Payments

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/social/feed` | Get social payment feed |
| POST | `/api/v1/social/post` | Create social feed post |
| POST | `/api/v1/social/react` | React to a feed post |
| GET | `/api/v1/social/friends` | List friends |
| POST | `/api/v1/social/friends/add` | Add friend |

### Sprint 4: Merchant Invoices and Payment Links

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/v1/merchant/invoices` | Create invoice |
| GET | `/api/v1/merchant/invoices` | List invoices |
| POST | `/api/v1/merchant/payment-links` | Create payment link |
| GET | `/api/v1/merchant/payment-links` | List payment links |

### Sprint 4: Admin

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/api/v1/admin/metrics` | Platform metrics |
| GET | `/api/v1/admin/users` | User management data |
| GET | `/api/v1/admin/fraud/review` | Fraud review queue |

### Trading Microservice API

Base URL:

```text
http://localhost:8081/api/v1
```

Protected trading routes require the same `Authorization: Bearer <jwt_token>` format.

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/health` | Trading service health check |
| GET | `/stocks/search` | Search stock symbols |
| GET | `/stocks/market-summary` | Get market summary |
| GET | `/stocks/:symbol/quote` | Get stock quote |
| GET | `/stocks/:symbol/details` | Get stock details |
| GET | `/stocks/:symbol/chart` | Get chart data |
| POST | `/trading/verify` | Verify trading account |
| GET | `/trading/account` | Get trading account |
| POST | `/trading/trade` | Execute trade |
| GET | `/trading/portfolio` | Get portfolio |
| GET | `/trading/orders` | Get order history |

## Frontend Routes

| Route | Feature |
| --- | --- |
| `/login` | User login |
| `/register` | User registration |
| `/dashboard` | Dashboard |
| `/wallet` | Wallet |
| `/transactions` | Transaction history |
| `/transfer` | Money transfer |
| `/bills` | Bill payments |
| `/rewards` | Rewards and offers |
| `/profile` | User profile |
| `/trading` | Trading |
| `/scan` | QR scan |
| `/loans` | Loans |
| `/cards` | Virtual cards |
| `/insights` | Sprint 4 financial insights |
| `/budget-planner` | Sprint 4 budget planner |
| `/subscriptions` | Sprint 4 subscription tracking |
| `/social` | Sprint 4 social payments |
| `/notifications` | Sprint 4 notifications |
| `/merchant-portal` | Sprint 4 merchant portal |
| `/admin` | Sprint 4 admin console |

## Test Case Inventory

The repository currently contains 82 Angular unit test cases, 22 Cypress E2E test cases, 64 main backend Go test cases, and 2 trading microservice Go test cases.

### Frontend Unit Test Cases

| File | Test cases |
| --- | --- |
| `core/guards/auth.guard.spec.ts` | allow auth route with token; redirect auth route without token; allow guest route without token; redirect guest route with token |
| `core/interceptors/auth.interceptor.spec.ts` | interceptor defined; interceptor function type; interceptor exists as function |
| `core/services/auth.service.spec.ts` | service creation; login endpoint; register endpoint; verify OTP endpoint; resend OTP endpoint; Google auth endpoint; no stored token; stored token; handle auth stores token and updates state; logout clears state; clear OTP state; unauthenticated without user; authenticated with user |
| `core/services/wallet.service.spec.ts` | service creation; add money endpoint; withdraw endpoint; paginated transaction endpoint; default transaction params; refresh wallet calls auth user refresh |
| `core/services/transfer.service.spec.ts` | service creation; send money endpoint; contacts endpoint; search endpoint with query; search query encoding |
| `core/services/bill.service.spec.ts` | service creation; categories endpoint; billers endpoint without category; billers endpoint with category; pay bill endpoint; saved billers endpoint; remove saved biller endpoint |
| `core/services/reward.service.spec.ts` | service creation; rewards summary endpoint; reward history with pagination; reward history default pagination; offers endpoint |
| `features/auth/login/login.component.spec.ts` | component creation; password visibility toggle; empty email invalid; invalid email invalid; valid email valid; submit rejects invalid email; submit calls login for valid email; login failure error; OTP rejects non-6-digit code; back to login resets OTP state |
| `features/auth/register/register.component.spec.ts` | component creation; password visibility toggle; empty email invalid; valid email valid; short phone invalid; 10-digit phone valid; formatted phone valid; empty form invalid; valid form accepted; short password invalid; short username invalid; invalid form does not call service; valid form calls register; back to form resets OTP state |
| `features/dashboard/dashboard.component.spec.ts` | component creation; five quick actions; expected quick action routes; balance formatted as USD; missing wallet returns zero balance; loads recent transactions on init; loading flag cleared; recent transaction signal populated |
| `features/trading/trading.component.spec.ts` | component creation; verification state mocked on init |
| `features/loans/loans.component.spec.ts` | component creation; active tab initialization |
| `features/cards/cards.component.spec.ts` | component creation; card number formatting |
| `features/scan/scan.component.spec.ts` | component creation |

### Cypress E2E Test Cases

| File | Test cases |
| --- | --- |
| `cypress/e2e/login.cy.ts` | displays login form; fills email and password; clicks login button; shows register link |
| `cypress/e2e/register.cy.ts` | displays registration form fields; allows typing in registration fields; shows login link |
| `cypress/e2e/dashboard.cy.ts` | displays welcome message; displays wallet balance; displays quick action links; lists recent transactions |
| `cypress/e2e/wallet.cy.ts` | displays wallet balance; displays add money and withdraw toggles; shows transaction history; displays reward statistics; displays active offers |
| `cypress/e2e/transfer.cy.ts` | displays transfer form; shows user balance; displays recent contacts |
| `cypress/e2e/bills.cy.ts` | displays bill categories; displays saved billers; navigates to pay a new bill |

### Backend Go Test Cases

| File | Test cases |
| --- | --- |
| `config/config_test.go` | default config loading; environment config loading; fallback env helper; env helper with value |
| `utils/response_test.go` | success response; success response with nil data; error response; unauthorized error response; internal server error response |
| `middleware/auth_test.go` | missing auth header; invalid auth header format; invalid token; valid token; bearer-only header |
| `services/token_service_test.go` | generate token; validate token; invalid signature; malformed token; empty token; different user IDs create different tokens; constructor; token round trip |
| `services/auth_service_test.go` | valid email; invalid email; valid phone; invalid phone; mask email; invalid mask email format; mask email with empty name; constructor |
| `services/reward_service_test.go` | offers retrieval; offer content; all offers active; offers have IDs; offers have descriptions; offers have discounts; offers have icons; constructor |
| `services/qr_service_test.go` | QR generation |
| `services/loan_service_test.go` | loan application; invalid loan application values |
| `services/card_service_test.go` | card creation |
| `services/statement_service_test.go` | statement generation |
| `handlers/auth_handler_test.go` | constructor; invalid register input; missing register fields; invalid login input; missing login password; invalid OTP input; invalid resend OTP input; unauthenticated current user; invalid Google auth input |
| `handlers/wallet_handler_test.go` | constructor; invalid add money input; missing add money fields; invalid withdraw input; missing withdraw bank account |
| `handlers/transfer_handler_test.go` | constructor; invalid send money input; missing recipient |
| `handlers/bill_handler_test.go` | constructor; invalid pay bill input; missing pay bill fields |
| `handlers/reward_handler_test.go` | constructor |

### Trading Microservice Go Test Cases

| File | Test cases |
| --- | --- |
| `flowpay-trading/services/trading_service_test.go` | verify trading account; execute trade |

## Useful Development Notes

- Run the backend before using frontend features that call `http://localhost:8080/api/v1`.
- Run the trading service before using the frontend trading screen.
- Cypress uses `http://localhost:4200` as its base URL.
- The main backend uses PostgreSQL and requires `DB_DSN`.
- The trading service uses local SQLite and creates or updates `trading.db` automatically.
- Protected APIs require a JWT token from the auth flow.
