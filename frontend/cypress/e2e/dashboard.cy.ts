describe('Dashboard Page', () => {
    beforeEach(() => {
        // Mock authentication state
        window.localStorage.setItem('gatorpay_token', 'fake-jwt-token');
        
        // Mock user profile response
        cy.intercept('GET', '**/api/v1/auth/me', {
            statusCode: 200,
            body: {
                success: true,
                message: "User retrieved successfully",
                data: {
                    user: {
                        id: 'user-1',
                        first_name: 'Test',
                        last_name: 'User',
                        email: 'test@example.com'
                    },
                    wallet: {
                        balance: 1540.50,
                        currency: 'USD'
                    }
                }
            }
        }).as('getMe');

        // Mock recent transactions
        cy.intercept('GET', '**/api/v1/wallet/transactions*', {
            statusCode: 200,
            body: {
                success: true,
                data: {
                    transactions: [
                        { id: '1', amount: 50.00, type: 'credit', status: 'completed', created_at: new Date().toISOString() }
                    ],
                    total: 1
                }
            }
        }).as('getTransactions');

        cy.visit('/dashboard');
    });

    it('should display the user welcome message', () => {
        cy.wait('@getMe');
        cy.contains(/Welcome|Hello|Hi/i).should('be.visible');
        cy.contains('Test').should('be.visible');
    });

    it('should display the wallet balance', () => {
        cy.wait('@getMe');
        cy.contains(/\$1,540\.50|1540\.50/).should('be.visible');
    });

    it('should display quick action links', () => {
        cy.get('a[href*="transfer"], button').contains(/send|transfer/i).should('be.visible');
        cy.get('a[href*="bills"], button').contains(/pay|bills/i).should('be.visible');
        cy.get('a[href*="wallet"], button').contains(/add|withdraw/i).should('be.visible');
    });

    it('should list recent transactions', () => {
        cy.wait('@getTransactions');
        cy.contains(/\$50\.00|50\.00/).should('be.visible');
    });
});
