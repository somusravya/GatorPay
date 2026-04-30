import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HelpCenterService, SupportTicket } from '../../core/services/help-center.service';

@Component({
    selector: 'app-help-center',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './help-center.component.html',
    styleUrl: './help-center.component.scss'
})
export class HelpCenterComponent implements OnInit {
    private helpCenterService = inject(HelpCenterService);

    articles = signal<any[]>([]);
    tickets = signal<SupportTicket[]>([]);
    query = '';
    ticketForm = {
        subject: '',
        category: 'payments',
        priority: 'medium',
        message: ''
    };

    ngOnInit(): void {
        this.helpCenterService.getArticles().subscribe(articles => this.articles.set(articles));
        this.helpCenterService.getTickets().subscribe(tickets => this.tickets.set(tickets));
    }

    filteredArticles(): any[] {
        const value = this.query.trim().toLowerCase();
        if (!value) return this.articles();

        return this.articles().filter(article => {
            return article.title.toLowerCase().includes(value)
                || article.category.toLowerCase().includes(value)
                || article.summary.toLowerCase().includes(value);
        });
    }

    submitTicket(): void {
        if (!this.ticketForm.subject.trim() || !this.ticketForm.message.trim()) return;

        this.helpCenterService.createTicket(this.ticketForm).subscribe(ticket => {
            this.tickets.update(tickets => [ticket, ...tickets]);
            this.ticketForm = {
                subject: '',
                category: 'payments',
                priority: 'medium',
                message: ''
            };
        });
    }
}
