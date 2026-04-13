import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { LoanService } from '../../core/services/loan.service';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-loans',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './loans.component.html',
  styleUrls: ['./loans.component.scss']
})
export class LoansComponent implements OnInit {
  activeTab: 'active' | 'offers' = 'active';
  
  loans: any[] = [];
  offers: any[] = [];
  
  // Apply Modal
  showApplyModal = false;
  selectedOffer: any = null;
  applyForm: FormGroup;
  calculatedEMI: number = 0;
  
  constructor(private loanService: LoanService, private fb: FormBuilder) {
    this.applyForm = this.fb.group({
       amount: [1000, Validators.required],
       term_months: [12, Validators.required]
    });
    
    this.applyForm.valueChanges.subscribe(() => {
      this.recalcEMI();
    });
  }

  ngOnInit() {
    this.loadLoans();
    this.loadOffers();
  }

  async loadLoans() {
    try {
      const res = await firstValueFrom(this.loanService.getUserLoans());
      this.loans = res.data || [];
    } catch(e) {}
  }

  async loadOffers() {
    try {
      const res = await firstValueFrom(this.loanService.getOffers());
      this.offers = res.data || [];
    } catch(e) {}
  }

  openApplyModal(offer: any) {
    this.selectedOffer = offer;
    this.applyForm.patchValue({
      amount: offer.min_amount,
      term_months: Math.min(12, offer.term_months)
    });
    this.showApplyModal = true;
    this.recalcEMI();
  }

  recalcEMI() {
    if (!this.selectedOffer) return;
    const amount = this.applyForm.value.amount || 0;
    const term = this.applyForm.value.term_months || 1;
    const rate = this.selectedOffer.interest_rate || 0;
    
    if (rate === 0) {
      this.calculatedEMI = amount / term;
      return;
    }
    
    const r = (rate / 100) / 12;
    this.calculatedEMI = amount * r * Math.pow(1+r, term) / (Math.pow(1+r, term) - 1);
  }

  async submitApplication() {
    if (this.applyForm.invalid) return;
    try {
       await firstValueFrom(this.loanService.applyForLoan({
         offer_id: this.selectedOffer.id,
         amount: this.applyForm.value.amount,
         term_months: this.applyForm.value.term_months
       }));
       alert("Loan application approved and funds disbursed to your wallet!");
       this.showApplyModal = false;
       this.activeTab = 'active';
       this.loadLoans();
    } catch (e: any) {
       alert("Application failed: " + e.error?.message);
    }
  }

  async payEMI(loanId: string) {
    if (!confirm("Confirm payment for your next EMI from your wallet balance?")) return;
    try {
      await firstValueFrom(this.loanService.payEMI(loanId));
      alert("EMI Payment successful!");
      this.loadLoans();
    } catch (e: any) {
      alert("Payment failed: " + e.error?.message);
    }
  }
}
