import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { RouterModule } from '@angular/router';
import { RegisterComponent } from './register.component';
import { AuthService } from '../../../core/services/auth.service';
import { of, throwError } from 'rxjs';

describe('RegisterComponent', () => {
    let component: RegisterComponent;
    let fixture: ComponentFixture<RegisterComponent>;
    let authServiceSpy: jasmine.SpyObj<AuthService>;

    beforeEach(async () => {
        authServiceSpy = jasmine.createSpyObj('AuthService', [
            'registerRequest', 'verifyOTP', 'resendOTP', 'handleAuth', 'getToken'
        ]);
        authServiceSpy.getToken.and.returnValue(null);

        await TestBed.configureTestingModule({
            imports: [RegisterComponent, RouterModule.forRoot([])],
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                { provide: AuthService, useValue: authServiceSpy }
            ]
        }).compileComponents();

        fixture = TestBed.createComponent(RegisterComponent);
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
    });

    // --- isEmailValid ---
    it('isEmailValid should return false for empty string', () => {
        component.email = '';
        expect(component.isEmailValid).toBeFalse();
    });

    it('isEmailValid should return true for valid email', () => {
        component.email = 'user@ufl.edu';
        expect(component.isEmailValid).toBeTrue();
    });

    // --- isPhoneValid ---
    it('isPhoneValid should return false for short phone', () => {
        component.phone = '12345';
        expect(component.isPhoneValid).toBeFalse();
    });

    it('isPhoneValid should return true for 10-digit phone', () => {
        component.phone = '3521234567';
        expect(component.isPhoneValid).toBeTrue();
    });

    it('isPhoneValid should handle formatted phone numbers', () => {
        component.phone = '(352) 123-4567';
        expect(component.isPhoneValid).toBeTrue();
    });

    // --- isFormValid ---
    it('isFormValid should return false when fields are empty', () => {
        expect(component.isFormValid).toBeFalse();
    });

    it('isFormValid should return true when all fields are valid', () => {
        component.email = 'test@ufl.edu';
        component.password = 'password123';
        component.username = 'testuser';
        component.phone = '3521234567';
        component.firstName = 'John';
        component.lastName = 'Doe';
        expect(component.isFormValid).toBeTrue();
    });

    it('isFormValid should return false when password too short', () => {
        component.email = 'test@ufl.edu';
        component.password = 'short';
        component.username = 'testuser';
        component.phone = '3521234567';
        component.firstName = 'John';
        component.lastName = 'Doe';
        expect(component.isFormValid).toBeFalse();
    });

    it('isFormValid should return false when username too short', () => {
        component.email = 'test@ufl.edu';
        component.password = 'password123';
        component.username = 'ab';
        component.phone = '3521234567';
        component.firstName = 'John';
        component.lastName = 'Doe';
        expect(component.isFormValid).toBeFalse();
    });

    // --- onSubmit ---
    it('onSubmit should not call service when form is invalid', () => {
        component.onSubmit();
        expect(authServiceSpy.registerRequest).not.toHaveBeenCalled();
    });

    it('onSubmit should call registerRequest when form is valid', () => {
        authServiceSpy.registerRequest.and.returnValue(of({
            success: true,
            message: 'OTP sent',
            data: { user_id: '2', email: 't***@ufl.edu', purpose: 'register' }
        }));

        component.email = 'test@ufl.edu';
        component.password = 'password123';
        component.username = 'testuser';
        component.phone = '3521234567';
        component.firstName = 'John';
        component.lastName = 'Doe';
        component.onSubmit();

        expect(authServiceSpy.registerRequest).toHaveBeenCalled();
        expect(component.otpStep()).toBeTrue();
    });

    // --- backToForm ---
    it('backToForm should reset OTP state', () => {
        component.otpStep.set(true);
        component.otpCode = '123456';
        component.error.set('err');
        component.success.set('ok');

        component.backToForm();

        expect(component.otpStep()).toBeFalse();
        expect(component.otpCode).toBe('');
        expect(component.error()).toBe('');
        expect(component.success()).toBe('');
    });
});
