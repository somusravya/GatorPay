import { Component, OnInit, ViewChild, ElementRef, AfterViewInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { createChart, LineSeries, Time } from 'lightweight-charts';
import { TradingService } from '../../core/services/trading.service';
import { StockService } from '../../core/services/stock.service';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-trading',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './trading.component.html',
  styleUrls: ['./trading.component.scss']
})
export class TradingComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('chartContainer') chartContainer!: ElementRef;
  private chart!: any;
  private lineSeries!: any;
  Math = Math;
  
  
  isVerified = false;
  showKycModal = false;
  kycForm: FormGroup;
  
  marketSummary: any[] = [];
  portfolio: any = null;
  buyingPower: number = 0;
  totalPortfolioValue: number = 0;
  
  // Search & Active Stock
  searchForm: FormGroup;
  searchResults: any[] = [];
  activeStock: any = null;
  activeStockDetails: any = null;
  
  // Trade Panel
  tradeForm: FormGroup;
  tradeType: 'buy' | 'sell' = 'buy';
  orderProcessing = false;

  watchlist: string[] = ['AAPL', 'TSLA', 'GOOGL', 'AMZN', 'MSFT'];
  watchlistData: any[] = [];

  // Live simulation
  private simulationInterval: any = null;
  priceFlashState: 'up' | 'down' | null = null;
  flashTimeout: any = null;

  constructor(
    private fb: FormBuilder,
    private tradingService: TradingService,
    private stockService: StockService
  ) {
    this.kycForm = this.fb.group({
      dob: ['', Validators.required],
      ssn: ['', [Validators.required, Validators.pattern('^[0-9]{9}$')]],
      sec_q1: ['', Validators.required],
      sec_a1: ['', Validators.required],
      sec_q2: ['', Validators.required],
      sec_a2: ['', Validators.required],
      risk_ack: [false, Validators.requiredTrue]
    });

    this.searchForm = this.fb.group({
      query: ['']
    });

    this.tradeForm = this.fb.group({
      quantity: [1, [Validators.required, Validators.min(1)]]
    });
  }

  async ngOnInit() {
    await this.checkKycStatus();
    await this.loadDashboard();
    this.startLiveSimulation();
  }

  ngAfterViewInit() {
    this.initChart();
  }

  ngOnDestroy() {
    if (this.chart) {
      this.chart.remove();
    }
    if (this.simulationInterval) {
      clearInterval(this.simulationInterval);
    }
    if (this.flashTimeout) {
      clearTimeout(this.flashTimeout);
    }
  }

  async checkKycStatus() {
    try {
      const res = await firstValueFrom(this.tradingService.getAccount());
      if (res.data && res.data.status === 'verified') {
        this.isVerified = true;
      }
    } catch (e) {
      this.isVerified = false;
    }
  }

  async loadDashboard() {
    if (this.isVerified) {
       await this.loadPortfolio();
    }
    try {
      const marketRes = await firstValueFrom(this.stockService.getMarketSummary());
      if (marketRes.data) this.marketSummary = marketRes.data;
    } catch(e) {}
    
    await this.loadWatchlistData();
    // Load default stock
    await this.selectStock('AAPL');
  }

  async loadPortfolio() {
    if (!this.isVerified) return;
    try {
      const res = await firstValueFrom(this.tradingService.getPortfolio());
      if (res.data) {
        this.portfolio = res.data.positions || [];
        this.buyingPower = res.data.buying_power || 0;
        this.calculatePortfolioValue();
      }
    } catch(e) {
      console.error("Failed to load portfolio", e);
    }
  }

  calculatePortfolioValue() {
    // Basic calculation of total portfolio value using avg costs if real prices aren't available yet
    let posValue = 0;
    if (this.portfolio) {
      for(let p of this.portfolio) {
        // if active stock is viewing p, use active price
        const currentPrice = this.activeStock && this.activeStock.symbol === p.symbol ? this.activeStock.price : p.avg_cost;
        posValue += (p.quantity * currentPrice);
      }
    }
    this.totalPortfolioValue = posValue + this.buyingPower;
  }

  async loadWatchlistData() {
    this.watchlistData = [];
    for (const sym of this.watchlist) {
      try {
        const q = await firstValueFrom(this.stockService.getQuote(sym));
        if(q.data) {
          this.watchlistData.push({
            symbol: sym,
            price: q.data.price,
            change: q.data.change_percent,
            flash: null
          });
        }
      } catch(e) {}
    }
  }

  async submitKyc() {
    if (this.kycForm.invalid) return;
    try {
      await firstValueFrom(this.tradingService.verifyAccount(this.kycForm.value));
      this.isVerified = true;
      this.showKycModal = false;
      await this.loadPortfolio(); // load portfolio once verified
    } catch (e) {
      alert("KYC Failed: " + JSON.stringify(e));
    }
  }

  async searchStock() {
    const query = this.searchForm.value.query;
    if (!query) return;
    try {
      const res = await firstValueFrom(this.stockService.search(query));
      this.searchResults = res.data || [];
    } catch(e) {}
  }

  async selectStock(symbol: string) {
    this.searchResults = [];
    this.searchForm.patchValue({query: ''});
    
    try {
      // Load quote
      const quoteRes = await firstValueFrom(this.stockService.getQuote(symbol));
      this.activeStock = quoteRes.data;

      // Ensure chart gets re-initialized if destroyed
      if(!this.chart) {
        this.initChart();
      }

      this.stockService.getDetails(symbol).subscribe(res => {
         this.activeStockDetails = res.data;
      });

      // Load chart
      const chartRes = await firstValueFrom(this.stockService.getChart(symbol));
      if (chartRes.data && this.lineSeries) {
        const d = chartRes.data as any[];
        // Transform data slightly to assure correct TS time format if needed
        const formatted = d.map(item => ({...item, time: this.formatTime(item.time)}));
        this.lineSeries.setData(formatted);
      }
      
      this.tradeForm.reset({quantity: 1});
      this.tradeType = 'buy';
    } catch(e) {
      console.warn("Failed to select stock", e);
    }
  }

  formatTime(timeStr: string | number): Time {
    // If it's already an epoch timestamp or lightweight string
    if (typeof timeStr === 'number') return timeStr as Time;
    return (new Date(timeStr).getTime() / 1000) as Time;
  }

  initChart() {
    if (!this.chartContainer) return;
    this.chart = createChart(this.chartContainer.nativeElement, {
      width: this.chartContainer.nativeElement.clientWidth,
      height: 480,
      layout: {
        background: { color: 'transparent' },
        textColor: '#94a3b8',
      },
      grid: {
        vertLines: { color: 'rgba(31, 41, 55, 0.4)', style: 1 },
        horzLines: { color: 'rgba(31, 41, 55, 0.4)', style: 1 },
      },
      timeScale: {
        borderColor: '#1f2937',
        timeVisible: true,
        secondsVisible: false,
      },
      rightPriceScale: {
        borderColor: '#1f2937',
      },
      crosshair: {
        mode: 0,
      }
    });

    this.lineSeries = this.chart.addSeries(LineSeries, {
      color: '#3b82f6',
      lineWidth: 2,
      crosshairMarkerVisible: true,
      crosshairMarkerRadius: 5,
      crosshairMarkerBorderColor: '#60a5fa',
      crosshairMarkerBackgroundColor: '#1d4ed8'
    });

    // Handle responsive resize
    new ResizeObserver(entries => {
      if (entries.length === 0 || entries[0].target !== this.chartContainer.nativeElement) { return; }
      const newRect = entries[0].contentRect;
      this.chart.applyOptions({ height: newRect.height, width: newRect.width });
    }).observe(this.chartContainer.nativeElement);
  }

  setTradeType(type: 'buy' | 'sell') {
    this.tradeType = type;
  }

  async executeTrade() {
    if (this.tradeForm.invalid || !this.activeStock || this.orderProcessing) return;
    
    // Popup KYC if not verified!
    if (!this.isVerified) {
       this.showKycModal = true;
       return;
    }
    
    this.orderProcessing = true;
    const payload = {
      symbol: this.activeStock.symbol,
      type: this.tradeType,
      quantity: this.tradeForm.value.quantity,
      price: this.activeStock.price
    };

    try {
      await firstValueFrom(this.tradingService.executeTrade(payload));
      // Show success animation overlay eventually
      this.tradeForm.reset({quantity: 1});
      await this.loadPortfolio();
    } catch (e: any) {
      console.error(e);
      alert("Trade failed: Insufficient funds or invalid stock.");
    } finally {
      this.orderProcessing = false;
    }
  }

  addToWatchlist(symbol: string) {
    if (!this.watchlist.includes(symbol)) {
      this.watchlist.push(symbol);
      this.loadWatchlistData();
    }
  }

  // --- Live Data Simulator ---
  startLiveSimulation() {
    this.simulationInterval = setInterval(() => {
      this.simulateMarketTick();
    }, 2000); // tick every 2 seconds
  }

  simulateMarketTick() {
    // Randomize active stock
    if (this.activeStock) {
      const volatility = 0.001; // max 0.1% move
      const move = this.activeStock.price * volatility * (Math.random() - 0.48); // slightly biased to go up over time optionally, or neutral
      const newPrice = this.activeStock.price + move;

      // Trigger flash
      if(this.flashTimeout) clearTimeout(this.flashTimeout);
      this.priceFlashState = move >= 0 ? 'up' : 'down';
      this.flashTimeout = setTimeout(() => { this.priceFlashState = null; }, 1000);

      this.activeStock.price = newPrice;
      
      // Add data point to chart
      if (this.lineSeries) {
         try {
           this.lineSeries.update({
             time: (Date.now() / 1000) as Time,
             value: newPrice
           });
         } catch(e) {} // ignore if time is slightly off
      }
    }

    // Randomize market summary
    if (this.marketSummary) {
      this.marketSummary.forEach(m => {
        if(Math.random() > 0.5) {
          const move = (Math.random() - 0.5) * 10;
          m.price += move;
          m.change_percent += (move / 100);
        }
      });
    }

    // Randomize watchlist
    if (this.watchlistData) {
      this.watchlistData.forEach(w => {
         if(Math.random() > 0.3) {
            const move = w.price * 0.001 * (Math.random() - 0.5);
            w.price += move;
            w.change += (move / w.price) * 100;
            w.flash = move >= 0 ? 'up' : 'down';
            setTimeout(() => { w.flash = null; }, 1000);
         }
      });
    }
    
    this.calculatePortfolioValue();
  }
}
