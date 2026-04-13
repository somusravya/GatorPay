import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter, Router } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { AuthService } from '../../core/services/auth.service';
import { WalletService } from '../../core/services/wallet.service';
import { of } from 'rxjs';
import { signal } from '@angular/core';

describe('DashboardComponent', () => {
    let component: DashboardComponent;
    let fixture: ComponentFixture<DashboardComponent>;
    let authServiceSpy: jasmine.SpyObj<AuthService>;
    let walletServiceSpy: jasmine.SpyObj<WalletService>;

    beforeEach(async () => {
        authServiceSpy = jasmine.createSpyObj('AuthService', ['getToken', 'getMe'], {
            currentUser: signal({ id: '1', username: 'testuser', first_name: 'Test', last_name: 'User' }),
            currentWallet: signal({ id: 'w1', balance: '1250.75', currency: 'USD' }),
            isAuthenticated: signal(true)
        });
        authServiceSpy.getToken.and.returnValue('token');

        walletServiceSpy = jasmine.createSpyObj('WalletService', ['getTransactions'], {
            transactions: signal([]),
            totalTransactions: signal(0),
            totalPages: signal(0)
        });
        walletServiceSpy.getTransactions.and.returnValue(of({
            success: true,
            message: 'ok',
            data: {
                transactions: [
                    { id: 'tx1', wallet_id: 'w1', from_user_id: 'u1', to_user_id: '', type: 'deposit', amount: '100.00', description: 'Test deposit', status: 'success', created_at: '2025-01-01' },
                    { id: 'tx2', wallet_id: 'w1', from_user_id: 'u1', to_user_id: 'u2', type: 'p2p_send', amount: '25.00', description: 'Transfer', status: 'success', created_at: '2025-01-02' }
                ],
                total: 2, page: 1, limit: 5, total_pages: 1
            }
        }));

        await TestBed.configureTestingModule({
            imports: [DashboardComponent],
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                provideRouter([]),
                { provide: AuthService, useValue: authServiceSpy },
                { provide: WalletService, useValue: walletServiceSpy }
            ]
        }).compileComponents();

        fixture = TestBed.createComponent(DashboardComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    // --- quickActions ---
    it('should have 5 quick action items', () => {
        expect(component.quickActions.length).toBe(5);
    });

    it('quickActions should include expected routes', () => {
        const routes = component.quickActions.map(a => a.route);
        expect(routes).toContain('/wallet');
        expect(routes).toContain('/transfer');
        expect(routes).toContain('/bills');
        expect(routes).toContain('/rewards');
        expect(routes).toContain('/transactions');
    });

    // --- getFormattedBalance ---
    it('getFormattedBalance should format balance as USD currency', () => {
        const formatted = component.getFormattedBalance();
        expect(formatted).toContain('1,250.75');
        expect(formatted).toContain('$');
    });

    it('getFormattedBalance should return $0.00 when no wallet', () => {
        (component as any).wallet = signal(null);
        const formatted = component.getFormattedBalance();
        expect(formatted).toContain('0.00');
    });

    // --- loadRecentTransactions ---
    it('should load recent transactions on init', () => {
        expect(walletServiceSpy.getTransactions).toHaveBeenCalledWith(1, 5);
    });

    it('should set loading to false after transactions load', () => {
        expect(component.loading()).toBeFalse();
    });

    it('should populate recentTransactions signal', () => {
        expect(component.recentTransactions().length).toBe(2);
    });
});
