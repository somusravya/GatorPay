describe('Wallet and Rewards Pages', () => {
    beforeEach(() => {
        // Mock authentication state
        window.localStorage.setItem('gatorpay_token', 'fake-jwt-token');
        
        // Mock user profile response
        cy.intercept('GET', '**/api/v1/auth/me', {
            statusCode: 200,
            body: {
                success: true,
                data: {
                    user: { id: 'user-1', first_name: 'Test', last_name: 'User', email: 'test@example.com' },
                    wallet: { balance: 1540.50, currency: 'USD' }
                }
            }
        }).as('getMe');
    });

    describe('Wallet Page', () => {
        beforeEach(() => {
            cy.intercept('GET', '**/api/v1/wallet/transactions*', {
                statusCode: 200,
                body: { success: true, data: { transactions: [], total: 0 } }
            }).as('getTransactions');
            
            cy.visit('/wallet');
        });

        it('should display the current wallet balance', () => {
            cy.wait('@getMe');
            cy.contains(/\$1,540\.50|1540\.50/).should('be.visible');
        });

        it('should display Add Money and Withdraw toggle buttons', () => {
            cy.get('button').contains(/add|deposit/i).should('be.visible');
            cy.get('button').contains(/withdraw|send/i).should('be.visible');
        });

        it('should show the transaction history section', () => {
            cy.contains(/history|transactions/i).should('be.visible');
        });
    });

    describe('Rewards Page', () => {
        beforeEach(() => {
            cy.intercept('GET', '**/api/v1/rewards', {
                statusCode: 200,
                body: {
                    success: true,
                    data: { total_points: 1250, total_cashback: 12.50, lifetime_earnings: 25.00, total_transactions: 10 }
                }
            }).as('getRewards');

            cy.intercept('GET', '**/api/v1/rewards/offers', {
                statusCode: 200,
                body: {
                    success: true,
                    data: [
                        { id: '1', title: '5% off Amazon', description: 'Get 5% off on your next Amazon purchase', discount: 5.0, is_active: true }
                    ]
                }
            }).as('getOffers');

            cy.visit('/rewards');
        });

        it('should display reward statistics', () => {
            cy.wait(['@getRewards', '@getOffers']);
            cy.contains(/1,250|1250/i).should('be.visible'); // Points
            cy.contains(/\$12\.50|12\.50/i).should('be.visible'); // Cashback
        });

        it('should display active promotional offers', () => {
            cy.wait('@getOffers');
            cy.contains(/amazon/i).should('be.visible');
            cy.contains(/5%/i).should('be.visible');
        });
    });
});
