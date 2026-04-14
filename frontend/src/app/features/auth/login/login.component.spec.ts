import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { RouterModule } from '@angular/router';
import { LoginComponent } from './login.component';
import { AuthService } from '../../../core/services/auth.service';
import { of, throwError } from 'rxjs';

describe('LoginComponent', () => {
    let component: LoginComponent;
    let fixture: ComponentFixture<LoginComponent>;
    let authServiceSpy: jasmine.SpyObj<AuthService>;

    beforeEach(async () => {
        authServiceSpy = jasmine.createSpyObj('AuthService', [
            'loginRequest', 'verifyOTP', 'resendOTP', 'handleAuth', 'getToken'
        ]);
        authServiceSpy.getToken.and.returnValue(null);

        await TestBed.configureTestingModule({
            imports: [LoginComponent, RouterModule.forRoot([])],
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                { provide: AuthService, useValue: authServiceSpy }
            ]
        }).compileComponents();

        fixture = TestBed.createComponent(LoginComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    // --- togglePassword ---
    it('togglePassword should toggle showPassword signal', () => {
        expect(component.showPassword()).toBeFalse();
        component.togglePassword();
        expect(component.showPassword()).toBeTrue();
        component.togglePassword();
        expect(component.showPassword()).toBeFalse();
    });

    // --- isEmailValid ---
    it('isEmailValid should return false for empty string', () => {
        component.email = '';
        expect(component.isEmailValid).toBeFalse();
    });

    it('isEmailValid should return false for invalid email', () => {
        component.email = 'invalid';
        expect(component.isEmailValid).toBeFalse();
    });

    it('isEmailValid should return true for valid email', () => {
        component.email = 'test@gator.edu';
        expect(component.isEmailValid).toBeTrue();
    });

    // --- onSubmit validation ---
    it('onSubmit should set error when email is invalid', () => {
        component.email = 'bad-email';
        component.password = 'password123';
        component.onSubmit();
        expect(component.error()).toBe('Please enter a valid email address');
    });

    it('onSubmit should call loginRequest when email is valid', () => {
        authServiceSpy.loginRequest.and.returnValue(of({
            success: true,
            message: 'OTP sent',
            data: { user_id: '1', email: 't***@gator.edu', purpose: 'login' }
        }));

        component.email = 'test@gator.edu';
        component.password = 'password123';
        component.onSubmit();

        expect(authServiceSpy.loginRequest).toHaveBeenCalledWith({
            email: 'test@gator.edu',
            password: 'password123'
        });
        expect(component.otpStep()).toBeTrue();
    });

    it('onSubmit should set error on login failure', () => {
        authServiceSpy.loginRequest.and.returnValue(throwError(() => ({
            error: { message: 'Invalid credentials' }
        })));

        component.email = 'test@gator.edu';
        component.password = 'wrong';
        component.onSubmit();

        expect(component.error()).toBe('Invalid credentials');
        expect(component.loading()).toBeFalse();
    });

    // --- onVerifyOTP ---
    it('onVerifyOTP should set error when code is not 6 digits', () => {
        component.otpCode = '123';
        component.onVerifyOTP();
        expect(component.error()).toBe('Please enter the 6-digit code');
    });

    // --- backToLogin ---
    it('backToLogin should reset OTP state', () => {
        component.otpStep.set(true);
        component.otpCode = '123456';
        component.error.set('some error');
        component.success.set('some success');

        component.backToLogin();

        expect(component.otpStep()).toBeFalse();
        expect(component.otpCode).toBe('');
        expect(component.error()).toBe('');
        expect(component.success()).toBe('');
    });
});
