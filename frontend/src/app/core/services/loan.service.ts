import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({
  providedIn: 'root'
})
export class LoanService {
  private apiUrl = `${environment.apiUrl}/loans`;

  constructor(private http: HttpClient) {}

  getOffers(): Observable<ApiResponse<any>> {
    return this.http.get<ApiResponse<any>>(`${this.apiUrl}/offers`);
  }

  applyForLoan(payload: { offer_id: string, amount: number, term_months: number }): Observable<ApiResponse<any>> {
    return this.http.post<ApiResponse<any>>(`${this.apiUrl}/apply`, payload);
  }

  getUserLoans(): Observable<ApiResponse<any>> {
    return this.http.get<ApiResponse<any>>(`${this.apiUrl}`);
  }

  getLoan(id: string): Observable<ApiResponse<any>> {
    return this.http.get<ApiResponse<any>>(`${this.apiUrl}/${id}`);
  }

  payEMI(id: string): Observable<ApiResponse<any>> {
    return this.http.post<ApiResponse<any>>(`${this.apiUrl}/${id}/pay`, {});
  }
}
