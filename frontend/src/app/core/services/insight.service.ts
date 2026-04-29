import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class InsightService {
    private apiUrl = `${environment.apiUrl}/insights`;

    constructor(private http: HttpClient) { }

    getSummary() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/summary`);
    }
}
