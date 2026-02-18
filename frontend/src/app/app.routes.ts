import { Routes } from '@angular/router';
import { authGuard, guestGuard } from './core/guards/auth.guard';

export const routes: Routes = [
    {
        path: '',
        redirectTo: 'dashboard',
        pathMatch: 'full'
    },
    {
        path: 'login',
        canActivate: [guestGuard],
        loadComponent: () =>
            import('./features/auth/login/login.component').then(m => m.LoginComponent)
    },
    {
        path: 'register',
        canActivate: [guestGuard],
        loadComponent: () =>
            import('./features/auth/register/register.component').then(m => m.RegisterComponent)
    },
    {
        path: '',
        canActivate: [authGuard],
        loadComponent: () =>
            import('./features/layout/layout.component').then(m => m.LayoutComponent),
        children: [
            {
                path: 'dashboard',
                loadComponent: () =>
                    import('./features/dashboard/dashboard.component').then(m => m.DashboardComponent)
            },
            {
                path: 'wallet',
                loadComponent: () =>
                    import('./features/wallet/wallet.component').then(m => m.WalletComponent)
            },
            {
                path: 'transactions',
                loadComponent: () =>
                    import('./features/wallet/transactions/transactions.component').then(m => m.TransactionsComponent)
            }
        ]
    },
    {
        path: '**',
        redirectTo: 'dashboard'
    }
];
