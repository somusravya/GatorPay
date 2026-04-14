import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { RewardService } from './reward.service';
import { environment } from '../../../environments/environment';

describe('RewardService', () => {
    let service: RewardService;
    let httpMock: HttpTestingController;
    const apiUrl = `${environment.apiUrl}/rewards`;

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                RewardService
            ]
        });

        service = TestBed.inject(RewardService);
        httpMock = TestBed.inject(HttpTestingController);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    // --- getSummary ---
    it('getSummary should GET /rewards', () => {
        const mockSummary = { total_points: 500, total_cashback: '25.00', lifetime_earnings: '25.00', total_transactions: 5 };
        service.getSummary().subscribe(res => {
            expect(res.data).toEqual(mockSummary);
        });

        const req = httpMock.expectOne(apiUrl);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: mockSummary });
    });

    // --- getHistory ---
    it('getHistory should GET /rewards/history with pagination', () => {
        service.getHistory(2, 20).subscribe(res => {
            expect(res.data.rewards.length).toBe(1);
        });

        const req = httpMock.expectOne(`${apiUrl}/history?page=2&limit=20`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: { rewards: [{ id: 'r1' }], total: 1, page: 2, limit: 20 } });
    });

    it('getHistory should use default page and limit', () => {
        service.getHistory().subscribe();

        const req = httpMock.expectOne(`${apiUrl}/history?page=1&limit=10`);
        expect(req.request.method).toBe('GET');
        req.flush({ success: true, message: 'ok', data: { rewards: [], total: 0, page: 1, limit: 10 } });
    });

    // --- getOffers ---
    it('getOffers should GET /rewards/offers', () => {
        service.getOffers().subscribe(res => {
            expect(res.data.length).toBe(4);
        });

        const req = httpMock.expectOne(`${apiUrl}/offers`);
        expect(req.request.method).toBe('GET');
        req.flush({
            success: true, message: 'ok',
            data: [
                { id: '1', title: 'Send & Save' },
                { id: '2', title: 'Bill Pay Bonus' },
                { id: '3', title: 'Points Multiplier' },
                { id: '4', title: 'Refer a Friend' }
            ]
        });
    });
});
