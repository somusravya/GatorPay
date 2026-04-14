import { ComponentFixture, TestBed } from '@angular/core/testing';
import { TradingComponent } from './trading.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { TradingService } from '../../core/services/trading.service';
import { StockService } from '../../core/services/stock.service';
import { of } from 'rxjs';

describe('TradingComponent', () => {
  let component: TradingComponent;
  let fixture: ComponentFixture<TradingComponent>;

  beforeEach(async () => {
    const mockTradingService = {
      getAccountStatus: jasmine.createSpy().and.returnValue(of({ success: true, data: { status: 'verified', buying_power: 10000 } })),
      verifyAccount: jasmine.createSpy().and.returnValue(of({ success: true })),
      executeTrade: jasmine.createSpy().and.returnValue(of({ success: true })),
      getPortfolio: jasmine.createSpy().and.returnValue(of({ success: true, data: { positions: [], buying_power: 10000 } }))
    };

    const mockStockService = {
      getQuote: jasmine.createSpy().and.returnValue(of({ success: true, data: { current_price: 150 } })),
      getChart: jasmine.createSpy().and.returnValue(of({ success: true, data: [] }))
    };

    await TestBed.configureTestingModule({
      imports: [TradingComponent, HttpClientTestingModule],
      providers: [
        { provide: TradingService, useValue: mockTradingService },
        { provide: StockService, useValue: mockStockService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(TradingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should mock check verify state on init', () => {
     expect(component.isVerified).toBeFalse(); 
  });
});
