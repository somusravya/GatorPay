import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { TransferService } from './transfer.service';
import { environment } from '../../../environments/environment';

describe('TransferService', () => {
    let service: TransferService;
    let httpMock: HttpTestingController;
    const apiUrl = `${environment.apiUrl}/transfer`;

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                TransferService
            ]
        });

        service = TestBed.inject(TransferService);
        httpMock = TestBed.inject(HttpTestingController);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    // --- sendMoney ---
    it('sendMoney should POST to /transfer/send', () => {
        const transferData = { recipient: 'john', amount: 25.50, note: 'Coffee' };
        service.sendMoney(transferData).subscribe(res => {
            expect(res.data.transaction_id).toBe('tx1');
        });

        const req = httpMock.expectOne(`${apiUrl}/send`);
        expect(req.request.method).toBe('POST');
        expect(req.request.body).toEqual(transferData);
        req.flush({
            success: true,
            message: 'Sent',
            data: { transaction_id: 'tx1', amount: '25.50', new_balance: '974.50', note: 'Coffee', recipient: {} }
        });
    });

    // --- getContacts ---
    it('getContacts should GET /transfer/contacts', () => {
        service.getContacts().subscribe(res => {
            expect(res.data.length).toBe(2);
        });

        const req = httpMock.expectOne(`${apiUrl}/contacts`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: [{ id: 'u1', username: 'alice' }, { id: 'u2', username: 'bob' }] });
    });

    // --- searchUsers ---
    it('searchUsers should GET /transfer/search with query', () => {
        service.searchUsers('alice').subscribe(res => {
            expect(res.data.length).toBe(1);
        });

        const req = httpMock.expectOne(`${apiUrl}/search?query=alice`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: [{ id: 'u1', username: 'alice' }] });
    });

    it('searchUsers should encode special characters in query', () => {
        service.searchUsers('user name').subscribe();

        const req = httpMock.expectOne(`${apiUrl}/search?query=user%20name`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: [] });
    });
});
