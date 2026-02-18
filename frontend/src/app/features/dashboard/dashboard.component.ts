import { Component, signal, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { WalletService } from '../../core/services/wallet.service';
import { Transaction } from '../../shared/models';

@Component({
    selector: 'app-dashboard',
    standalone: true,
    imports: [CommonModule, RouterLink],
    templateUrl: './dashboard.component.html',
    styleUrl: './dashboard.component.scss'
})
export class DashboardComponent implements OnInit {
    private authService = inject(AuthService);
    private walletService = inject(WalletService);

    user = this.authService.currentUser;
    wallet = this.authService.currentWallet;
    recentTransactions = signal<Transaction[]>([]);
    loading = signal(true);

    quickActions = [
        { label: 'Add Money', icon: '💳', route: '/wallet', color: '#3b82f6' },
        { label: 'Withdraw', icon: '🏦', route: '/wallet', fragment: 'withdraw', color: '#8b5cf6' },
        { label: 'Transactions', icon: '📜', route: '/transactions', color: '#06b6d4' }
    ];

    ngOnInit(): void {
        this.loadRecentTransactions();
    }

    loadRecentTransactions(): void {
        this.walletService.getTransactions(1, 5).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.recentTransactions.set(res.data.transactions || []);
                }
                this.loading.set(false);
            },
            error: () => {
                this.loading.set(false);
            }
        });
    }

    getFormattedBalance(): string {
        const balance = this.wallet()?.balance || '0';
        return parseFloat(balance).toLocaleString('en-US', {
            style: 'currency',
            currency: 'USD'
        });
    }
}
