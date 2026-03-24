import { Component, signal, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { BillService } from '../../core/services/bill.service';
import { AuthService } from '../../core/services/auth.service';
import { Biller, SavedBiller } from '../../shared/models';

@Component({
    selector: 'app-bills',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './bills.component.html',
    styleUrl: './bills.component.scss'
})
export class BillsComponent implements OnInit {
    private billService = inject(BillService);
    private authService = inject(AuthService);

    categories = signal<string[]>([]);
    billers = signal<Biller[]>([]);
    savedBillers = signal<SavedBiller[]>([]);
    selectedCategory = signal<string>('');
    selectedBiller = signal<Biller | null>(null);
    accountNumber = '';
    amount = 0;
    saveBiller = false;
    loading = signal(false);
    showSuccess = signal(false);
    successMessage = signal('');
    errorMessage = signal('');

    ngOnInit(): void {
        this.loadCategories();
        this.loadBillers();
        this.loadSavedBillers();
    }

    loadCategories(): void {
        this.billService.getCategories().subscribe({
            next: (res) => {
                if (res.success) {
                    this.categories.set(res.data || []);
                }
            }
        });
    }

    loadBillers(category?: string): void {
        this.billService.getBillers(category).subscribe({
            next: (res) => {
                if (res.success) {
                    this.billers.set(res.data || []);
                }
            }
        });
    }

    loadSavedBillers(): void {
        this.billService.getSavedBillers().subscribe({
            next: (res) => {
                if (res.success) {
                    this.savedBillers.set(res.data || []);
                }
            }
        });
    }

    selectCategory(category: string): void {
        this.selectedCategory.set(category);
        this.loadBillers(category || undefined);
    }

    selectBiller(biller: Biller): void {
        this.selectedBiller.set(biller);
        this.accountNumber = '';
        this.amount = 0;
        this.saveBiller = false;
    }

    selectSavedBiller(saved: SavedBiller): void {
        this.selectedBiller.set(saved.biller);
        this.accountNumber = saved.account_number;
        this.amount = 0;
    }

    clearBiller(): void {
        this.selectedBiller.set(null);
        this.accountNumber = '';
        this.amount = 0;
    }

    payBill(): void {
        const biller = this.selectedBiller();
        if (!biller || !this.accountNumber || this.amount <= 0) return;

        this.loading.set(true);
        this.errorMessage.set('');

        this.billService.payBill({
            biller_id: biller.id,
            account_number: this.accountNumber,
            amount: this.amount,
            save_biller: this.saveBiller
        }).subscribe({
            next: (res) => {
                if (res.success) {
                    this.successMessage.set(`$${this.amount.toFixed(2)} paid to ${biller.name}`);
                    this.showSuccess.set(true);
                    this.clearBiller();
                    this.authService.getMe();
                    this.loadSavedBillers();
                }
                this.loading.set(false);
            },
            error: (err) => {
                this.errorMessage.set(err.error?.message || 'Payment failed');
                this.loading.set(false);
            }
        });
    }

    removeSavedBiller(id: string): void {
        this.billService.removeSavedBiller(id).subscribe({
            next: (res) => {
                if (res.success) {
                    this.loadSavedBillers();
                }
            }
        });
    }

    closeSuccess(): void {
        this.showSuccess.set(false);
    }
}
