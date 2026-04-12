export interface User {
  id: string;
  email: string;
  username: string;
  phone: string;
  first_name: string;
  last_name: string;
  avatar_url: string;
  auth_provider: string;
  email_verified: boolean;
  kyc_status: string;
  credit_score: number;
  created_at: string;
  updated_at: string;
}

export interface Wallet {
  id: string;
  user_id: string;
  balance: string;
  currency: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: string;
  wallet_id: string;
  from_user_id: string;
  to_user_id: string;
  type: string;
  amount: string;
  description: string;
  status: string;
  created_at: string;
  from_user?: User;
  to_user?: User;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface AuthResponse {
  token: string;
  user: User;
  wallet: Wallet;
}

export interface OTPSentResponse {
  user_id: string;
  email: string;
  purpose: string;
}

export interface TransactionListResponse {
  transactions: Transaction[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface RegisterRequest {
  email: string;
  password: string;
  username: string;
  phone: string;
  first_name: string;
  last_name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface VerifyOTPRequest {
  user_id: string;
  code: string;
  purpose: string;
}

export interface ResendOTPRequest {
  user_id: string;
  purpose: string;
}

export interface AddMoneyRequest {
  amount: number;
  source: string;
  description: string;
}

export interface WithdrawRequest {
  amount: number;
  bank_account: string;
}
// Sprint 2 interfaces

export interface TransferRequest {
  recipient: string;
  amount: number;
  note: string;
}

export interface TransferResponse {
  transaction_id: string;
  recipient: User;
  amount: string;
  note: string;
  new_balance: string;
}

export interface Biller {
  id: string;
  name: string;
  category: string;
  icon: string;
  is_active: boolean;
  created_at: string;
}

export interface SavedBiller {
  id: string;
  user_id: string;
  biller_id: string;
  account_number: string;
  nickname: string;
  created_at: string;
  biller: Biller;
}

export interface BillPayRequest {
  biller_id: string;
  account_number: string;
  amount: number;
  save_biller: boolean;
}

export interface BillPayResponse {
  payment_id: string;
  biller: Biller;
  amount: string;
  new_balance: string;
}

export interface RewardSummary {
  total_points: number;
  total_cashback: string;
  lifetime_earnings: string;
  total_transactions: number;
}

export interface Reward {
  id: string;
  user_id: string;
  type: string;
  amount: string;
  points: number;
  transaction_id: string;
  description: string;
  created_at: string;
}

export interface RewardHistoryResponse {
  rewards: Reward[];
  total: number;
  page: number;
  limit: number;
}

export interface Offer {
  id: string;
  title: string;
  description: string;
  discount: string;
  icon: string;
  is_active: boolean;
}
