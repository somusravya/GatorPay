import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { WalletService } from './wallet.service';
import { AuthService } from './auth.service';
import { environment } from '../../../environments/environment';

describe('WalletService', () => {
    let service: WalletService;
    let httpMock: HttpTestingController;
    let authServiceSpy: jasmine.SpyObj<AuthService>;
    const apiUrl = `${environment.apiUrl}/wallet`;

    beforeEach(() => {
        authServiceSpy = jasmine.createSpyObj('AuthService', ['getMe', 'getToken'], {
            currentUser: jasmine.createSpy().and.returnValue(null),
            currentWallet: jasmine.createSpy().and.returnValue(null),
            isAuthenticated: jasmine.createSpy().and.returnValue(false)
        });
        authServiceSpy.getToken.and.returnValue(null);

        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                WalletService,
                { provide: AuthService, useValue: authServiceSpy },
                { provide: Router, useValue: jasmine.createSpyObj('Router', ['navigate']) }
            ]
        });

        service = TestBed.inject(WalletService);
        httpMock = TestBed.inject(HttpTestingController);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    // --- addMoney ---
    it('addMoney should POST to /wallet/add', () => {
        const addData = { amount: 100, source: 'bank', description: 'Test deposit' };
        service.addMoney(addData).subscribe(res => {
            expect(res.data.balance).toBe('1100.00');
        });

        const req = httpMock.expectOne(`${apiUrl}/add`);
        expect(req.request.method).toBe('POST');
        expect(req.request.body).toEqual(addData);
        req.flush({ success: true, message: 'Added', data: { id: 'w1', balance: '1100.00', currency: 'USD' } });
    });

    // --- withdraw ---
    it('withdraw should POST to /wallet/withdraw', () => {
        const withdrawData = { amount: 50, bank_account: '****1234' };
        service.withdraw(withdrawData).subscribe(res => {
            expect(res.data.balance).toBe('950.00');
        });

        const req = httpMock.expectOne(`${apiUrl}/withdraw`);
        expect(req.request.method).toBe('POST');
        expect(req.request.body).toEqual(withdrawData);
        req.flush({ success: true, message: 'Withdrawn', data: { id: 'w1', balance: '950.00', currency: 'USD' } });
    });

    // --- getTransactions ---
    it('getTransactions should GET /wallet/transactions with pagination', () => {
        service.getTransactions(1, 10).subscribe(res => {
            expect(res.data.transactions.length).toBe(2);
            expect(res.data.total).toBe(2);
        });

        const req = httpMock.expectOne(`${apiUrl}/transactions?page=1&limit=10`);
        expect(req.request.method).toBe('GET');
        req.flush({
            success: true,
            message: 'ok',
            data: {
                transactions: [{ id: 'tx1' }, { id: 'tx2' }],
                total: 2, page: 1, limit: 10, total_pages: 1
            }
        });
    });

    it('getTransactions should use defaults when no params', () => {
        service.getTransactions().subscribe();

        const req = httpMock.expectOne(`${apiUrl}/transactions?page=1&limit=10`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: { transactions: [], total: 0, page: 1, limit: 10, total_pages: 0 } });
    });

    // --- refreshWallet ---
    it('refreshWallet should call authService.getMe()', () => {
        service.refreshWallet();
        expect(authServiceSpy.getMe).toHaveBeenCalled();
    });
});
