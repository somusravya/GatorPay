import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LoansComponent } from './loans.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { LoanService } from '../../core/services/loan.service';
import { of } from 'rxjs';

describe('LoansComponent', () => {
  let component: LoansComponent;
  let fixture: ComponentFixture<LoansComponent>;

  beforeEach(async () => {
    const mockLoanService = {
      getUserLoans: jasmine.createSpy().and.returnValue(of({ success: true, data: [] })),
      getOffers: jasmine.createSpy().and.returnValue(of({ success: true, data: [] })),
      applyForLoan: jasmine.createSpy().and.returnValue(of({ success: true })),
      payEMI: jasmine.createSpy().and.returnValue(of({ success: true }))
    };

    await TestBed.configureTestingModule({
      imports: [LoansComponent, HttpClientTestingModule],
      providers: [
        { provide: LoanService, useValue: mockLoanService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(LoansComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize active tab implicitly', () => {
    expect(component.activeTab).toBe('active');
  });
});
