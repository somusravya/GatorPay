import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import {
    ApiResponse, RewardSummary, RewardHistoryResponse, Offer
} from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class RewardService {
    private apiUrl = `${environment.apiUrl}/rewards`;

    constructor(private http: HttpClient) { }

    /** Get reward summary */
    getSummary() {
        return this.http.get<ApiResponse<RewardSummary>>(`${this.apiUrl}`);
    }

    /** Get paginated reward history */
    getHistory(page: number = 1, limit: number = 10) {
        return this.http.get<ApiResponse<RewardHistoryResponse>>(
            `${this.apiUrl}/history?page=${page}&limit=${limit}`
        );
    }

    /** Get available offers */
    getOffers() {
        return this.http.get<ApiResponse<Offer[]>>(`${this.apiUrl}/offers`);
    }
}
