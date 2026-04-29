import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MerchantService } from '../../core/services/merchant.service';

@Component({
    selector: 'app-merchant-portal',
    standalone: true,
    imports: [CommonModule, ReactiveFormsModule],
    templateUrl: './merchant-portal.component.html',
    styleUrl: './merchant-portal.component.scss'
})
export class MerchantPortalComponent implements OnInit {
    private merchantService = inject(MerchantService);
    private fb = inject(FormBuilder);

    invoices = signal<any[]>([]);
    paymentLinks = signal<any[]>([]);
    loading = signal(true);
    activeTab = 'invoices';
    showInvoiceForm = false;
    showLinkForm = false;

    tabs = [
        { id: 'invoices', label: 'Invoices', icon: '📄' },
        { id: 'links', label: 'Payment Links', icon: '🔗' },
        { id: 'analytics', label: 'Analytics', icon: '📊' }
    ];

    invoiceForm: FormGroup;
    linkForm: FormGroup;

    stats = { totalRevenue: 0, totalTransactions: 0, qrPayments: 0, avgOrderValue: 0 };

    constructor() {
        this.invoiceForm = this.fb.group({
            customer_name: ['', Validators.required],
            customer_email: ['', [Validators.required, Validators.email]],
            amount: [0, [Validators.required, Validators.min(0.01)]],
            description: [''],
            due_date: ['']
        });
        this.linkForm = this.fb.group({
            amount: [0, [Validators.required, Validators.min(0.01)]],
            description: [''],
            max_uses: [0]
        });
    }

    ngOnInit(): void {
        this.loadInvoices();
        this.loadPaymentLinks();
        this.calculateStats();
    }

    loadInvoices(): void {
        this.merchantService.getInvoices().subscribe({
            next: (res: any) => {
                if (res.success) this.invoices.set(res.data || []);
                this.loading.set(false);
                this.calculateStats();
            },
            error: () => this.loading.set(false)
        });
    }

    loadPaymentLinks(): void {
        this.merchantService.getPaymentLinks().subscribe({
            next: (res: any) => { if (res.success) this.paymentLinks.set(res.data || []); }
        });
    }

    createInvoice(): void {
        if (this.invoiceForm.invalid) return;
        this.merchantService.createInvoice(this.invoiceForm.value).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.invoices.update(l => [res.data, ...l]);
                    this.showInvoiceForm = false;
                    this.invoiceForm.reset();
                    this.calculateStats();
                }
            }
        });
    }

    createPaymentLink(): void {
        if (this.linkForm.invalid) return;
        this.merchantService.createPaymentLink(this.linkForm.value).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.paymentLinks.update(l => [res.data, ...l]);
                    this.showLinkForm = false;
                    this.linkForm.reset();
                }
            }
        });
    }

    calculateStats(): void {
        const invs = this.invoices();
        const paid = invs.filter(i => i.status === 'paid');
        this.stats.totalRevenue = paid.reduce((s, i) => s + parseFloat(i.amount || 0), 0);
        this.stats.totalTransactions = paid.length;
        this.stats.avgOrderValue = paid.length ? this.stats.totalRevenue / paid.length : 0;
        this.stats.qrPayments = Math.floor(paid.length * 0.35);
    }

    getStatusColor(status: string): string {
        const colors: { [key: string]: string } = {
            pending: '#f59e0b', paid: '#22c55e', partial: '#06b6d4', overdue: '#ef4444', cancelled: '#64748b'
        };
        return colors[status] || '#94a3b8';
    }

    copyLink(linkCode: string): void {
        navigator.clipboard.writeText(`https://gatorpay.com/pay/${linkCode}`);
    }
}
