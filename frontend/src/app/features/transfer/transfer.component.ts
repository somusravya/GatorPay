import { Component, signal, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TransferService } from '../../core/services/transfer.service';
import { AuthService } from '../../core/services/auth.service';
import { User } from '../../shared/models';
import { Subject, debounceTime, distinctUntilChanged } from 'rxjs';

@Component({
    selector: 'app-transfer',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './transfer.component.html',
    styleUrl: './transfer.component.scss'
})
export class TransferComponent implements OnInit {
    private transferService = inject(TransferService);
    private authService = inject(AuthService);

    searchQuery = '';
    searchResults = signal<User[]>([]);
    recentContacts = signal<User[]>([]);
    selectedUser = signal<User | null>(null);
    amount = 0;
    note = '';
    loading = signal(false);
    searchLoading = signal(false);
    showSuccess = signal(false);
    successMessage = signal('');
    errorMessage = signal('');

    private searchSubject = new Subject<string>();

    ngOnInit(): void {
        this.loadRecentContacts();

        // Debounced search
        this.searchSubject.pipe(
            debounceTime(300),
            distinctUntilChanged()
        ).subscribe(query => {
            if (query.length >= 2) {
                this.performSearch(query);
            } else {
                this.searchResults.set([]);
            }
        });
    }

    onSearchInput(): void {
        this.searchSubject.next(this.searchQuery);
    }

    private performSearch(query: string): void {
        this.searchLoading.set(true);
        this.transferService.searchUsers(query).subscribe({
            next: (res) => {
                if (res.success) {
                    this.searchResults.set(res.data || []);
                }
                this.searchLoading.set(false);
            },
            error: () => {
                this.searchLoading.set(false);
            }
        });
    }

    loadRecentContacts(): void {
        this.transferService.getContacts().subscribe({
            next: (res) => {
                if (res.success) {
                    this.recentContacts.set(res.data || []);
                }
            }
        });
    }

    selectUser(user: User): void {
        this.selectedUser.set(user);
        this.searchQuery = '';
        this.searchResults.set([]);
    }

    clearSelection(): void {
        this.selectedUser.set(null);
        this.amount = 0;
        this.note = '';
    }

    sendMoney(): void {
        const recipient = this.selectedUser();
        if (!recipient || this.amount <= 0) return;

        this.loading.set(true);
        this.errorMessage.set('');

        this.transferService.sendMoney({
            recipient: recipient.username,
            amount: this.amount,
            note: this.note
        }).subscribe({
            next: (res) => {
                if (res.success) {
                    this.successMessage.set(`$${this.amount.toFixed(2)} sent to ${recipient.first_name} ${recipient.last_name}`);
                    this.showSuccess.set(true);
                    this.clearSelection();
                    this.authService.getMe(); // refresh wallet
                    this.loadRecentContacts();
                }
                this.loading.set(false);
            },
            error: (err) => {
                this.errorMessage.set(err.error?.message || 'Transfer failed');
                this.loading.set(false);
            }
        });
    }

    closeSuccess(): void {
        this.showSuccess.set(false);
    }

    getUserInitials(user: User): string {
        return (user.first_name[0] + user.last_name[0]).toUpperCase();
    }
}
