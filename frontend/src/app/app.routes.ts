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
            },
            {
                path: 'transfer',
                loadComponent: () =>
                    import('./features/transfer/transfer.component').then(m => m.TransferComponent)
            },
            {
                path: 'bills',
                loadComponent: () =>
                    import('./features/bills/bills.component').then(m => m.BillsComponent)
            },
            {
                path: 'rewards',
                loadComponent: () =>
                    import('./features/rewards/rewards.component').then(m => m.RewardsComponent)
            },
            {
                path: 'help-center',
                loadComponent: () =>
                    import('./features/help-center/help-center.component').then(m => m.HelpCenterComponent)
            },
            {
                path: 'smart-tools',
                loadComponent: () =>
                    import('./features/smart-tools/smart-tools.component').then(m => m.SmartToolsComponent)
            },
            {
                path: 'profile',
                loadComponent: () =>
                    import('./features/profile/profile.component').then(m => m.ProfileComponent)
            },
            {
                path: 'analytics',
                loadComponent: () =>
                    import('./features/spending-analytics/spending-analytics.component').then(m => m.SpendingAnalyticsComponent)
            },
            {
                path: 'savings-goals',
                loadComponent: () =>
                    import('./features/savings-goals/savings-goals.component').then(m => m.SavingsGoalsComponent)
            }
        ]
    },
    {
        path: '**',
        redirectTo: 'dashboard'
    }
];
