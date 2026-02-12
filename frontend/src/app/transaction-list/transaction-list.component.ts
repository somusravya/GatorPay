import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

interface Transaction {
  id: number;
  payer: string;
  amount: number;
  date: string;
}

@Component({
  selector: 'app-transaction-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './transaction-list.component.html',
  styleUrls: ['./transaction-list.component.css'],
})
export class TransactionListComponent implements OnInit {
  transactions: Transaction[] = [];
  newTransaction = { payer: '', amount: 0, date: '' };
  apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadTransactions();
  }

  loadTransactions() {
    this.http.get<{ transactions: Transaction[] }>(`${this.apiUrl}/transactions`)
      .subscribe({
        next: (data) => {
          this.transactions = data.transactions;
        },
        error: (err) => {
          console.error('Error loading transactions:', err);
          // Show sample data if API fails
          this.transactions = [
            { id: 1, payer: 'Alice', amount: 50, date: '2026-02-12' },
            { id: 2, payer: 'Bob', amount: 30, date: '2026-02-11' },
          ];
        },
      });
  }

  addTransaction() {
    if (!this.newTransaction.payer || this.newTransaction.amount <= 0) {
      alert('Please fill in all fields');
      return;
    }

    this.http.post<Transaction>(`${this.apiUrl}/transactions`, this.newTransaction)
      .subscribe({
        next: (transaction) => {
          this.transactions.push(transaction);
          this.newTransaction = { payer: '', amount: 0, date: '' };
        },
        error: (err) => {
          console.error('Error adding transaction:', err);
          alert('Failed to add transaction');
        },
      });
  }

  deleteTransaction(id: number) {
    this.transactions = this.transactions.filter(t => t.id !== id);
  }

  getTotalAmount(): number {
    return this.transactions.reduce((sum, t) => sum + t.amount, 0);
  }
}
