import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import {
    ApiResponse, Biller, SavedBiller, BillPayRequest, BillPayResponse
} from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class BillService {
    private apiUrl = `${environment.apiUrl}/bills`;

    constructor(private http: HttpClient) { }

    /** Get bill categories */
    getCategories() {
        return this.http.get<ApiResponse<string[]>>(`${this.apiUrl}/categories`);
    }

    /** Get billers, optionally filtered by category */
    getBillers(category?: string) {
        const url = category
            ? `${this.apiUrl}/billers?category=${encodeURIComponent(category)}`
            : `${this.apiUrl}/billers`;
        return this.http.get<ApiResponse<Biller[]>>(url);
    }

    /** Pay a bill */
    payBill(data: BillPayRequest) {
        return this.http.post<ApiResponse<BillPayResponse>>(`${this.apiUrl}/pay`, data);
    }

    /** Get saved billers */
    getSavedBillers() {
        return this.http.get<ApiResponse<SavedBiller[]>>(`${this.apiUrl}/saved`);
    }

    /** Remove a saved biller */
    removeSavedBiller(id: string) {
        return this.http.delete<ApiResponse<null>>(`${this.apiUrl}/saved/${id}`);
    }
}
