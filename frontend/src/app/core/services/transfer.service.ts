import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import {
    ApiResponse, TransferRequest, TransferResponse, User
} from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class TransferService {
    private apiUrl = `${environment.apiUrl}/transfer`;

    constructor(private http: HttpClient) { }

    /** Send money to another user */
    sendMoney(data: TransferRequest) {
        return this.http.post<ApiResponse<TransferResponse>>(`${this.apiUrl}/send`, data);
    }

    /** Get recent transfer contacts */
    getContacts() {
        return this.http.get<ApiResponse<User[]>>(`${this.apiUrl}/contacts`);
    }

    /** Search for users by query */
    searchUsers(query: string) {
        return this.http.get<ApiResponse<User[]>>(`${this.apiUrl}/search?query=${encodeURIComponent(query)}`);
    }
}
