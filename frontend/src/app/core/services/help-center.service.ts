import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';

export interface SupportTicket {
    id: string;
    subject: string;
    category: string;
    priority: string;
    status: string;
    message: string;
    created_at: string;
}

@Injectable({ providedIn: 'root' })
export class HelpCenterService {
    private readonly storageKey = 'gatorpay_support_tickets';

    getArticles() {
        return of([
            {
                title: 'Account verification',
                category: 'Account',
                summary: 'Check KYC status, upload missing details, and understand review timing.',
                icon: '🛡️'
            },
            {
                title: 'Transfers and limits',
                category: 'Payments',
                summary: 'Review transfer timing, failed payments, and daily limit behavior.',
                icon: '💸'
            },
            {
                title: 'Cards and security',
                category: 'Cards',
                summary: 'Freeze cards, report suspicious activity, and manage card settings.',
                icon: '💳'
            },
            {
                title: 'Rewards and cashback',
                category: 'Rewards',
                summary: 'Learn how points, cashback, and merchant offers are calculated.',
                icon: '🎁'
            }
        ]);
    }

    getTickets(): Observable<SupportTicket[]> {
        return of(this.readTickets());
    }

    createTicket(ticket: Pick<SupportTicket, 'subject' | 'category' | 'priority' | 'message'>): Observable<SupportTicket> {
        const tickets = this.readTickets();
        const created: SupportTicket = {
            ...ticket,
            id: globalThis.crypto?.randomUUID?.() ?? Date.now().toString(),
            status: 'open',
            created_at: new Date().toISOString()
        };

        localStorage.setItem(this.storageKey, JSON.stringify([created, ...tickets]));
        return of(created);
    }

    private readTickets(): SupportTicket[] {
        const raw = localStorage.getItem(this.storageKey);
        if (!raw) return [];

        try {
            return JSON.parse(raw) as SupportTicket[];
        } catch {
            return [];
        }
    }
}
