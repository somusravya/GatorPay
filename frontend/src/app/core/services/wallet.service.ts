import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import {
    Wallet, Transaction, ApiResponse,
    AddMoneyRequest, WithdrawRequest, TransactionListResponse
} from '../../shared/models';
import { AuthService } from './auth.service';

@Injectable({ providedIn: 'root' })
export class WalletService {
    private apiUrl = `${environment.apiUrl}/wallet`;

    transactions = signal<Transaction[]>([]);
    totalTransactions = signal<number>(0);
    totalPages = signal<number>(0);

    constructor(private http: HttpClient, private authService: AuthService) { }

    /** Add money to wallet */
    addMoney(data: AddMoneyRequest) {
        return this.http.post<ApiResponse<Wallet>>(`${this.apiUrl}/add`, data);
    }

    /** Withdraw money from wallet */
    withdraw(data: WithdrawRequest) {
        return this.http.post<ApiResponse<Wallet>>(`${this.apiUrl}/withdraw`, data);
    }

    /** Get paginated transactions */
    getTransactions(page: number = 1, limit: number = 10) {
        return this.http.get<ApiResponse<TransactionListResponse>>(
            `${this.apiUrl}/transactions?page=${page}&limit=${limit}`
        );
    }

    /** Refresh wallet after operations */
    refreshWallet(): void {
        this.authService.getMe();
    }
}
