import { Component, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../core/services/auth.service';
import { WalletService } from '../../core/services/wallet.service';

@Component({
    selector: 'app-wallet',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './wallet.component.html',
    styleUrl: './wallet.component.scss'
})
export class WalletComponent {
    private authService = inject(AuthService);
    private walletService = inject(WalletService);

    wallet = this.authService.currentWallet;
    activeTab = signal<'add' | 'withdraw'>('add');

    // Add Money form
    addAmount = 0;
    addSource = 'Bank Account';
    addDescription = '';
    addLoading = signal(false);
    addError = signal('');
    addSuccess = signal('');

    // Withdraw form
    withdrawAmount = 0;
    withdrawBankAccount = '';
    withdrawLoading = signal(false);
    withdrawError = signal('');
    withdrawSuccess = signal('');

    getFormattedBalance(): string {
        const balance = this.wallet()?.balance || '0';
        return parseFloat(balance).toLocaleString('en-US', {
            style: 'currency',
            currency: 'USD'
        });
    }

    switchTab(tab: 'add' | 'withdraw'): void {
        this.activeTab.set(tab);
        this.clearMessages();
    }

    clearMessages(): void {
        this.addError.set('');
        this.addSuccess.set('');
        this.withdrawError.set('');
        this.withdrawSuccess.set('');
    }

    onAddMoney(): void {
        if (this.addAmount <= 0) {
            this.addError.set('Amount must be greater than 0');
            return;
        }

        this.clearMessages();
        this.addLoading.set(true);

        this.walletService.addMoney({
            amount: this.addAmount,
            source: this.addSource,
            description: this.addDescription || `Deposit from ${this.addSource}`
        }).subscribe({
            next: (res: any) => {
                this.addLoading.set(false);
                if (res.success) {
                    this.addSuccess.set(`$${this.addAmount.toFixed(2)} added successfully!`);
                    this.authService.currentWallet.set(res.data);
                    this.addAmount = 0;
                    this.addDescription = '';
                }
            },
            error: (err: any) => {
                this.addLoading.set(false);
                this.addError.set(err.error?.message || 'Failed to add money');
            }
        });
    }

    onWithdraw(): void {
        if (this.withdrawAmount <= 0) {
            this.withdrawError.set('Amount must be greater than 0');
            return;
        }
        if (!this.withdrawBankAccount) {
            this.withdrawError.set('Bank account is required');
            return;
        }

        this.clearMessages();
        this.withdrawLoading.set(true);

        this.walletService.withdraw({
            amount: this.withdrawAmount,
            bank_account: this.withdrawBankAccount
        }).subscribe({
            next: (res: any) => {
                this.withdrawLoading.set(false);
                if (res.success) {
                    this.withdrawSuccess.set(`$${this.withdrawAmount.toFixed(2)} withdrawn successfully!`);
                    this.authService.currentWallet.set(res.data);
                    this.withdrawAmount = 0;
                    this.withdrawBankAccount = '';
                }
            },
            error: (err: any) => {
                this.withdrawLoading.set(false);
                this.withdrawError.set(err.error?.message || 'Failed to withdraw');
            }
        });
    }
}
