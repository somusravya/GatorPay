import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';

@Component({
    selector: 'app-login',
    standalone: true,
    imports: [CommonModule, FormsModule, RouterLink],
    templateUrl: './login.component.html',
    styleUrl: './login.component.scss'
})
export class LoginComponent {
    email = '';
    password = '';
    showPassword = signal(false);
    loading = signal(false);
    error = signal('');
    success = signal('');

    // OTP flow
    otpStep = signal(false);
    otpCode = '';
    pendingUserID = '';
    maskedEmail = '';
    resendCooldown = signal(0);
    private cooldownTimer: any;

    constructor(private authService: AuthService, private router: Router) { }

    togglePassword(): void {
        this.showPassword.update(v => !v);
    }

    get isEmailValid(): boolean {
        return /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(this.email);
    }

    onSubmit(): void {
        if (!this.isEmailValid) {
            this.error.set('Please enter a valid email address');
            return;
        }

        this.error.set('');
        this.loading.set(true);

        this.authService.loginRequest({ email: this.email, password: this.password })
            .subscribe({
                next: (res: any) => {
                    this.loading.set(false);
                    if (res.success) {
                        this.pendingUserID = res.data.user_id;
                        this.maskedEmail = res.data.email;
                        this.otpStep.set(true);
                        this.success.set('Verification code sent to ' + res.data.email);
                    }
                },
                error: (err: any) => {
                    this.loading.set(false);
                    this.error.set(err.error?.message || 'Login failed. Please try again.');
                }
            });
    }

    onVerifyOTP(): void {
        if (this.otpCode.length !== 6) {
            this.error.set('Please enter the 6-digit code');
            return;
        }

        this.error.set('');
        this.success.set('');
        this.loading.set(true);

        this.authService.verifyOTP({
            user_id: this.pendingUserID,
            code: this.otpCode,
            purpose: 'login'
        }).subscribe({
            next: (res: any) => {
                this.loading.set(false);
                if (res.success) {
                    this.authService.handleAuth(res.data);
                    this.router.navigate(['/dashboard']);
                }
            },
            error: (err: any) => {
                this.loading.set(false);
                this.error.set(err.error?.message || 'Invalid code. Please try again.');
                this.otpCode = '';
            }
        });
    }

    resendOTP(): void {
        if (this.resendCooldown() > 0) return;

        this.error.set('');
        this.loading.set(true);

        this.authService.resendOTP({
            user_id: this.pendingUserID,
            purpose: 'login'
        }).subscribe({
            next: (res: any) => {
                this.loading.set(false);
                if (res.success) {
                    this.success.set('New code sent to ' + res.data.email);
                    this.startCooldown();
                }
            },
            error: (err: any) => {
                this.loading.set(false);
                this.error.set(err.error?.message || 'Failed to resend code');
            }
        });
    }

    backToLogin(): void {
        this.otpStep.set(false);
        this.otpCode = '';
        this.error.set('');
        this.success.set('');
    }

    private startCooldown(): void {
        this.resendCooldown.set(30);
        this.cooldownTimer = setInterval(() => {
            const val = this.resendCooldown();
            if (val <= 1) {
                this.resendCooldown.set(0);
                clearInterval(this.cooldownTimer);
            } else {
                this.resendCooldown.set(val - 1);
            }
        }, 1000);
    }
}
