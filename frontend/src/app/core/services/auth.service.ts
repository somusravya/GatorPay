import { Injectable, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment';
import {
    User, Wallet, ApiResponse, AuthResponse, OTPSentResponse,
    LoginRequest, RegisterRequest, VerifyOTPRequest, ResendOTPRequest
} from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class AuthService {
    private apiUrl = `${environment.apiUrl}/auth`;

    // Signals for reactive state
    currentUser = signal<User | null>(null);
    currentWallet = signal<Wallet | null>(null);
    isAuthenticated = computed(() => !!this.currentUser());

    // OTP flow state
    pendingUserID = signal<string | null>(null);
    pendingEmail = signal<string>('');
    pendingPurpose = signal<string>('');

    constructor(private http: HttpClient, private router: Router) {
        this.restoreSession();
    }

    /** Restore session from localStorage on app load */
    private restoreSession(): void {
        const token = localStorage.getItem('gatorpay_token');
        if (token) {
            this.getMe();
        }
    }

    /** Login — step 1: send credentials, receive OTP prompt */
    loginRequest(data: LoginRequest) {
        return this.http.post<ApiResponse<OTPSentResponse>>(`${this.apiUrl}/login`, data);
    }

    /** Register — step 1: send form, receive OTP prompt */
    registerRequest(data: RegisterRequest) {
        return this.http.post<ApiResponse<OTPSentResponse>>(`${this.apiUrl}/register`, data);
    }

    /** Verify OTP — step 2: complete auth */
    verifyOTP(data: VerifyOTPRequest) {
        return this.http.post<ApiResponse<AuthResponse>>(`${this.apiUrl}/verify-otp`, data);
    }

    /** Resend OTP */
    resendOTP(data: ResendOTPRequest) {
        return this.http.post<ApiResponse<OTPSentResponse>>(`${this.apiUrl}/resend-otp`, data);
    }

    /** Google OAuth */
    googleAuthRequest(data: { google_id: string; email: string; name: string; avatar: string }) {
        return this.http.post<ApiResponse<AuthResponse>>(`${this.apiUrl}/google`, data);
    }

    /** Get current user profile */
    getMe(): void {
        this.http.get<ApiResponse<AuthResponse>>(`${this.apiUrl}/me`)
            .subscribe({
                next: (res) => {
                    if (res.success) {
                        this.currentUser.set(res.data.user);
                        this.currentWallet.set(res.data.wallet);
                    }
                },
                error: () => {
                    this.logout();
                }
            });
    }

    /** Handle successful authentication (after OTP verified) */
    handleAuth(data: AuthResponse): void {
        localStorage.setItem('gatorpay_token', data.token);
        this.currentUser.set(data.user);
        this.currentWallet.set(data.wallet);
        this.clearOTPState();
    }

    /** Clear OTP flow state */
    clearOTPState(): void {
        this.pendingUserID.set(null);
        this.pendingEmail.set('');
        this.pendingPurpose.set('');
    }

    /** Logout and clear state */
    logout(): void {
        localStorage.removeItem('gatorpay_token');
        this.currentUser.set(null);
        this.currentWallet.set(null);
        this.clearOTPState();
        this.router.navigate(['/login']);
    }

    /** Get stored token */
    getToken(): string | null {
        return localStorage.getItem('gatorpay_token');
    }
}
