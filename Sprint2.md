# Sprint 2 — GatorPay

## Work Completed in Sprint 2

### Backend Features Implemented
- **P2P Transfers** — Send money between users by username/email/phone, with cashback rewards
- **Bill Payments** — Pay bills to saved billers, category-based biller discovery, saved billers management
- **Rewards System** — Cashback on transfers (1%) and bill payments (2%), points tracking, promotional offers
- **Transaction History** — Paginated transaction listing with sender/receiver details
- **OTP Verification** — Secure 6-digit OTP for registration and login with email delivery
- **Wallet Operations** — Deposit, withdraw, and real-time balance tracking
- **JWT Authentication** — Token-based auth with middleware protection on all secure endpoints
- **Google OAuth** — Third-party social authentication support

### Frontend Features Implemented
- **Login & Register** — Full auth flow with OTP verification step
- **Dashboard** — Balance overview, quick actions, recent transactions
- **Wallet** — Add money, withdraw, transaction history with pagination
- **Transfer** — P2P send money with user search and recent contacts
- **Bills** — Category browsing, biller selection, bill payment, saved billers
- **Rewards** — Points summary, cashback history, promotional offers
- **Profile** — User profile display
- **Auth Guard & Interceptor** — Route protection and automatic JWT injection

---

## Frontend Unit Tests (Jasmine/Karma)

### Service Tests

| Test File | Tests | Description |
|-----------|-------|-------------|
| `auth.service.spec.ts` | 12 | loginRequest, registerRequest, verifyOTP, resendOTP, googleAuthRequest, getToken (×2), handleAuth, logout, clearOTPState, isAuthenticated (×2) |
| `bill.service.spec.ts` | 7 | getCategories, getBillers (×2), payBill, getSavedBillers, removeSavedBiller, creation |
| `reward.service.spec.ts` | 5 | getSummary, getHistory (×2), getOffers, creation |
| `transfer.service.spec.ts` | 5 | sendMoney, getContacts, searchUsers (×2), creation |
| `wallet.service.spec.ts` | 6 | addMoney, withdraw, getTransactions (×2), refreshWallet, creation |

### Guard & Interceptor Tests

| Test File | Tests | Description |
|-----------|-------|-------------|
| `auth.guard.spec.ts` | 4 | authGuard allow/deny, guestGuard allow/deny |
| `auth.interceptor.spec.ts` | 3 | interceptor definition, function signature, type check |

### Component Tests

| Test File | Tests | Description |
|-----------|-------|-------------|
| `login.component.spec.ts` | 10 | creation, togglePassword, isEmailValid (×3), onSubmit validation (×3), onVerifyOTP, backToLogin |
| `register.component.spec.ts` | 12 | creation, togglePassword, isEmailValid (×2), isPhoneValid (×3), isFormValid (×4), onSubmit (×2), backToForm |
| `dashboard.component.spec.ts` | 7 | creation, quickActions count, quickActions routes, getFormattedBalance (×2), loadRecentTransactions, loading state |

**Run frontend unit tests:**
```bash
cd frontend
npx ng test --watch=false --browsers=ChromeHeadless
```

---

## Frontend Cypress E2E Test

| Test File | Tests | Description |
|-----------|-------|-------------|
| `cypress/e2e/login.cy.ts` | 4 | display login form, fill email/password, click submit button, register link exists |

**Cypress config:** `frontend/cypress.config.js`

**Run Cypress tests:**
```bash
cd frontend
npx cypress run          # headless
npx cypress open         # interactive
```

> **Note:** Cypress E2E tests require both the Angular dev server (`ng serve`) and Go backend to be running.

---

## Backend Unit Tests (Go)

| Test File | Tests | Description |
|-----------|-------|-------------|
| `config/config_test.go` | 4 | Load default values, Load from env, getEnv with fallback, getEnv with value |
| `utils/response_test.go` | 5 | SuccessResponse, SuccessResponse nil data, ErrorResponse, ErrorResponse 401, ErrorResponse 500 |
| `services/token_service_test.go` | 8 | GenerateToken, ValidateToken, invalid signature, malformed token, empty string, different user IDs, NewTokenService, round-trip |
| `services/auth_service_test.go` | 8 | validateEmail valid (×1), validateEmail invalid (×1), validatePhone valid (×1), validatePhone invalid (×1), maskEmail (×1), maskEmail invalid format, maskEmail empty name, NewAuthService |
| `services/reward_service_test.go` | 8 | GetOffers count, content, all active, have IDs, have descriptions, have discounts, have icons, NewRewardService |
| `handlers/auth_handler_test.go` | 9 | NewAuthHandler, Register invalid JSON, Register missing fields, Login invalid JSON, Login missing password, VerifyOTP invalid, ResendOTP invalid, GetMe unauthenticated, GoogleAuth invalid |
| `handlers/bill_handler_test.go` | 3 | NewBillHandler, PayBill invalid JSON, PayBill missing fields |
| `handlers/transfer_handler_test.go` | 3 | NewTransferHandler, SendMoney invalid JSON, SendMoney missing recipient |
| `handlers/wallet_handler_test.go` | 5 | NewWalletHandler, AddMoney invalid JSON, AddMoney missing fields, Withdraw invalid JSON, Withdraw missing bank_account |
| `handlers/reward_handler_test.go` | 1 | NewRewardHandler |
| `middleware/auth_test.go` | 5 | missing header, invalid format, invalid token, valid token with userID, bare Bearer |

**Total: 59 backend tests — all passing ✅**

**Run backend unit tests:**
```bash
cd backend
go test ./... -v
```

---

## Backend API Documentation

**Base URL:** `http://localhost:8080/api/v1`

All responses use a standard envelope:
```json
{
  "success": true,
  "message": "Description",
  "data": { ... }
}
```

### Authentication

#### POST `/auth/register`
Register a new user. Sends OTP to provided email.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | ✅ | User email |
| password | string | ✅ | Min 8 characters |
| username | string | ✅ | Min 3 characters |
| phone | string | ✅ | 10-digit phone |
| first_name | string | ✅ | First name |
| last_name | string | ✅ | Last name |

**Response (201):** `{ user_id, email (masked), purpose: "register" }`

---

#### POST `/auth/login`
Login with credentials. Sends OTP to registered email.

| Field | Type | Required |
|-------|------|----------|
| email | string | ✅ |
| password | string | ✅ |

**Response (200):** `{ user_id, email (masked), purpose: "login" }`

---

#### POST `/auth/verify-otp`
Verify OTP code to complete authentication. Returns JWT token.

| Field | Type | Required |
|-------|------|----------|
| user_id | string | ✅ |
| code | string | ✅ | 6 digits |
| purpose | string | ✅ | "login" or "register" |

**Response (200):** `{ token, user, wallet }`

---

#### POST `/auth/resend-otp`
Resend a new OTP code.

| Field | Type | Required |
|-------|------|----------|
| user_id | string | ✅ |
| purpose | string | ✅ |

**Response (200):** `{ user_id, email (masked), purpose }`

---

#### POST `/auth/google`
Google OAuth authentication (no OTP required).

| Field | Type | Required |
|-------|------|----------|
| google_id | string | ✅ |
| email | string | ✅ |
| name | string | ✅ |
| avatar | string | ❌ |

**Response (200):** `{ token, user, wallet }`

---

#### GET `/auth/me` 🔒
Get current authenticated user profile with wallet.

**Headers:** `Authorization: Bearer <token>`

**Response (200):** `{ user, wallet }`

---

### Wallet 🔒

> All wallet endpoints require `Authorization: Bearer <token>` header.

#### POST `/wallet/add`
Add money to wallet.

| Field | Type | Required |
|-------|------|----------|
| amount | float | ✅ |
| source | string | ✅ |
| description | string | ❌ |

**Response (200):** `{ id, user_id, balance, currency, is_active }`

---

#### POST `/wallet/withdraw`
Withdraw money from wallet.

| Field | Type | Required |
|-------|------|----------|
| amount | float | ✅ |
| bank_account | string | ✅ |

**Response (200):** `{ id, user_id, balance, currency, is_active }`

---

#### GET `/wallet/transactions?page=1&limit=10`
Get paginated transaction history.

| Param | Type | Default |
|-------|------|---------|
| page | int | 1 |
| limit | int | 20 (max 100) |

**Response (200):** `{ transactions[], total, page, limit, total_pages }`

---

### Transfers 🔒

#### POST `/transfer/send`
Send money to another user (P2P transfer). Awards 1% cashback.

| Field | Type | Required |
|-------|------|----------|
| recipient | string | ✅ | username, email, or phone |
| amount | float | ✅ |
| note | string | ❌ |

**Response (200):** `{ transaction_id, recipient, amount, note, new_balance }`

---

#### GET `/transfer/contacts`
Get 10 most recent unique transfer recipients.

**Response (200):** `[ { id, username, email, first_name, last_name, ... } ]`

---

#### GET `/transfer/search?query=<search>`
Search users by username, email, or full name.

**Response (200):** `[ { id, username, email, first_name, last_name, ... } ]`

---

### Bills 🔒

#### GET `/bills/categories`
Get distinct bill categories.

**Response (200):** `[ "electricity", "internet", "phone", ... ]`

---

#### GET `/bills/billers?category=<category>`
Get billers, optionally filtered by category.

**Response (200):** `[ { id, name, category, icon, is_active } ]`

---

#### POST `/bills/pay`
Pay a bill. Awards 2% cashback.

| Field | Type | Required |
|-------|------|----------|
| biller_id | string | ✅ |
| account_number | string | ✅ |
| amount | float | ✅ |
| save_biller | bool | ❌ |

**Response (200):** `{ payment_id, biller, amount, new_balance }`

---

#### GET `/bills/saved`
Get saved billers for the current user.

**Response (200):** `[ { id, user_id, biller_id, account_number, nickname, biller } ]`

---

#### DELETE `/bills/saved/:id`
Remove a saved biller.

**Response (200):** `null`

---

### Rewards 🔒

#### GET `/rewards`
Get reward summary (total points, cashback, lifetime earnings).

**Response (200):** `{ total_points, total_cashback, lifetime_earnings, total_transactions }`

---

#### GET `/rewards/history?page=1&limit=10`
Get paginated reward history.

**Response (200):** `{ rewards[], total, page, limit }`

---

#### GET `/rewards/offers`
Get promotional offers.

**Response (200):** `[ { id, title, description, discount, icon, is_active } ]`
