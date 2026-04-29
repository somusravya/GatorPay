import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class MerchantService {
    private apiUrl = `${environment.apiUrl}/merchant`;

    constructor(private http: HttpClient) { }

    getInvoices() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/invoices`);
    }

    createInvoice(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/invoices`, data);
    }

    createPaymentLink(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/payment-links`, data);
    }

    getPaymentLinks() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/payment-links`);
    }
}
