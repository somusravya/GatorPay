import { ComponentFixture, TestBed } from '@angular/core/testing';
import { CardsComponent } from './cards.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { CardService } from '../../core/services/card.service';
import { of } from 'rxjs';

describe('CardsComponent', () => {
  let component: CardsComponent;
  let fixture: ComponentFixture<CardsComponent>;

  beforeEach(async () => {
    const mockCardService = {
      getCards: jasmine.createSpy().and.returnValue(of({ success: true, data: [] })),
      createCard: jasmine.createSpy().and.returnValue(of({ success: true })),
      freezeCard: jasmine.createSpy().and.returnValue(of({ success: true })),
      requestOTP: jasmine.createSpy().and.returnValue(of({ success: true })),
      getCardDetails: jasmine.createSpy().and.returnValue(of({ success: true, data: {} }))
    };

    await TestBed.configureTestingModule({
      imports: [CardsComponent, HttpClientTestingModule],
      providers: [
        { provide: CardService, useValue: mockCardService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(CardsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should format card numbers properly', () => {
    expect(component.formatCardNumber('1234567812345678')).toBe('1234 5678 1234 5678');
  });
});
