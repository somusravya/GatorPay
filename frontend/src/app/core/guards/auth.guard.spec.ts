import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { authGuard, guestGuard } from './auth.guard';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';

describe('Auth Guards', () => {
    let authServiceSpy: jasmine.SpyObj<AuthService>;
    let routerSpy: jasmine.SpyObj<Router>;
    const mockRoute = {} as ActivatedRouteSnapshot;
    const mockState = {} as RouterStateSnapshot;

    beforeEach(() => {
        authServiceSpy = jasmine.createSpyObj('AuthService', ['getToken']);
        routerSpy = jasmine.createSpyObj('Router', ['navigate']);

        TestBed.configureTestingModule({
            providers: [
                { provide: AuthService, useValue: authServiceSpy },
                { provide: Router, useValue: routerSpy }
            ]
        });
    });

    describe('authGuard', () => {
        it('should allow access when token exists', () => {
            authServiceSpy.getToken.and.returnValue('valid-token');
            const result = TestBed.runInInjectionContext(() => authGuard(mockRoute, mockState));
            expect(result).toBeTrue();
        });

        it('should redirect to /login when no token', () => {
            authServiceSpy.getToken.and.returnValue(null);
            const result = TestBed.runInInjectionContext(() => authGuard(mockRoute, mockState));
            expect(result).toBeFalse();
            expect(routerSpy.navigate).toHaveBeenCalledWith(['/login']);
        });
    });

    describe('guestGuard', () => {
        it('should allow access when no token (guest)', () => {
            authServiceSpy.getToken.and.returnValue(null);
            const result = TestBed.runInInjectionContext(() => guestGuard(mockRoute, mockState));
            expect(result).toBeTrue();
        });

        it('should redirect to /dashboard when token exists', () => {
            authServiceSpy.getToken.and.returnValue('valid-token');
            const result = TestBed.runInInjectionContext(() => guestGuard(mockRoute, mockState));
            expect(result).toBeFalse();
            expect(routerSpy.navigate).toHaveBeenCalledWith(['/dashboard']);
        });
    });
});
