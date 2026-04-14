import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { BillService } from './bill.service';
import { environment } from '../../../environments/environment';

describe('BillService', () => {
    let service: BillService;
    let httpMock: HttpTestingController;
    const apiUrl = `${environment.apiUrl}/bills`;

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                BillService
            ]
        });

        service = TestBed.inject(BillService);
        httpMock = TestBed.inject(HttpTestingController);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    // --- getCategories ---
    it('getCategories should GET /bills/categories', () => {
        const mockCategories = ['electricity', 'internet', 'phone'];
        service.getCategories().subscribe(res => {
            expect(res.data).toEqual(mockCategories);
        });

        const req = httpMock.expectOne(`${apiUrl}/categories`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: mockCategories });
    });

    // --- getBillers without category ---
    it('getBillers should GET /bills/billers without category', () => {
        service.getBillers().subscribe(res => {
            expect(res.data.length).toBe(2);
        });

        const req = httpMock.expectOne(`${apiUrl}/billers`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: [{ id: '1' }, { id: '2' }] });
    });

    // --- getBillers with category ---
    it('getBillers should include category query param', () => {
        service.getBillers('electricity').subscribe(res => {
            expect(res.data.length).toBe(1);
        });

        const req = httpMock.expectOne(`${apiUrl}/billers?category=electricity`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: [{ id: '1', category: 'electricity' }] });
    });

    // --- payBill ---
    it('payBill should POST to /bills/pay', () => {
        const payData = { biller_id: 'b1', account_number: '12345', amount: 50, save_biller: true };
        service.payBill(payData).subscribe(res => {
            expect(res.data.payment_id).toBe('p1');
        });

        const req = httpMock.expectOne(`${apiUrl}/pay`);
        expect(req.request.method).toBe('POST');
        expect(req.request.body).toEqual(payData);
        req.flush({ success: true, message: 'paid', data: { payment_id: 'p1', amount: '50', new_balance: '950' } });
    });

    // --- getSavedBillers ---
    it('getSavedBillers should GET /bills/saved', () => {
        service.getSavedBillers().subscribe(res => {
            expect(res.data.length).toBe(1);
        });

        const req = httpMock.expectOne(`${apiUrl}/saved`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: [{ id: 'sb1', biller: { name: 'AT&T' } }] });
    });

    // --- removeSavedBiller ---
    it('removeSavedBiller should DELETE /bills/saved/:id', () => {
        service.removeSavedBiller('sb1').subscribe(res => {
            expect(res.success).toBeTrue();
        });

        const req = httpMock.expectOne(`${apiUrl}/saved/sb1`);
        expect(req.request.method).toBe('DELETE');
        req.flush({ success: true, message: 'removed', data: null });
    });
});
