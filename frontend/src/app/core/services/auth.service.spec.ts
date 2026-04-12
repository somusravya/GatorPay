import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { AuthService } from './auth.service';
import { environment } from '../../../environments/environment';

describe('AuthService', () => {
    let service: AuthService;
    let httpMock: HttpTestingController;
    let routerSpy: jasmine.SpyObj<Router>;
    const apiUrl = `${environment.apiUrl}/auth`;

    beforeEach(() => {
        routerSpy = jasmine.createSpyObj('Router', ['navigate']);
        localStorage.clear();

        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                AuthService,
                { provide: Router, useValue: routerSpy }
            ]
        });

        service = TestBed.inject(AuthService);
        httpMock = TestBed.inject(HttpTestingController);

        // Flush the restoreSession call if no token
        httpMock.expectNone(`${apiUrl}/me`);
    });

    afterEach(() => {
        httpMock.verify();
        localStorage.clear();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    // --- loginRequest ---
    it('loginRequest should POST to /auth/login', () => {
        const loginData = { email: 'test@gator.edu', password: 'password123' };
        service.loginRequest(loginData).subscribe(res => {
            expect(res.success).toBeTrue();
        });

        const req = httpMock.expectOne(`${apiUrl}/login`);
        expect(req.request.method).toBe('POST');
        expect(req.request.body).toEqual(loginData);
        req.flush({ success: true, message: 'OTP sent', data: { user_id: '1', email: 't***@gator.edu', purpose: 'login' } });
    });

    // --- registerRequest ---
    it('registerRequest should POST to /auth/register', () => {
        const registerData = {
            email: 'new@gator.edu', password: 'password123',
            username: 'newuser', phone: '1234567890',
            first_name: 'First', last_name: 'Last'
        };
        service.registerRequest(registerData).subscribe(res => {
            expect(res.success).toBeTrue();
        });

        const req = httpMock.expectOne(`${apiUrl}/register`);
        expect(req.request.method).toBe('POST');
        req.flush({ success: true, message: 'OTP sent', data: { user_id: '2', email: 'n***@gator.edu', purpose: 'register' } });
    });

    // --- verifyOTP ---
    it('verifyOTP should POST to /auth/verify-otp', () => {
        const otpData = { user_id: '1', code: '123456', purpose: 'login' };
        service.verifyOTP(otpData).subscribe(res => {
            expect(res.success).toBeTrue();
        });

        const req = httpMock.expectOne(`${apiUrl}/verify-otp`);
        expect(req.request.method).toBe('POST');
        req.flush({ success: true, message: 'Verified', data: { token: 'jwt-token', user: {}, wallet: {} } });
    });

    // --- resendOTP ---
    it('resendOTP should POST to /auth/resend-otp', () => {
        const resendData = { user_id: '1', purpose: 'login' };
        service.resendOTP(resendData).subscribe(res => {
            expect(res.success).toBeTrue();
        });

        const req = httpMock.expectOne(`${apiUrl}/resend-otp`);
        expect(req.request.method).toBe('POST');
        req.flush({ success: true, message: 'Resent', data: { user_id: '1', email: 't***@gator.edu', purpose: 'login' } });
    });

    // --- googleAuthRequest ---
    it('googleAuthRequest should POST to /auth/google', () => {
        const googleData = { google_id: 'g123', email: 'g@gmail.com', name: 'Google User', avatar: '' };
        service.googleAuthRequest(googleData).subscribe(res => {
            expect(res.success).toBeTrue();
        });

        const req = httpMock.expectOne(`${apiUrl}/google`);
        expect(req.request.method).toBe('POST');
        req.flush({ success: true, message: 'Auth success', data: { token: 'jwt', user: {}, wallet: {} } });
    });

    // --- getToken ---
    it('getToken should return null when no token stored', () => {
        expect(service.getToken()).toBeNull();
    });

    it('getToken should return token when stored', () => {
        localStorage.setItem('gatorpay_token', 'test-token');
        expect(service.getToken()).toBe('test-token');
    });

    // --- handleAuth ---
    it('handleAuth should store token and update signals', () => {
        const authData = {
            token: 'new-jwt',
            user: { id: '1', email: 'test@gator.edu', username: 'tester' } as any,
            wallet: { id: 'w1', balance: '100.00' } as any
        };

        service.handleAuth(authData);

        expect(localStorage.getItem('gatorpay_token')).toBe('new-jwt');
        expect(service.currentUser()).toEqual(authData.user);
        expect(service.currentWallet()).toEqual(authData.wallet);
    });

    // --- logout ---
    it('logout should clear state and navigate to /login', () => {
        localStorage.setItem('gatorpay_token', 'token');
        service.currentUser.set({ id: '1' } as any);
        service.currentWallet.set({ id: 'w1' } as any);

        service.logout();

        expect(localStorage.getItem('gatorpay_token')).toBeNull();
        expect(service.currentUser()).toBeNull();
        expect(service.currentWallet()).toBeNull();
        expect(routerSpy.navigate).toHaveBeenCalledWith(['/login']);
    });

    // --- clearOTPState ---
    it('clearOTPState should reset OTP signals', () => {
        service.pendingUserID.set('123');
        service.pendingEmail.set('test@gator.edu');
        service.pendingPurpose.set('login');

        service.clearOTPState();

        expect(service.pendingUserID()).toBeNull();
        expect(service.pendingEmail()).toBe('');
        expect(service.pendingPurpose()).toBe('');
    });

    // --- isAuthenticated ---
    it('isAuthenticated should be false when no user', () => {
        expect(service.isAuthenticated()).toBeFalse();
    });

    it('isAuthenticated should be true when user is set', () => {
        service.currentUser.set({ id: '1' } as any);
        expect(service.isAuthenticated()).toBeTrue();
    });
});
