import { Component, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { WalletService } from '../../../core/services/wallet.service';
import { Transaction } from '../../../shared/models';

@Component({
    selector: 'app-transactions',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './transactions.component.html',
    styleUrl: './transactions.component.scss'
})
export class TransactionsComponent implements OnInit {
    transactions = signal<Transaction[]>([]);
    loading = signal(true);
    currentPage = signal(1);
    totalPages = signal(1);
    total = signal(0);

    constructor(private walletService: WalletService) { }

    ngOnInit(): void {
        this.loadTransactions();
    }

    loadTransactions(): void {
        this.loading.set(true);
        this.walletService.getTransactions(this.currentPage(), 10).subscribe({
            next: (res) => {
                if (res.success) {
                    this.transactions.set(res.data.transactions || []);
                    this.totalPages.set(res.data.total_pages);
                    this.total.set(res.data.total);
                }
                this.loading.set(false);
            },
            error: () => {
                this.loading.set(false);
            }
        });
    }

    prevPage(): void {
        if (this.currentPage() > 1) {
            this.currentPage.update(p => p - 1);
            this.loadTransactions();
        }
    }

    nextPage(): void {
        if (this.currentPage() < this.totalPages()) {
            this.currentPage.update(p => p + 1);
            this.loadTransactions();
        }
    }
}
