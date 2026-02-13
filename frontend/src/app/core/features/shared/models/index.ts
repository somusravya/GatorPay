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
  type: string;
  amount: string;
  description: string;
  status: string;
  created_at: string;
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
