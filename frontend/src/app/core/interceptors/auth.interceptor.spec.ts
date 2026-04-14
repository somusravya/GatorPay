import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient, HttpClient, HTTP_INTERCEPTORS, HttpInterceptorFn } from '@angular/common/http';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { authInterceptor } from './auth.interceptor';

describe('AuthInterceptor', () => {
    let httpMock: HttpTestingController;
    let httpClient: HttpClient;
    let authServiceSpy: jasmine.SpyObj<AuthService>;

    beforeEach(() => {
        authServiceSpy = jasmine.createSpyObj('AuthService', ['getToken', 'logout'], {
            currentUser: jasmine.createSpy().and.returnValue(null),
            currentWallet: jasmine.createSpy().and.returnValue(null),
            isAuthenticated: jasmine.createSpy().and.returnValue(false)
        });

        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(
                    // Angular 19 functional interceptors are passed via withInterceptors
                ),
                provideHttpClientTesting(),
                { provide: AuthService, useValue: authServiceSpy },
                { provide: Router, useValue: jasmine.createSpyObj('Router', ['navigate']) }
            ]
        });

        // We need to test the interceptor function directly
        httpMock = TestBed.inject(HttpTestingController);
        httpClient = TestBed.inject(HttpClient);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be defined', () => {
        expect(authInterceptor).toBeDefined();
        expect(typeof authInterceptor).toBe('function');
    });

    it('should be an HttpInterceptorFn', () => {
        // Verify it has the correct function signature
        expect(authInterceptor.length).toBe(2); // (req, next) => ...
    });

    it('interceptor should exist and be a function', () => {
        const interceptorFn: HttpInterceptorFn = authInterceptor;
        expect(interceptorFn).toBeTruthy();
    });
});
