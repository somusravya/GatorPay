import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from '../services/auth.service';

/** Protects routes that require authentication */
export const authGuard: CanActivateFn = () => {
    const authService = inject(AuthService);
    const router = inject(Router);

    if (authService.getToken()) {
        return true;
    }

    router.navigate(['/login']);
    return false;
};

/** Protects routes that should only be accessible to guests (login, register) */
export const guestGuard: CanActivateFn = () => {
    const authService = inject(AuthService);
    const router = inject(Router);

    if (!authService.getToken()) {
        return true;
    }

    router.navigate(['/dashboard']);
    return false;
};
